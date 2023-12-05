package http

import (
	"io"
	"net/http"
	"net/url"

	"github.com/adrianriobo/qe-eventmanager/pkg/util/logging"
)

func GetFile(fileUrl string) ([]byte, error) {
	_, err := url.Parse(fileUrl)
	if err != nil {
		return nil, err
	}
	logging.Debugf("downloading %s", fileUrl)
	resp, err := httpClient().Get(fileUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// TODO pooling or delegate to lib
func httpClient() *http.Client {
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	return &client
}
