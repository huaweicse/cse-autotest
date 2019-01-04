package main

import _ "github.com/huaweicse/auth/adaptor/gochassis"

import (
	_ "net/http/pprof"

	"code.huawei.com/cse/assets/provider-gosdk/schemas/rest"
	"github.com/go-chassis/go-chassis"
	_ "github.com/go-chassis/go-chassis/configcenter"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/server"
	_ "github.com/huaweicse/cse-collector"
	"net/http"
)

func main() {
	go func() {
		http.ListenAndServe(":30111", nil)
	}()
	chassis.RegisterSchema("rest", &rest.ATProvider{},
		server.WithSchemaID("ATProvider"))
	//start all server you register in server/schemas.
	if err := chassis.Init(); err != nil {
		lager.Logger.Errorf("Init failed, err: %s", err)
		return
	}

	chassis.Run()
}
