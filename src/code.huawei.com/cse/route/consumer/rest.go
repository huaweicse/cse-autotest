package consumer

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	providerRestApi "code.huawei.com/cse/api/provider/rest"
	"code.huawei.com/cse/common"
	"code.huawei.com/cse/config"
	"code.huawei.com/cse/dispatcher"
	"code.huawei.com/cse/model"
	"github.com/go-mesh/openlogging"
	"github.com/gorilla/mux"
	"github.com/huaweicse/auth"
)

var signRequest auth.SignRequest

type Invoker interface {
	Invoke(*model.InvokeOption) error
}

type Server struct {
	Invoker Invoker
}

type CreateConfigApi struct {
	DimensionInfo string                 `json:"dimensionsInfo"`
	Items         map[string]interface{} `json:"items"`
}

type DeleteConfigApi struct {
	DimensionInfo string   `json:"dimensionsInfo"`
	Keys          []string `json:"keys"`
}

type Param struct {
	Data          map[string]interface{} `json:"data"`
	Keys          []string               `json:"key"`
	DimensionInfo string                 `json:"dimensionInfo"`
	Type          string                 `json:"type"`
}

func (s *Server) Svc(w http.ResponseWriter, r *http.Request) {
	p, err := dispatcher.GetQueryParams(r)
	if err != nil {
		errorMessage(w, http.StatusBadRequest, err.Error())
		return
	}
	p.URI = r.URL.Path

	err = s.Invoker.Invoke(p)
	if err != nil {
		errorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp, ok := p.Reply.(*model.InstanceInfoResponse)
	if !ok {
		errorMessage(w, http.StatusInternalServerError, "invoke response is not type *model.InstanceInfoResponse")
		return
	}
	jsonByte, err := json.Marshal(resp)
	if err != nil {
		errorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	tmpBuf := &bytes.Buffer{}
	if err = json.Indent(tmpBuf, jsonByte, "", "  "); err != nil {
		errorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set(common.HeaderContentType, common.ContentTypeJson)
	w.WriteHeader(http.StatusOK)
	w.Write(tmpBuf.Bytes())
}

func (s *Server) ConfigCenter(w http.ResponseWriter, r *http.Request) {
	// read Param
	param := Param{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = json.Unmarshal(body, &param)
	if err != nil {
		errorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}
	var resp *http.Response
	if param.Type == "add" {
		resp, err = ccAdd(param)
	} else if param.Type == "delete" {
		resp, err = ccDelete(param)
	} else {
		errorMessage(w, http.StatusBadRequest, "type of %s not supper , use add or delete ")
		return
	}

	//	 send data to cc config center
	if err != nil {
		errorMessage(w, http.StatusInternalServerError, err.Error())
		return
	}
	if resp == nil && resp.StatusCode > 300 && resp.StatusCode < 200 {
		m := fmt.Sprintf("%s config form cc failed , please check param", param.Type)
		openlogging.GetLogger().Error(m)
		errorMessage(w, http.StatusInternalServerError, m)
		return
	}

	body, _ = ioutil.ReadAll(resp.Body)
	openlogging.GetLogger().Infof("%s config form cc success : message %s", param.Type, string(body))

	w.Header().Set(common.HeaderContentType, common.ContentTypeJson)
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func ccAdd(param Param) (*http.Response, error) {

	// get url
	ccurl := fmt.Sprintf("%v/v3/%s/configuration/items", config.GlobalDef.ConfigCenterUrl, config.GlobalDef.ProjectId)
	// param
	configApi := CreateConfigApi{
		DimensionInfo: param.DimensionInfo,
		Items:         param.Data,
	}
	cb, _ := json.Marshal(configApi)
	payload := strings.NewReader(string(cb))
	return ccHTTPDO(payload, ccurl, http.MethodPost)
}

func ccDelete(param Param) (*http.Response, error) {
	if len(param.Keys) == 0 {
		openlogging.GetLogger().Error("not key need to delete , please check keys")
		return nil, errors.New("not key need to delete , please check keys")
	}
	// ccurl
	ccurl := fmt.Sprintf("%v/v3/%s/configuration/items", config.GlobalDef.ConfigCenterUrl, config.GlobalDef.ProjectId)

	configApi := DeleteConfigApi{
		DimensionInfo: param.DimensionInfo,
		Keys:          param.Keys,
	}
	cb, _ := json.Marshal(configApi)
	payload := strings.NewReader(string(cb))

	return ccHTTPDO(payload, ccurl, http.MethodDelete)
}
func ccHTTPDO(data *strings.Reader, ccurl, method string) (*http.Response, error) {
	req, _ := http.NewRequest(method, ccurl, data)
	if signRequest == nil {
		signRequest = newSignRequest()
	}
	err := signRequest(req)

	if err != nil {
		openlogging.GetLogger().Errorf("sign request failed,err :%v", ccurl, err)
		return nil, err
	}
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
	}
	client := &http.Client{Transport: transCfg}
	resp, err := client.Do(req)
	if err != nil {
		openlogging.GetLogger().Errorf("call url[%s] failed,err :%v", ccurl, err)
		return nil, err
	}
	return resp, nil
}

// errorMessage
func errorMessage(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	w.Write([]byte(message))
}
func newSignRequest() auth.SignRequest {
	request, _ := auth.GetShaAKSKSignFunc(config.GlobalDef.Cse.Credentials.AccessKey,
		config.GlobalDef.Cse.Credentials.SecretKey,
		config.GlobalDef.Region)
	return request
}
func NewRouter(i Invoker) *mux.Router {
	Router := mux.NewRouter()
	s := &Server{
		Invoker: i,
	}
	Router.HandleFunc(providerRestApi.Svc, s.Svc).Methods(http.MethodGet)
	Router.HandleFunc(providerRestApi.Fail, s.Svc).Methods(http.MethodGet)
	Router.HandleFunc(providerRestApi.FailTwice, s.Svc).Methods(http.MethodGet)
	Router.HandleFunc(providerRestApi.FailInstance, s.Svc).Methods(http.MethodGet)
	Router.HandleFunc(providerRestApi.Delay, s.Svc).Methods(http.MethodGet)
	Router.HandleFunc(providerRestApi.DelayInstance, s.Svc).Methods(http.MethodGet)
	Router.HandleFunc(providerRestApi.ConfigCenterAdd, s.ConfigCenter).Methods(http.MethodPost)
	return Router
}
