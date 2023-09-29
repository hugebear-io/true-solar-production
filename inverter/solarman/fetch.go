package solarman

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
	"github.com/hugebear-io/true-solar-production/util"
)

func prepareHttpRequest(method, url string, headers map[string]string, data interface{}) (*http.Request, error) {
	var req *http.Request
	var err error

	if data != nil {
		encoded_data, err := sonic.Marshal(data)
		if err != nil {
			return nil, err
		}

		buffered_data := bytes.NewBuffer(encoded_data)
		req, err = http.NewRequest(method, url, buffered_data)
		if err != nil {
			return nil, err
		}
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, err
		}
	}

	req.Header.Set("Content-Type", "application/json")
	for key, val := range headers {
		req.Header.Set(key, val)
	}

	return req, nil
}

func prepareHttpResponse[R interface{}](req *http.Request) (*R, int, error) {
	// request to endpoint
	client := &http.Client{
		Timeout: time.Duration(300) * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, res.StatusCode, err
	}
	defer res.Body.Close()

	// read a bytes of data
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// check empty response body
	if len(resBody) == 0 {
		return nil, res.StatusCode, nil
	}

	// decode response body
	var result R
	if err := util.Recast(resBody, &result); err != nil {
		if err := checkHTMLResponse(resBody); err != nil {
			return nil, http.StatusInternalServerError, err
		}

		return nil, http.StatusInternalServerError, err
	}

	return &result, res.StatusCode, nil
}
