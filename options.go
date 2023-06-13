package aliyunoss_restapi

import "context"

type Option struct {
	EndPoint        string
	AccessKeyID     string
	AccessKeySecret string

	Bucket     string
	BucketName string

	Context context.Context
}

type ProxyRequest struct {
	AccessEndpoint  string `form:"access_endpoint"`
	AccessKeyID     string `form:"access_key_id"`
	AccessKeySecret string `form:"access_key_secret"`
	BucketName      string `form:"bucket_name"`
	Namespace       string `form:"namespace"`
	Channel         string `form:"channel"`
	Filename        string `form:"filename"`
}
