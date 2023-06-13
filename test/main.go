package main

import (
	"flag"
	"github.com/Shanghai-Lunara/pkg/zaplogger"
	ossrest "github.com/TyrandeCloud/aliyunoss-restapi"
	"github.com/TyrandeCloud/signals/pkg/signals"
)

func main() {
	var listenPort = flag.Int64("listenPort", 8085, "the port which was watched by the service")
	flag.Parse()
	stopCh := signals.SetupSignalHandler()
	s := ossrest.Run(int32(*listenPort))
	zaplogger.Sugar().Info("oss-proxy/service is running")
	<-stopCh
	zaplogger.Sugar().Info("oss-proxy/service trigger shutdown")
	s.Shutdown()
	<-stopCh
	zaplogger.Sugar().Info("oss-proxy/service shutdown gracefully")
}
