package main

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
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

func parseBody(c *gin.Context, s interface{}) error {
	d := "json"
	format := c.Request.URL.Query()["d"]
	if len(format) > 0 {
		d = format[0]
	}

	b, err := getBodyAsBytes(c.Request.Body)
	if err != nil {
		return err
	}

	if d == "yaml" {
		//err = yaml.Unmarshal(b, s)
		err = errors.New("POST/PUT of yaml is not yet supported.")
	} else {
		err = json.Unmarshal(b, s)
	}

	return err
}

func renderHttp(code int, data interface{}, c *gin.Context) {
	f := "json"
	format := c.Request.URL.Query()["f"]
	if len(format) > 0 {
		f = format[0]
	}

	if f == "yaml" {
		out, _ := yaml.Marshal(data)
		c.Data(code, "text/yaml", out)
	} else {
		c.JSON(code, data)
	}
}
