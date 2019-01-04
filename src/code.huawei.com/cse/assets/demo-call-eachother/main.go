package main

import (
	providerRestApi "code.huawei.com/cse/api/provider/rest"
	rest2 "code.huawei.com/cse/assets/consumer-mesher/schemas/rest"
	"code.huawei.com/cse/pkg/exchanger"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/golang/glog"
	"net/http"
	"os"
	"strconv"
	"time"

	"code.huawei.com/cse/assets/demo-call-eachother/schemas/rest"
	"code.huawei.com/cse/config"
)

const EnvServiceNum = "SERVICE_NUM"
const EnvTotalServiceNum = "TOTAL_SERVICE_NUM"
const EnvServiceNamePrefix = "SERVICE_NAME_PREFIX"

func main() {
	flag.Parse()
	e, err := beforeCall()
	if err != nil {
		panic(err)
	}
	err = config.Init()
	if err != nil {
		panic(err)
	}
	rest.InitRouter()
	addr := config.GlobalDef.Http.Listen
	glog.Warningf("Listening at: %s", addr)
	glog.Info("API: GET /hello")
	go func() {
		for {
			callEachOther(e)
			time.Sleep(5 * time.Minute)
		}
	}()
	err = http.ListenAndServe(addr, rest.Router)
	if err != nil {
		glog.Error("Server start failed", err)
	}
}

func callEachOther(e *exchanger.InterServiceExchger) {
	for _, v := range e.Level2Map[config.GlobalDef.Name] {
		time.Sleep(1 * time.Second)
		r := rest2.Call(v, providerRestApi.Hello)
		s, _ := json.Marshal(r)
		glog.Infof("%s -> %s, api: %s, result: %s", config.GlobalDef.Name, v, providerRestApi.Hello, string(s))
	}
	for k, v := range e.Level3Map[config.GlobalDef.Name] {
		for _, vv := range v {
			time.Sleep(1 * time.Second)
			targetApi := "/proxyto/" + vv
			r := rest2.Call(k, targetApi)
			s, _ := json.Marshal(r)
			glog.Infof("%s -> %s, api: %s, result: %s", config.GlobalDef.Name, k, targetApi, string(s))
		}
	}
}

func beforeCall() (*exchanger.InterServiceExchger, error) {
	totalNumString := os.Getenv(EnvTotalServiceNum)
	numString := os.Getenv(EnvServiceNum)
	if numString == "" || totalNumString == "" {
		return nil, fmt.Errorf("please set env %s and %s", EnvTotalServiceNum, EnvServiceNum)
	}
	totalNum, err := strconv.Atoi(totalNumString)
	if err != nil {
		return nil, err
	}
	num, err := strconv.Atoi(numString)
	if err != nil {
		return nil, err
	}
	if totalNum < num {
		return nil, fmt.Errorf("%s must bigger or equal to %s", EnvTotalServiceNum, EnvServiceNum)
	}
	prefix := "service"
	name := os.Getenv(EnvServiceNamePrefix)
	if name != "" {
		prefix = name
	}
	os.Setenv(common.ServiceName, prefix+numString)
	elems := make([]string, 0)
	for i := 1; i <= totalNum; i++ {
		elems = append(elems, prefix+strconv.Itoa(i))
	}
	fmt.Println(elems)
	return exchanger.NewExchanger(elems...), nil
}
