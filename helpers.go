package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

func getBodyAsBytes(body io.ReadCloser) ([]byte, error) {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

func getBodyAsString(body io.ReadCloser) (string, error) {
	b, err := getBodyAsBytes(body)
	if err != nil {
		return "", err
	}
	out := string(b)
	return out, nil
}

func getJSONBodyAsStruct(body io.ReadCloser, s interface{}) error {
	b, err := getBodyAsBytes(body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, s)
	if err != nil {
		return err
	}
	return nil
}
