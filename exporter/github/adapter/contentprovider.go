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

type File struct {
	Content      string
	Name         string
	Type         string `json:"type"`
	Path         string
	Download_url string `json:"download_url"`
}

type queryParams struct {
	page      int    // number of requested result page
	path      string // where to start searching for files
	per_page  int    // number of files returned per request, max. 100
	recursive bool   // false - search only in the provided path, true - check all children directories
	ref       string // branchName
}

func GetMdFiles(config map[string]string) ([]File, error) {
	baseurl := config["baseurl"]
	pattern := readPatternParam(config)

	client := &http.Client{}

	files, err := getMatchingFiles(baseurl, pattern)
	if err != nil {
		return []File{}, errors.Wrap(err, "failed GET files")
	}

	for idx, file := range files {
		content, err := getFileContent(client, file)
		if err != nil {
			return []File{}, errors.Wrap(err, "failed GET raw file from github")
		}
		file.Content = content
		files[idx] = file
	}
	return files, nil
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

func getMatchingFiles(baseurl string, pattern string) ([]File, error) {
	var files []File
	var rawFiles []File
	parsedURL, err := url.Parse(baseurl)
	if err != nil {
		return []File{}, err
	}
	req := configureRequest(parsedURL)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []File{}, errors.Wrapf(err, "failed GET request")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []File{}, errors.Wrap(err, "failed to read response body")
	}

	logrus.Info("RECEIVED BODY: ", string(body))

	err = json.Unmarshal(body, &rawFiles)
	if err != nil {
		return []File{}, errors.Wrap(err, "failed to parse json")
	}

	files = filterFiles(rawFiles, pattern)
	return files, nil
}

func filterFiles(files []File, pattern string) []File {
	var filteredFiles []File
	for _, file := range files {
		if file.Type == "file" {
			match, _ := regexp.MatchString(pattern, file.Name)
			if match {
				filteredFiles = append(filteredFiles, file)
			}
		}
	}
	return filteredFiles
}

func configureRequest(url *url.URL) *http.Request {
	req, _ := http.NewRequest("GET", url.String(), nil)
	req = setAuthToken(req)

	return req
}

func setAuthToken(req *http.Request) *http.Request {
	token := os.Getenv("AUTH_TOKEN")
	req.Header.Set("Authorization", fmt.Sprintf("token %v", token))
	return req
}

func getFileContent(client *http.Client, file File) (string, error) {
	req, _ := http.NewRequest("GET", file.Download_url, nil)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return string(body), nil
}
