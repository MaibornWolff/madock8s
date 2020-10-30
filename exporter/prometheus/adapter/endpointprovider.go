package adapter

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

func GetEndpointBody(address string) (string, error) {
	ep := fmt.Sprintf("http://%v/metrics", address)
	u, err := url.Parse(ep)
	if nil != err {
		return "", errors.Wrap(err, "failed to parse url")
	}
	resp, err := http.Get(u.String())
	if err != nil {
		return "", errors.Wrap(err, "connect to prometheus endpoint failed")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "read prometheus body failed")
	}

	return string(body), nil
}
