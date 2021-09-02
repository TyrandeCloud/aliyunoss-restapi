package aliyunoss_restapi

import (
	"io"
	"io/ioutil"
)

func ReadBody(body io.ReadCloser) ([]byte, error) {
	data, err := ioutil.ReadAll(body)
	defer body.Close()
	if err != nil {
		return nil, err
	}
	return data, nil
}
