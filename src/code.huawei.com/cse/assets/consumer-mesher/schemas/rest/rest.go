package rest

import (
	"encoding/json"
	"errors"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"

	"code.huawei.com/cse/common"
	"code.huawei.com/cse/model"
	"code.huawei.com/cse/pkg"
	"code.huawei.com/cse/util"
)

func Handle(p *model.InvokeOption, args ...interface{}) (*model.InstanceInfoResponse, error) {
	if p == nil {
		return nil, errors.New("InvokeOption is nil")
	}
	t := &model.InstanceInfoResponse{
		Result: make([]*model.SingleInstanceInfoResponse, 0),
	}

	for i := 1; i <= p.Times; i++ {
		var c *model.SingleInstanceInfoResponse
		c = Call(p.Service.ServiceName, p.URI)
		c.Num = i
		c.Time = util.Now()
		t.Result = append(t.Result, c)
	}
	return t, nil
}

func Call(provider, apiPath string) *model.SingleInstanceInfoResponse {
	u := util.FormatURL(common.SchemeHttp, "%s%s", provider, apiPath)
	c := &model.SingleInstanceInfoResponse{}
	resp, err := pkg.HTTPDo(http.MethodGet, u, nil, nil)
	if err != nil {
		c.Error = err.Error()
		c.Time = util.Now()
		return c
	}
	c.StatusCode = resp.StatusCode
	providerInfo := &model.InstanceStruct{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error("Read response body failed", err)
		c.Error = err.Error()
		c.Time = util.Now()
		return c
	}
	if resp.StatusCode == http.StatusOK {
		if util.IsContentTypeJson(resp.Header) {
			err = json.Unmarshal(body, providerInfo)
			if err == nil {
				c.Provider = providerInfo
			} else {
				c.Error = err.Error()
				glog.Error("Unmarshal communiction resp failed", err)
			}
		} else {
			c.Body = string(body)
		}
	} else {
		c.Error = string(body)
	}
	c.Time = util.Now()
	return c
}

type Invoker interface {
	Invoke(*model.InvokeOption) error
}

type MesherInvoker struct {
}

func (m *MesherInvoker) Invoke(p *model.InvokeOption) error {
	Handle(p)
	return nil
}
