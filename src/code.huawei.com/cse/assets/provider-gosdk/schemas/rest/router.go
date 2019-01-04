package rest

import (
	"errors"
	"fmt"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/go-mesh/openlogging"
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
	"github.com/go-chassis/go-chassis/core/config"
	rf "github.com/go-chassis/go-chassis/server/restful"
)

var (
	Platform string
	Counter2 *pkg.Counter
)

type ATProvider struct {
	instance *model.InstanceStruct
	once     sync.Once
}

func (a *ATProvider) MySelf() (*model.InstanceStruct, error) {
	a.once.Do(func() {
		instanceName := runtime.HostName
		for name, proto := range config.GlobalDefinition.Cse.Protocols {
			instanceName = instanceName + fmt.Sprintf("_%s_%s", name, proto.Listen)
		}
		instanceName = strings.Replace(instanceName, ":", "_", -1)

		a.instance = &model.InstanceStruct{
			MicroService: &model.ServiceStruct{
				ServiceName: config.MicroserviceDefinition.ServiceDescription.Name,
				Version:     config.MicroserviceDefinition.ServiceDescription.Version,
				Application: config.MicroserviceDefinition.AppID,
			},
			InstanceName:  instanceName,
			InstanceAlias: os.Getenv(common.EnvInstanceAlias),
		}
	})
	return a.instance, nil
}

func (a *ATProvider) SayHello(b *rf.Context) {
	id := b.ReadPathParameter(providerRestApi.Id)
	msg := fmt.Sprintf("Hello: %s", id)
	openlogging.GetLogger().Info(msg)
	b.WriteHeader(http.StatusOK)
	b.Write([]byte(msg))
}

func (a *ATProvider) Communication(b *rf.Context) {
	openlogging.GetLogger().Info("Communication accept request")
	selfInfo, _ := a.MySelf()
	if err := b.WriteJSON(selfInfo, common.ContentTypeJson); err != nil {
		openlogging.GetLogger().Errorf("Unexpected err in response, err: %s", err)
	}
}

func (a *ATProvider) Fail(b *rf.Context) {
	a.failWithMsg(b, "")
}

func (a *ATProvider) failWithMsg(b *rf.Context, msg string) {
	status := b.ReadPathParameter(providerRestApi.StatusCode)
	code, err := strconv.Atoi(status)
	if err != nil {
		msg := fmt.Sprintf("bad path param: %s, %s", status, err.Error())
		openlogging.GetLogger().Error(msg)
		b.WriteError(http.StatusBadRequest, errors.New(msg))
		return
	}
	b.WriteHeader(code)
	if msg != "" {
		b.Write([]byte(msg))
	}
}

func (a *ATProvider) FailTwice(b *rf.Context) {
	if err := Counter2.Increase(); err != nil {
		a.Communication(b)
		return
	}
	msg := fmt.Sprintf("Failed! Please retry [%d] times", common.FailNum2+1-Counter2.Value())
	openlogging.GetLogger().Error(msg)
	a.failWithMsg(b, msg)
}

func (a *ATProvider) FailInstance(b *rf.Context) {
	targetInst := b.ReadPathParameter(providerRestApi.InstanceName)
	s, _ := a.MySelf()
	if s.IsMyName(targetInst) {
		msg := fmt.Sprintf("Failed! Because I am %s", targetInst)
		openlogging.GetLogger().Error(msg)
		a.failWithMsg(b, msg)
		return
	}
	a.Communication(b)
}

func (a *ATProvider) DelayMs(b *rf.Context) {
	ms := b.ReadPathParameter(providerRestApi.Ms)
	n, err := strconv.Atoi(ms)
	if err != nil {
		msg := fmt.Sprintf("bad path param: %s, %s", ms, err.Error())
		openlogging.GetLogger().Error(msg)
		b.WriteHeader(http.StatusBadRequest)
		b.Write([]byte(msg))
		return
	}
	openlogging.GetLogger().Infof("Delay %d ms", n)
	time.Sleep(time.Duration(n) * time.Millisecond)
	a.Communication(b)
}

func (a *ATProvider) DelayInstance(b *rf.Context) {
	targetInst := b.ReadPathParameter(providerRestApi.InstanceName)
	m, _ := a.MySelf()
	if m.IsMyName(targetInst) {
		a.DelayMs(b)
		return
	}
	a.Communication(b)
}
func (a *ATProvider) FailV3(b *rf.Context) {
	if config.SelfVersion != common.Version30 {
		a.Communication(b)
		return
	}
	msg := fmt.Sprintf("Failed! Because provider version == %s", common.Version30)
	openlogging.GetLogger().Error(msg)
	b.WriteHeader(http.StatusInternalServerError)
	b.Write([]byte(msg))
}

//func (a *ATProvider) ProxyTo(b *rf.Context) {
//	backendService := b.ReadPathParameter(providerRestApi.Service)
//	if backendService == "" {
//		msg := fmt.Sprintf("The service name to proxied is empty")
//		openlogging.GetLogger().Error(msg)
//		b.WriteHeader(http.StatusBadRequest)
//		b.Write([]byte(msg))
//		return
//	}
//	//c := &service.CommunicationProject{}
//	path := "/sayhello/1"
//	//openlogging.GetLogger().Infof("proxy to: %s, path: %s", backendService, path)
//	//resp, err := c.Communicate(common.PlatformGoSDK, backendService, path, "rest", 1)
//	//if err != nil {
//	//	openlogging.GetLogger().Error("Unexpected err in calling", err)
//	//	b.WriteHeader(http.StatusInternalServerError)
//	//	b.Write([]byte(err.Error()))
//	//	return
//	//}
//	//b.WriteJSON(resp, common.ContentTypeJson)
//
//	u := util.FormatURL(common.SchemeCse, "%s%s", backendService, path)
//	c := &model.CommunicateResponse{}
//	req, err := rest.NewRequest(http.MethodGet, u)
//	if err != nil {
//		c.Error = err.Error()
//		c.Time = util.Now()
//		b.WriteHeaderAndJSON(http.StatusInternalServerError, c, common.ContentTypeJson)
//		return
//	}
//	defer req.Close()
//
//	resp, err := service.RestInvoker.ContextDo(context.Background(), req)
//	if err != nil {
//		c.Error = err.Error()
//		c.Time = util.Now()
//		b.WriteHeaderAndJSON(http.StatusInternalServerError, c, common.ContentTypeJson)
//		return
//	}
//	defer resp.Close()
//
//	c.StatusCode = resp.GetStatusCode()
//	respBody := resp.ReadBody()
//	if resp.GetStatusCode() == http.StatusOK {
//		c.Message = string(respBody)
//	} else {
//		c.Error = string(respBody)
//	}
//	c.Time = util.Now()
//	b.WriteJSON(c, common.ContentTypeJson)
//}

func (a *ATProvider) URLPatterns() []rf.Route {
	return []rf.Route{
		{http.MethodGet, providerRestApi.SayHello, "SayHello"},
		{http.MethodGet, providerRestApi.Svc, "Communication"},
		{http.MethodGet, providerRestApi.Fail, "Fail"},
		{http.MethodGet, providerRestApi.FailTwice, "FailTwice"},
		{http.MethodGet, providerRestApi.FailInstance, "FailInstance"},
		{http.MethodGet, providerRestApi.Delay, "DelayMs"},
		{http.MethodGet, providerRestApi.DelayInstance, "DelayInstance"},
		//{http.MethodGet, providerRestApi.FailV3, "FailV3"},
	}
}

func init() {
	Counter2 = pkg.NewCounter(common.FailNum2 + 1)
}
