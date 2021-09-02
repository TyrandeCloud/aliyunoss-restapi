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
