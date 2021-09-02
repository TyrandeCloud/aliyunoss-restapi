package aliyunoss_restapi

import (
	"errors"
	"os"
)

const (
	AccessKeyID     = "ALIYUNOSS_ACCESS_KEY_ID"
	AccessKeySecret = "ALIYUNOSS_ACCESS_SECRET"
)

var (
	ErrNoAccessKeyID     = errors.New("unable to load configuration, ALIYUNOSS_ACCESS_KEY_ID must be defined")
	ErrNoAccessKeySecret = errors.New("unable to load configuration, ALIYUNOSS_ACCESS_SECRET must be defined")
)

func GetAccessKeyID() (string, error) {
	s := os.Getenv(AccessKeyID)
	if s == "" {
		return "", ErrNoAccessKeyID
	}
	return s, nil
}

func GetAccessKeySecret() (string, error) {
	s := os.Getenv(AccessKeySecret)
	if s == "" {
		return "", ErrNoAccessKeySecret
	}
	return s, nil
}
