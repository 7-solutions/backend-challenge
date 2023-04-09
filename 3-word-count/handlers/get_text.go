package handlers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type externalService interface {
	getText() (string, error)
}

type config struct {
	url string
}

func newExternal(url string) externalService {
	return config{url: url}
}

func (cfg config) getText() (string, error) {

	resp, err := http.Get(cfg.url)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Faile to call [GET] %s, got %v", cfg.url, err))
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error reading response body: %v", err))
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("Error from response external url got status code: %d", resp.StatusCode))
	}

	return string(body), nil

}
