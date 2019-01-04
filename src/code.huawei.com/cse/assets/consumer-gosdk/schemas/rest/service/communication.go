package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"code.huawei.com/cse/common"
	"code.huawei.com/cse/model"
	"code.huawei.com/cse/util"

	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/core"
	chassisCommon "github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/lager"
	"golang.org/x/net/context"
)

func Handle(p *model.InvokeOption, args ...interface{}) (*model.InstanceInfoResponse, error) {
	t := &model.InstanceInfoResponse{
		Result: make([]*model.SingleInstanceInfoResponse, 0),
	}

	for i := 1; i <= p.Times; i++ {
		var c *model.SingleInstanceInfoResponse
		c = invoke(p.Service.Application, p.Service.ServiceName, p.Service.Version, p.URI)
		c.Num = i
		t.Result = append(t.Result, c)
	}
	p.Reply = t
	return t, nil
}

func invoke(app, provider, version, apiPath string) *model.SingleInstanceInfoResponse {
	c := &model.SingleInstanceInfoResponse{}
	defer func() { c.Time = util.Now() }()
	u := util.FormatURL(common.SchemeCse, "%s%s", provider, apiPath)
	req, err := rest.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		c.Error = err.Error()
		return c
	}

	tags := map[string]string{}
	if app != "" {
		tags[chassisCommon.BuildinTagApp] = app
	}
	if version != "" {
		tags[chassisCommon.BuildinTagVersion] = version
	}
	var tagsOp core.InvocationOption = func(o *core.InvokeOptions) { return }
	// if tag is not nil but length is zero, route rule takes no effect
	if len(tags) > 0 {
		tagsOp = core.WithRouteTags(tags)
	}
	resp, err := RestInvoker.ContextDo(context.TODO(), req, tagsOp)
	if err != nil {
		c.Error = err.Error()
		return c
	}
	defer resp.Body.Close()

	c.StatusCode = resp.StatusCode
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Error = err.Error()
		return c
	}
	if resp.StatusCode == http.StatusOK {
		if util.IsContentTypeJson(resp.Header) {
			providerInfo := &model.InstanceStruct{}
			err = json.Unmarshal(b, providerInfo)
			if err == nil {
				c.Provider = providerInfo
			} else {
				c.Error = err.Error()
				lager.Logger.Errorf("Unmarshal communiction resp failed: %v", err)
			}
		} else {
			c.Body = string(b)
		}
	} else {
		c.Error = string(b)
	}
	return c
}
