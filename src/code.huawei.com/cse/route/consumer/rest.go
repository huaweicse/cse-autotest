package consumer

import (
	"bytes"
	"code.huawei.com/cse/dispatcher"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"

	providerRestApi "code.huawei.com/cse/api/provider/rest"
	"code.huawei.com/cse/common"
	"code.huawei.com/cse/model"
)

type Invoker interface {
	Invoke(*model.InvokeOption) error
}

type Server struct {
	Invoker Invoker
}

func (s *Server) Svc(w http.ResponseWriter, r *http.Request) {
	p, err := dispatcher.GetQueryParams(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	p.URI = r.URL.Path

	err = s.Invoker.Invoke(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	resp, ok := p.Reply.(*model.InstanceInfoResponse)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("invoke response is not type *model.InstanceInfoResponse"))
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
	return Router
}
