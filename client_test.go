package aliyunoss_restapi

import (
	"bytes"
	"context"
	"testing"
)

const (
	BucketName = "lunara-dev"
)

func newAliyunOss(t *testing.T) *RestOssClient {
	key, err := GetAccessKeyID()
	if err != nil {
		t.Fatal(err)
	}
	sec, err := GetAccessKeySecret()
	if err != nil {
		t.Fatal(err)
	}
	opt := &Option{
		EndPoint:        "oss-cn-hangzhou.aliyuncs.com",
		AccessKeyID:     key,
		AccessKeySecret: sec,
		Context:         context.Background(),
	}
	a, err := NewRestClient(opt)
	if err != nil {
		t.Fatal(err)
	}
	return a
}

func TestListObjects(t *testing.T) {
	a := newAliyunOss(t)
	obj := []byte("2121xxass1")
	objKey := "test-put/luacripts.zip"
	if err := a.PutObject(BucketName, objKey, bytes.NewReader(obj)); err != nil {
		t.Error(err)
		return
	}
	if _, err := a.ListObjects(BucketName); err != nil {
		t.Error(err)
		return
	}
	if err := a.DeleteObject(BucketName, objKey); err != nil {
		t.Error(err)
	}
}

func TestPutAndGetObject(t *testing.T) {
	a := newAliyunOss(t)
	obj := []byte("2121xxass1")
	objKey := "test-put/luacripts.zip"
	if err := a.PutObject(BucketName, objKey, bytes.NewReader(obj)); err != nil {
		t.Error(err)
		return
	}
	reader, err := a.GetObject(BucketName, objKey)
	if err != nil {
		t.Error(err)
		return
	}
	data, err := ReadBody(reader)
	if err != nil {
		t.Error(err)
		return
	}
	if !bytes.Equal(data, obj) {
		t.Errorf("expect:%v got:%v", obj, data)
	}
	if err := a.DeleteObject(BucketName, objKey); err != nil {
		t.Error(err)
	}
}

func TestPutAndDeleteObject(t *testing.T) {
	a := newAliyunOss(t)
	obj := []byte("2121xxass1")
	objKey := "test-put/luacripts.zip"
	if err := a.PutObject(BucketName, objKey, bytes.NewReader(obj)); err != nil {
		t.Error(err)
		return
	}
	if err := a.DeleteObject(BucketName, objKey); err != nil {
		t.Error(err)
		return
	}
	_, err := a.GetObject(BucketName, objKey)
	if err == nil {
		t.Errorf("key: %s exist", objKey)
		return
	}
}
