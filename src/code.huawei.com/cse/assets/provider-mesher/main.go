package main

import (
	"flag"
	"github.com/golang/glog"
	"net/http"
	"os"

	"code.huawei.com/cse/assets/provider-mesher/schemas/rest"
	"code.huawei.com/cse/config"
)

func main() {
	os.Args = append(os.Args, "-logtostderr", "true")
	flag.Parse()
	err := config.Init()
	if err != nil {
		panic(err)
	}

	rest.InitRouter()
	addr := config.GlobalDef.Http.Listen
	glog.Warningf("Listening at: %s", addr)
	glog.Info("API: GET /hello")
	err = http.ListenAndServe(addr, rest.Router)
	if err != nil {
		glog.Error("Server start failed", err)
	}
}
