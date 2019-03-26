package main

import _ "github.com/huaweicse/auth/adaptor/gochassis"

import (
	"net/http"
	_ "net/http/pprof"

	_ "code.huawei.com/cse/assets/consumer-gosdk/schemas/rest"
	"code.huawei.com/cse/config"
	"github.com/go-chassis/go-chassis"
	_ "github.com/go-chassis/go-chassis/configcenter"
	"github.com/go-chassis/go-chassis/core/lager"
	_ "github.com/huaweicse/cse-collector"
)

func main() {
	go func() {
		http.ListenAndServe(":30110", nil)
	}()
	//start all server you register in server/schemas.
	if err := chassis.Init(); err != nil {
		lager.Logger.Errorf("Init failed, err: %s", err)
		return
	}
	err := config.Init()
	if err != nil {
		panic(err)
	}
	chassis.Run()
}
