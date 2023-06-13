package aliyunoss_restapi

import (
	"errors"
	"os"
)

const (
	AccessEndpoint  = "ALIYUNOSS_ACCESS_ENDPOINT"
	AccessKeyID     = "ALIYUNOSS_ACCESS_KEY_ID"
	AccessKeySecret = "ALIYUNOSS_ACCESS_SECRET"
	CustomProxy     = "ALIYUNOSS_CUSTOM_PROXY"
)

var (
	ErrNoAccessEndpoint  = errors.New("unable to load configuration, ALIYUNOSS_ACCESS_ENDPOINT must be defined")
	ErrNoAccessKeyID     = errors.New("unable to load configuration, ALIYUNOSS_ACCESS_KEY_ID must be defined")
	ErrNoAccessKeySecret = errors.New("unable to load configuration, ALIYUNOSS_ACCESS_SECRET must be defined")
)

func GetAccessEndpoint() (string, error) {
	s := os.Getenv(AccessEndpoint)
	if s == "" {
		return "", ErrNoAccessEndpoint
	}
	return s, nil
}

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

func GetCustomProxy() string {
	return os.Getenv(CustomProxy)
}
