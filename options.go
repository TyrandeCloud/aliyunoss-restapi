package aliyunoss_restapi

import (
	"context"
	"fmt"
)

type Option struct {
	EndPoint        string
	AccessKeyID     string
	AccessKeySecret string

	Bucket     string
	BucketName string
	BaseUrl    string

	Context context.Context
}

type ProxyRequestParams struct {
	AccessEndpoint  string `form:"access_endpoint"`
	AccessKeyID     string `form:"access_key_id"`
	AccessKeySecret string `form:"access_key_secret"`
	BucketName      string `form:"bucket_name"`
	Namespace       string `form:"namespace"`
	Channel         string `form:"channel"`
	Filename        string `form:"filename"`
}

func (prp ProxyRequestParams) ToUrlParams() string {
	return fmt.Sprintf("access_endpoint=%s&access_key_id=%s&access_key_secret=%s&bucket_name=%s&namespace=%s&channel=%s&filename=%s",
		prp.AccessEndpoint, prp.AccessKeyID, prp.AccessKeySecret, prp.BucketName, prp.Namespace, prp.Channel, prp.Filename)
}
