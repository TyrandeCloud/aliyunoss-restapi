package aliyunoss_restapi

import (
	"bytes"
	"testing"
)

func TestRequest(t *testing.T) {
	restOssClient := newAliyunOss(t)
	obj := []byte("saklsakslasas1")
	name := "abc.tz"
	if err := restOssClient.Verb("").Bucket(BucketName).Namespace(defaultNamespace).Channel(defaultChannel).Name(name).Body(obj).Upload(); err != nil {
		t.Error(err)
		return
	}

	reader, err := restOssClient.Verb("").Bucket(BucketName).Namespace(defaultNamespace).Channel(defaultChannel).Name(name).Download()
	if err != nil {
		t.Error(err)
		return
	}
	getObj, err := ReadBody(reader)
	if err != nil {
		t.Error(err)
		return
	}
	if !bytes.Equal(obj, getObj) {
		t.Errorf("expect:%v got:%v", string(obj), string(getObj))
		return
	}

	if err := restOssClient.Verb("").Bucket(BucketName).Namespace(defaultNamespace).Channel(defaultChannel).Name(name).Delete(); err != nil {
		t.Error(err)
		return
	}
}

func TestRequestDelete(t *testing.T) {
	restOssClient := newAliyunOss(t)
	obj := []byte("saklsakslasas1")
	name := "abc222.tz"
	if err := restOssClient.Verb("").Bucket(BucketName).Namespace(defaultNamespace).Channel(defaultChannel).Name(name).Body(obj).Upload(); err != nil {
		t.Error(err)
		return
	}

	if err := restOssClient.Verb("").Bucket(BucketName).Namespace(defaultNamespace).Channel(defaultChannel).Name(name).Delete(); err != nil {
		t.Error(err)
		return
	}

	_, err := restOssClient.Verb("").Bucket(BucketName).Namespace(defaultNamespace).Channel(defaultChannel).Name(name).Download()
	if err == nil {
		t.Errorf("key: %s exist", name)
		return
	}
}
