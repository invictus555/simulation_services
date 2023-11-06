package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// GET
func DoHttpGetMethod(url, env string, body []byte) ([]byte, error) {
	client := http.Client{
		Timeout: time.Second * 3,
	}
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-tt-env", env)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ret, err
	}
	if resp.StatusCode != 200 {
		return ret, fmt.Errorf("resp status is %d, %s. url is %s", resp.StatusCode, url, string(body))
	}
	return ret, err
}

// GET
func DoHttpGetMethodV2(url string, body []byte) ([]byte, error) {
	client := http.Client{
		Timeout: time.Second * 3,
	}
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ret, err
	}
	if resp.StatusCode != 200 {
		return ret, fmt.Errorf("resp status is %d, %s. url is %s", resp.StatusCode, url, string(body))
	}
	return ret, err
}
