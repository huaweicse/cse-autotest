package pkg

import (
	"bytes"
	"net/http"
)

var c http.Client

//HTTPDo is a method used for http connection
func HTTPDo(method string, rawURL string, headers http.Header, body []byte) (resp *http.Response, err error) {
	req, err := http.NewRequest(method, rawURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err = c.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func init() {
	c = http.Client{
		Transport: http.DefaultTransport,
	}
}
