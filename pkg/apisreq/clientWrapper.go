package apisreq

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

//This wrapper will return as text the request so it can be processed by the AI
type HTTPClientWrapper struct {
	headers http.Header
}

func NewHTTPClientWrapper(headers http.Header) *HTTPClientWrapper {
	return &HTTPClientWrapper{
		headers: headers,
	}
}

func (w *HTTPClientWrapper) Get(url string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header = w.headers
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (w *HTTPClientWrapper) Post(url string, data map[string]interface{}) (string, error) {
	client := &http.Client{}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header = w.headers
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
