package aliyunoss_restapi

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/Shanghai-Lunara/pkg/zaplogger"
	"github.com/gin-gonic/gin"
)

type Proxy struct {
	server *http.Server
	ctx    context.Context
	cancel context.CancelFunc
}

func Run(listenPort int32) *Proxy {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Proxy{
		ctx:    ctx,
		cancel: cancel,
	}
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	s.registerRouter(router)
	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", listenPort),
		Handler: router,
	}
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				zaplogger.Sugar().Info("Server closed under request")
			} else {
				zaplogger.Sugar().Info("Server closed unexpected err:", err)
			}
		}
	}()
	return s
}

func (p *Proxy) Shutdown() {
	p.cancel()
}

const (
	routerGet    = "/get"
	routerUpload = "/upload"
)

func (p *Proxy) registerRouter(router *gin.Engine) {
	router.GET(routerGet, p.get)
	router.GET(routerUpload, p.upload)
}

func (p *Proxy) checkRequest(req *ProxyRequest) error {
	if req.AccessEndpoint == "" {
		return fmt.Errorf("empty AccessEndpoint")
	}
	if req.AccessKeyID == "" {
		return fmt.Errorf("empty AccessKeyID")
	}
	if req.AccessKeySecret == "" {
		return fmt.Errorf("empty AccessKeySecret")
	}
	if req.Namespace == "" {
		return fmt.Errorf("empty Namespace")
	}
	if req.Channel == "" {
		return fmt.Errorf("empty Channel")
	}
	if req.Filename == "" {
		return fmt.Errorf("empty Filename")
	}
	return nil
}

func (p *Proxy) get(c *gin.Context) {
	req := &ProxyRequest{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	if err := p.checkRequest(req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	restClient, err := NewRestClient(&Option{
		EndPoint:        req.AccessEndpoint,
		AccessKeyID:     req.AccessKeyID,
		AccessKeySecret: req.AccessKeySecret,
		Context:         p.ctx,
	})
	if err != nil {
		zaplogger.Sugar().Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
	}
	downReader, err := restClient.Get().Verb("").Bucket(req.BucketName).Namespace(req.Namespace).Channel(req.Channel).Name(req.Filename).Download()
	if err != nil {
		zaplogger.Sugar().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	fileContentDisposition := fmt.Sprintf("attachment;filename=%s", req.Filename)
	c.Header("Content-Type", "text/plain")
	c.Header("Content-Disposition", fileContentDisposition)
	if _, err = io.Copy(c.Writer, downReader); err != nil {
		zaplogger.Sugar().Error(err)
	}
}

func (p *Proxy) upload(c *gin.Context) {
	req := &ProxyRequest{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	if err := p.checkRequest(req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	restClient, err := NewRestClient(&Option{
		EndPoint:        req.AccessEndpoint,
		AccessKeyID:     req.AccessKeyID,
		AccessKeySecret: req.AccessKeySecret,
		Context:         p.ctx,
	})
	if err != nil {
		zaplogger.Sugar().Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
	}
	if err := restClient.Get().Bucket(req.BucketName).Namespace(req.Namespace).Channel(req.Channel).Name(req.Filename).Body(c.Request.Body).Upload(); err != nil {
		zaplogger.Sugar().Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
