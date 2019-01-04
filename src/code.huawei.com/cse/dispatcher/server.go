package dispatcher

import (
	"bytes"
	"code.huawei.com/cse/common"
	"code.huawei.com/cse/model"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Handler func(c *model.InvokeOption, args ...interface{}) (*model.InstanceInfoResponse, error)

type Dispatcher struct {
	H Handler
}

func (this *Dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p, err := GetQueryParams(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	p.URI = r.URL.Path

	resp, err := this.H(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	jsonByte, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	tmpBuf := &bytes.Buffer{}
	if err = json.Indent(tmpBuf, jsonByte, "", "  "); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set(common.HeaderContentType, common.ContentTypeJson)
	w.WriteHeader(http.StatusOK)
	w.Write(tmpBuf.Bytes())
}

func GetQueryParams(r *http.Request) (*model.InvokeOption, error) {
	p := &model.InvokeOption{
		Service: &model.ServiceStruct{},
	}

	if p.Protocol = r.FormValue(common.ParamProtocol); p.Protocol == "" {
		p.Protocol = common.ProtocolRest
	}

	if p.Service.ServiceName = r.FormValue(common.ParamProvider); p.Service.ServiceName == "" {
		p.Service.ServiceName = common.ProviderGoSDK
	}

	p.Service.Application = r.FormValue(common.ParamApp)
	p.Service.Version = r.FormValue(common.ParamVersion)

	t := r.FormValue(common.ParamTimes)
	if t == "" {
		p.Times = common.DefaultTimes
	} else {
		n, err := strconv.Atoi(t)
		if err != nil || n < 1 {
			return nil, fmt.Errorf("times must be integer > 1, %v", err)
		}
		p.Times = n
	}
	return p, nil
}
