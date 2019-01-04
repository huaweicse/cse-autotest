package main

import (
	"code.huawei.com/cse/assets/consumer-mesher/schemas/rest"
	"code.huawei.com/cse/route/consumer"
	"flag"
	"net/http"
	"os"

	"code.huawei.com/cse/config"
	"github.com/golang/glog"
)

func main() {
	os.Args = append(os.Args, "-logtostderr", "true")
	flag.Parse()
	err := config.Init()
	if err != nil {
		panic(err)
	}

	h := consumer.NewRouter(&rest.MesherInvoker{})
	addr := config.GlobalDef.Http.Listen
	glog.Warningf("Listening at: %s", addr)
	err = http.ListenAndServe(addr, h)
	if err != nil {
		glog.Error("Server start failed", err)
	}
}
