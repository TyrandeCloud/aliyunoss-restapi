package aliyunoss_restapi

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
)

const (
	defaultTimeoutInSeconds = 10
	defaultMaxRetries       = 1

	defaultNamespace = "default"
	defaultChannel   = "channel"
)

const (
	ModelNative = false
	ModelProxy  = true
)

func NewRequest(c *RestOssClient) *Request {
	return &Request{
		c:          c,
		timeout:    time.Second * defaultTimeoutInSeconds,
		maxRetries: defaultMaxRetries,

		namespace: defaultNamespace,
		channel:   defaultChannel,

		isProxy: ModelNative,
	}
}

type Request struct {
	c *RestOssClient

	timeout    time.Duration
	maxRetries int

	// generic components accessible via method setters
	verb string

	bucketNameSet bool
	// bucketName will be used to specify the object's name or parent directory inside an aliyun-oss-client request
	bucketName string

	namespaceSet bool
	// namespace was used to make up the sub path
	namespace string

	channelSet bool
	channel    string

	resource     string
	resourceName string

	// output
	err  error
	body io.Reader

	// model
	isProxy bool
	// Only when the isProxy is ModelProxy, then opt should be set.
	// Otherwise, there is no need to specify it.
	opt Option
}

func (r *Request) SetModel(bo bool) {
	r.isProxy = bo
}

func (r *Request) SetOption(opt Option) {
	r.opt = opt
}

// Verb sets the verb this request will use.
func (r *Request) Verb(verb string) *Request {
	r.verb = verb
	return r
}

// Bucket will specify the ossOption Bucket
func (r *Request) Bucket(bucketName string) *Request {
	if r.err != nil {
		return r
	}
	if r.bucketNameSet {
		r.err = fmt.Errorf("namespace already set to %q, cannot change to %q", r.bucketName, bucketName)
		return r
	}
	if msgs := rest.IsValidPathSegmentName(bucketName); len(msgs) != 0 {
		r.err = fmt.Errorf("invalid bucketName %q: %v", bucketName, msgs)
		return r
	}
	r.bucketNameSet = true
	r.bucketName = bucketName
	return r
}

// Name sets the name of a resource to access (<resource>/[ns/<namespace>/]/[ch/<channel>/]<name>)
func (r *Request) Name(resourceName string) *Request {
	if r.err != nil {
		return r
	}
	if len(resourceName) == 0 {
		r.err = fmt.Errorf("resource name may not be empty")
		return r
	}
	if len(r.resourceName) != 0 {
		r.err = fmt.Errorf("resource name already set to %q, cannot change to %q", r.resourceName, resourceName)
		return r
	}
	if msgs := rest.IsValidPathSegmentName(resourceName); len(msgs) != 0 {
		r.err = fmt.Errorf("invalid resource name %q: %v", resourceName, msgs)
		return r
	}
	r.resourceName = resourceName
	return r
}

// Namespace applies the namespace scope to a request (<resource>/[ns/<namespace>/]/[ch/<channel>/]<name>)
func (r *Request) Namespace(namespace string) *Request {
	if r.err != nil {
		return r
	}
	if r.namespaceSet {
		r.err = fmt.Errorf("namespace already set to %q, cannot change to %q", r.namespace, namespace)
		return r
	}
	if msgs := rest.IsValidPathSegmentName(namespace); len(msgs) != 0 {
		r.err = fmt.Errorf("invalid namespace %q: %v", namespace, msgs)
		return r
	}
	r.namespaceSet = true
	r.namespace = namespace
	return r
}

// Channel applies the channel scope to a request (<resource>/[ns/<namespace>/]/[ch/<channel>/]<name>)
func (r *Request) Channel(channel string) *Request {
	if r.err != nil {
		return r
	}
	if r.channelSet {
		r.err = fmt.Errorf("channel already set to %q, cannot change to %q", r.channel, channel)
		return r
	}
	if msgs := rest.IsValidPathSegmentName(channel); len(msgs) != 0 {
		r.err = fmt.Errorf("invalid channel %q: %v", channel, msgs)
		return r
	}
	r.channelSet = true
	r.channel = channel
	return r
}

// Resource sets the resource to access (<resource>/[ns/<namespace>/]/[ch/<channel>/]<name>)
func (r *Request) Resource(resource string) *Request {
	if r.err != nil {
		return r
	}
	if len(r.resource) != 0 {
		r.err = fmt.Errorf("resource already set to %q, cannot change to %q", r.resource, resource)
		return r
	}
	if msgs := rest.IsValidPathSegmentName(resource); len(msgs) != 0 {
		r.err = fmt.Errorf("invalid resource %q: %v", resource, msgs)
		return r
	}
	r.resource = resource
	return r
}

// Timeout makes the request use the given duration as an overall timeout for the
// request. Additionally, if set passes the value as "timeout" parameter in URL.
func (r *Request) Timeout(d time.Duration) *Request {
	if r.err != nil {
		return r
	}
	r.timeout = d
	return r
}

// MaxRetries makes the request use the given integer as a ceiling of retrying upon receiving
// "Retry-After" headers and 429 status-code in the response. The default is 10 unless this
// function is specifically called with a different value.
// A zero maxRetries prevent it from doing retires and return an error immediately.
func (r *Request) MaxRetries(maxRetries int) *Request {
	if maxRetries < 0 {
		maxRetries = 0
	}
	r.maxRetries = maxRetries
	return r
}

// Body makes the request use obj as the body. Optional.
// If obj is a string, try to read a file of that name.
// If obj is a []byte, send it directly.
// If obj is an io.Reader, use it directly.
// If obj is a runtime.Object, marshal it correctly, and set Content-Type header.
// If obj is a runtime.Object and nil, do nothing.
// Otherwise, set an error.
func (r *Request) Body(obj interface{}) *Request {
	if r.err != nil {
		return r
	}
	switch t := obj.(type) {
	case string:
		data, err := os.ReadFile(t)
		if err != nil {
			r.err = err
			return r
		}
		glogBody("Request Body", data)
		r.body = bytes.NewReader(data)
	case []byte:
		glogBody("Request Body", t)
		r.body = bytes.NewReader(t)
	case io.Reader:
		r.body = t
	default:
		r.err = fmt.Errorf("unknown type used for body: %+v", obj)
	}
	return r
}

// truncateBody decides if the body should be truncated, based on the glog Verbosity.
func truncateBody(body string) string {
	max := 0
	switch {
	case bool(klog.V(10).Enabled()):
		return body
	case bool(klog.V(9).Enabled()):
		max = 10240
	case bool(klog.V(8).Enabled()):
		max = 1024
	}

	if len(body) <= max {
		return body
	}

	return body[:max] + fmt.Sprintf(" [truncated %d chars]", len(body)-max)
}

// glogBody logs a body output that could be either JSON or protobuf. It explicitly guards against
// allocating a new string for the body output unless necessary. Uses a simple heuristic to determine
// whether the body is printable.
func glogBody(prefix string, body []byte) {
	if klog.V(8).Enabled() {
		if bytes.IndexFunc(body, func(r rune) bool {
			return r < 0x0a
		}) != -1 {
			klog.Infof("%s:\n%s", prefix, truncateBody(hex.Dump(body)))
		} else {
			klog.Infof("%s: %s", prefix, truncateBody(string(body)))
		}
	}
}

func (r *Request) List(opts ...oss.Option) (oss.ListObjectsResult, error) {
	if r.err != nil {
		return oss.ListObjectsResult{}, r.err
	}
	if r.namespaceSet {
		opts = append(opts, oss.Prefix(r.namespace))
	}
	switch r.isProxy {
	case ModelNative:
		return r.c.ListObjects(r.bucketName, opts...)
	case ModelProxy:
	}
	return oss.ListObjectsResult{}, fmt.Errorf("invalid list")
}

func (r *Request) FullRootPath() string {
	return path.Join(r.namespace, r.channel, r.resourceName)
}

func (r *Request) Upload() error {
	if r.err != nil {
		return r.err
	}
	switch r.isProxy {
	case ModelNative:
		return r.c.PutObject(r.bucketName, r.FullRootPath(), r.body)
	case ModelProxy:
		params := &ProxyRequestParams{
			AccessEndpoint:  r.opt.EndPoint,
			AccessKeyID:     r.opt.AccessKeyID,
			AccessKeySecret: r.opt.AccessKeySecret,
			BucketName:      r.bucketName,
			Namespace:       r.namespace,
			Channel:         r.channel,
			Filename:        r.resource,
		}
		resp, err := http.NewRequest("GET", fmt.Sprintf("%s/get?%s", r.opt.BaseUrl, params.ToUrlParams()), r.body)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		return nil
	}
	return fmt.Errorf("invalid Upload")
}

func (r *Request) Download() (io.ReadCloser, error) {
	if r.err != nil {
		return nil, r.err
	}
	switch r.isProxy {
	case ModelNative:
		return r.c.GetObject(r.bucketName, r.FullRootPath())
	case ModelProxy:
		params := &ProxyRequestParams{
			AccessEndpoint:  r.opt.EndPoint,
			AccessKeyID:     r.opt.AccessKeyID,
			AccessKeySecret: r.opt.AccessKeySecret,
			BucketName:      r.bucketName,
			Namespace:       r.namespace,
			Channel:         r.channel,
			Filename:        r.resource,
		}
		resp, err := http.NewRequest("GET", fmt.Sprintf("%s/upload?%s", r.opt.BaseUrl, params.ToUrlParams()), r.body)
		if err != nil {
			return nil, err
		}
		return resp.Body, nil
	}
	return nil, fmt.Errorf("invalid Download")
}

func (r *Request) Delete() error {
	if r.err != nil {
		return r.err
	}
	switch r.isProxy {
	case ModelNative:
		return r.c.DeleteObject(r.bucketName, r.FullRootPath())
	case ModelProxy:
	}
	return fmt.Errorf("invalid Delete")
}

func (r *Request) ForceDelete(key string) error {
	if r.err != nil {
		return r.err
	}
	switch r.isProxy {
	case ModelNative:
		return r.c.DeleteObject(r.bucketName, key)
	case ModelProxy:
	}
	return fmt.Errorf("invalid ForceDelete")
}
