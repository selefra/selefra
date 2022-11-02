package utils

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type Header struct {
	Key   string
	Value string
}

func Request(ctx context.Context, method string, _url string, body []byte, headers ...Header) ([]byte, error) {
	client := &http.Client{}
	sBody := strings.NewReader(string(body))
	request, err := http.NewRequestWithContext(ctx, method, _url, sBody)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	for _, header := range headers {
		request.Header.Add(header.Key, header.Value)
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("code not equal 200")
	}
	rByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("read body err :" + err.Error())
	}
	return rByte, err
}
