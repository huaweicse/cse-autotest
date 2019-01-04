package rest

import (
	"code.huawei.com/cse/config"
	"code.huawei.com/cse/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	providerRestApi "code.huawei.com/cse/api/provider/rest"
	"code.huawei.com/cse/common"
	"code.huawei.com/cse/model"
	"code.huawei.com/cse/pkg"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

var (
	Router           *mux.Router
	Counter2         *pkg.Counter
	providerInstance *model.InstanceStruct
	once             sync.Once
)

func mySelf() (*model.InstanceStruct, error) {
	once.Do(func() {
		instanceName := fmt.Sprintf("%s_%s_%s", util.HostName(), "http", config.GlobalDef.Http.Listen)
		instanceName = strings.Replace(instanceName, ":", "_", -1)

		providerInstance = &model.InstanceStruct{
			MicroService: &model.ServiceStruct{
				ServiceName: config.GlobalDef.Name,
				Version:     config.GlobalDef.Version,
				Application: config.GlobalDef.AppID,
			},
			InstanceName:  instanceName,
			InstanceAlias: os.Getenv(common.EnvInstanceAlias),
		}
	})
	return providerInstance, nil
}

func communication(w http.ResponseWriter, r *http.Request) {
	glog.Info("Communication accept request")
	selfInfo, _ := mySelf()
	jsonByte, err := json.Marshal(selfInfo)
	if err != nil {
		glog.Error("Json marshal failed", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set(common.HeaderContentType, common.ContentTypeJson)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonByte)
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params[providerRestApi.Id]
	msg := fmt.Sprintf("Hello: %s", id)
	glog.Info(msg)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msg))
}

func fail(w http.ResponseWriter, r *http.Request) {
	failWithMsg(w, r, "")
}

func failWithMsg(w http.ResponseWriter, r *http.Request, msg string) {
	params := mux.Vars(r)
	status := params[providerRestApi.StatusCode]
	statusCode, err := strconv.Atoi(status)
	if err != nil {
		msg := fmt.Sprintf("bad path param: %s, %s", status, err.Error())
		glog.Error(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg))
		return
	}
	w.WriteHeader(statusCode)
	if msg != "" {
		w.Write([]byte(msg))
	}
}

func failtwice(w http.ResponseWriter, r *http.Request) {
	if err := Counter2.Increase(); err != nil {
		communication(w, r)
		return
	}
	msg := fmt.Sprintf("Failed! Please retry [%d] times", common.FailNum2+1-Counter2.Value())
	glog.Error(msg)
	failWithMsg(w, r, msg)
}

func failInstance(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	targetInst := params[providerRestApi.InstanceName]
	selfInfo, _ := mySelf()
	if selfInfo.IsMyName(targetInst) {
		msg := fmt.Sprintf("Failed, because I am %s", targetInst)
		glog.Error(msg)
		failWithMsg(w, r, msg)
		return
	}
	communication(w, r)
}

func delayMs(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ms := params[providerRestApi.Ms]
	n, err := strconv.Atoi(ms)
	if err != nil {
		msg := fmt.Sprintf("bad path param: %s, %s", ms, err.Error())
		glog.Error(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg))
		return
	}
	glog.Infof("Delay %d ms", n)
	time.Sleep(time.Duration(n) * time.Millisecond)
	communication(w, r)
}

func delayInstance(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	targetInst := params[providerRestApi.InstanceName]
	selfInfo, _ := mySelf()
	if selfInfo.IsMyName(targetInst) {
		delayMs(w, r)
		return
	}
	communication(w, r)
	return
}

var p *model.InstanceStruct

func pushProviderInfo(w http.ResponseWriter, r *http.Request) {
	providerInfo := &model.InstanceStruct{}
	jsonByte, err := ioutil.ReadAll(r.Body)
	if err != nil {
		glog.Error("Read body failed", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	err = json.Unmarshal(jsonByte, providerInfo)
	if err != nil {
		glog.Error("unmarshal provider failed", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	p = providerInfo
	glog.Info("pushProviderInfo success")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success"))
}

func popProviderInfo(w http.ResponseWriter, r *http.Request) {
	jsonByte, err := json.Marshal(p)
	if err != nil {
		glog.Error("Json marshal failed", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	glog.Info("popProviderInfo success")
	w.Header().Set(common.HeaderContentType, common.ContentTypeJson)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonByte)
}

func failV3(w http.ResponseWriter, r *http.Request) {
	selfInfo, err := mySelf()
	if err != nil {
		glog.Error("Get self info failed", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if selfInfo.MicroService.Version != common.Version30 {
		communication(w, r)
		return
	}
	msg := fmt.Sprintf("Failed! Because provider version == %s", common.Version30)
	glog.Error(msg)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(msg))
}

func hello(w http.ResponseWriter, r *http.Request) {
	msg := "Hello, Mesher!"
	glog.Info(msg)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msg + "\n"))
}

func InitRouter() {
	Router = mux.NewRouter()
	Router.HandleFunc(providerRestApi.Hello, hello).Methods(http.MethodGet)
	Router.HandleFunc(providerRestApi.SayHello, sayHello).Methods(http.MethodGet)
	Router.HandleFunc(providerRestApi.Svc, communication).Methods(http.MethodGet)
	Router.HandleFunc(providerRestApi.Fail, fail).Methods(http.MethodGet)
	Router.HandleFunc(providerRestApi.FailTwice, failtwice).Methods(http.MethodGet)
	Router.HandleFunc(providerRestApi.FailInstance, failInstance).Methods(http.MethodGet)
	Router.HandleFunc(providerRestApi.Delay, delayMs).Methods(http.MethodGet)
	Router.HandleFunc(providerRestApi.DelayInstance, delayInstance).Methods(http.MethodGet)
	Router.HandleFunc(providerRestApi.FailV3, failV3).Methods(http.MethodGet)
	Router.HandleFunc(providerRestApi.Provider, pushProviderInfo).Methods(http.MethodPost)
	Router.HandleFunc(providerRestApi.Provider, popProviderInfo).Methods(http.MethodGet)
}

func init() {
	Counter2 = pkg.NewCounter(common.FailNum2 + 1)
	p = &model.InstanceStruct{}
}
