package adapter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type GitlabFile struct {
	ID      string
	Name    string
	Type    string `json:"type"`
	Path    string
	Mode    string
	Content string
}

type queryParams struct {
	page      int    // number of requested result page
	path      string // where to start searching for files
	per_page  int    // number of files returned per request, max. 100
	recursive bool   // false - search only in the provided path, true - check all children directories
	ref       string // branchName
}

func GetMdFiles(config map[string]string) ([]GitlabFile, error) {
	baseurl := config["baseurl"]
	path := config["path"]
	ref := config["ref"]
	pattern := readPatternParam(config)
	recursive := readRecursiveParam(config)

	client := &http.Client{}

	params := queryParams{
		page:      0,
		path:      path,
		per_page:  100,
		recursive: recursive,
		ref:       ref,
	}

	files, err := listRepositoryTree(client, baseurl, params)
	if err != nil {
		return []GitlabFile{}, errors.Wrap(err, "failed GET repository tree")
	}

	mdFiles := filterByFileName(files, pattern)

	for idx, file := range mdFiles {
		content, err := getRawFileContent(client, baseurl, file, params.ref)
		if err != nil {
			return []GitlabFile{}, errors.Wrap(err, "failed GET raw file from gitlab")
		}
		file.Content = content
		mdFiles[idx] = file
	}
	return mdFiles, nil
}

func readPatternParam(config map[string]string) string {
	pattern := `.\.md`
	if value, ok := config["pattern"]; ok {
		pattern = value
	}
	return pattern
}

func readRecursiveParam(config map[string]string) bool {
	recursive := true
	if value, ok := config["recursive"]; ok {
		parsedRecursive, err := strconv.ParseBool(value)
		if err != nil {
			logrus.Errorf("failed to parse config[\"recursive\"] with", err)
		} else {
			recursive = parsedRecursive
		}
	}
	return recursive
}

// Get files where names match provided pattern
func filterByFileName(files []GitlabFile, pattern string) []GitlabFile {
	var result []GitlabFile
	for _, file := range files {
		match, _ := regexp.MatchString(pattern, file.Name)
		if match {
			result = append(result, file)
		}
	}
	return result
}

// Get the COMPLETE list of repository files and directories in a project for all pages.
// NOTE: Gitlab API uses paging to limit the number of objects in response body.
// Reference: https://docs.gitlab.com/ee/api/repositories.html#list-repository-tree
func listRepositoryTree(client *http.Client, baseurl string, params queryParams) ([]GitlabFile, error) {
	var files []GitlabFile
	currentPage := 0
	totalPages := 1

	for currentPage != totalPages {
		parsedURL, _ := url.Parse(baseurl + "tree")
		req := configureRequest(parsedURL, params)
		resp, err := client.Do(req)
		if err != nil {
			return []GitlabFile{}, errors.Wrapf(err, "failed GET request")
		}

		var newFiles []GitlabFile
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return []GitlabFile{}, errors.Wrap(err, "failed to read response body")
		}

		err = json.Unmarshal(body, &newFiles)
		if err != nil {
			return []GitlabFile{}, errors.Wrap(err, "failed to parse json")
		}

		files = append(files, newFiles...)

		currentPage, _ = strconv.Atoi(resp.Header.Get("x-page"))
		totalPages, _ = strconv.Atoi(resp.Header.Get("x-total-pages"))
		params.page, _ = strconv.Atoi(resp.Header.Get("x-next-page"))
	}

	return files, nil
}

func configureRequest(url *url.URL, queryParams queryParams) *http.Request {
	req, _ := http.NewRequest("GET", url.String(), nil)
	req = setAuthToken(req)

	q := req.URL.Query()
	q.Add("path", queryParams.path)
	q.Add("recursive", strconv.FormatBool(queryParams.recursive))
	q.Add("ref", queryParams.ref)
	q.Add("per_page", strconv.Itoa(queryParams.per_page))
	q.Set("page", strconv.Itoa(queryParams.page))
	req.URL.RawQuery = q.Encode()

	return req
}

func setAuthToken(req *http.Request) *http.Request {
	token := os.Getenv("GITLAB_TOKEN")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	return req
}

// Get raw file from repository
// reference: https://docs.gitlab.com/ee/api/repository_files.html#get-raw-file-from-repository
func getRawFileContent(client *http.Client, baseurl string, file GitlabFile, ref string) (string, error) {
	rawurl := baseurl + "files/"
	rawurl += url.QueryEscape(file.Path)
	rawurl += "/raw"
	parsedURL, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}

	req, _ := http.NewRequest("GET", parsedURL.String(), nil)
	req = setAuthToken(req)

	q := req.URL.Query()
	q.Add("ref", ref)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return string(body), nil
}
