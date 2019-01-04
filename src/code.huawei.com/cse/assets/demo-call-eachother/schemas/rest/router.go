package rest

import (
	providerRestApi "code.huawei.com/cse/api/provider/rest"
	"code.huawei.com/cse/common"
	"code.huawei.com/cse/config"
	"code.huawei.com/cse/pkg"
	"code.huawei.com/cse/util"
	"fmt"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

var (
	Router *mux.Router
)

func proxyTo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	service := params[providerRestApi.Service]
	u := util.FormatURL(common.SchemeHttp, "%s%s", service, "/hello")
	resp, err := pkg.HTTPDo(http.MethodGet, u, nil, nil)
	if err != nil {
		glog.Error(err)
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(err.Error()))
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	msg := fmt.Sprintf("%s - %s, api: %s, result: %s", config.GlobalDef.Name, service, providerRestApi.Hello, string(body))
	glog.Infof(msg)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msg))
}

func hello(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("This is %s", config.GlobalDef.Name)
	glog.Info(msg)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msg + "\n"))
}

func InitRouter() {
	Router = mux.NewRouter()
	Router.HandleFunc(providerRestApi.Hello, hello).Methods(http.MethodGet)
	Router.HandleFunc(providerRestApi.ProxyTo, proxyTo).Methods(http.MethodGet)
}
