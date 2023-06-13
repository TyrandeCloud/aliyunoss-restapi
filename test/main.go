package main

import (
	"github.com/Shanghai-Lunara/pkg/zaplogger"
	ossrest "github.com/TyrandeCloud/aliyunoss-restapi"
	"github.com/TyrandeCloud/signals/pkg/signals"
	flag "github.com/spf13/pflag"
)

func main() {
	var listenPort = flag.Int32("listenPort", 8085, "the port which was watched by the service")
	flag.Parse()
	stopCh := signals.SetupSignalHandler()
	s := ossrest.Run(*listenPort)
	zaplogger.Sugar().Info("oss-proxy/service is running")
	<-stopCh
	zaplogger.Sugar().Info("oss-proxy/service trigger shutdown")
	s.Shutdown()
	<-stopCh
	zaplogger.Sugar().Info("oss-proxy/service shutdown gracefully")

}
