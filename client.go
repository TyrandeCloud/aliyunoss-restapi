package aliyunoss_restapi

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/Shanghai-Lunara/pkg/zaplogger"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func NewRestClient(opt *Option) (*RestOssClient, error) {
	client, err := oss.New(opt.EndPoint, opt.AccessKeyID, opt.AccessKeySecret)
	if err != nil {
		zaplogger.Sugar().Error(err)
		return nil, err
	}
	ctx, cancel := context.WithCancel(opt.Context)
	return &RestOssClient{
		option: opt,
		ctx:    ctx,
		cancel: cancel,
		client: client,
	}, nil
}

type RestOssClient struct {
	mu     sync.RWMutex
	option *Option
	ctx    context.Context
	cancel context.CancelFunc

	client *oss.Client
}

func (a *RestOssClient) ListBucket() error {
	lsRes, err := a.client.ListBuckets()
	if err != nil {
		return err
	}
	for _, bucket := range lsRes.Buckets {
		fmt.Println("Buckets:", bucket.Name)
	}
	return nil
}

func (a *RestOssClient) ListObjects(bucketName string, opts ...oss.Option) (oss.ListObjectsResult, error) {
	bucket, err := a.client.Bucket(bucketName)
	if err != nil {
		return oss.ListObjectsResult{}, err
	}
	return bucket.ListObjects(opts...)
}

func (a *RestOssClient) PutObject(bucketName, objectKey string, reader io.Reader) error {
	bucket, err := a.client.Bucket(bucketName)
	if err != nil {
		return err
	}
	return bucket.PutObject(objectKey, reader)
}

func (a *RestOssClient) GetObject(bucketName, objectKey string) (io.ReadCloser, error) {
	bucket, err := a.client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}
	return bucket.GetObject(objectKey)
}

func (a *RestOssClient) DeleteObject(bucketName, objectKey string) error {
	bucket, err := a.client.Bucket(bucketName)
	if err != nil {
		return err
	}
	return bucket.DeleteObject(objectKey)
}

func (a *RestOssClient) Verb(verb string) *Request {
	r := NewRequest(a).Verb(verb)
	if a.option.BaseUrl != "" {
		r.SetModel(ModelProxy)
		r.SetOption(r.opt)
	}
	return r
}

// Post begins a POST request. Short for a.Verb("POST").
func (a *RestOssClient) Post() *Request {
	return a.Verb("POST")
}

// Put begins a PUT request. Short for a.Verb("PUT").
func (a *RestOssClient) Put() *Request {
	return a.Verb("PUT")
}

// Get begins a GET request. Short for a.Verb("GET").
func (a *RestOssClient) Get() *Request {
	return a.Verb("GET")
}

// Delete begins a DELETE request. Short for a.Verb("DELETE").
func (a *RestOssClient) Delete() *Request {
	return a.Verb("DELETE")
}
