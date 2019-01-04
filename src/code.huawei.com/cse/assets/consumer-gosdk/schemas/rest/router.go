package rest

import (
	"code.huawei.com/cse/assets/consumer-gosdk/schemas/rest/service"
	"code.huawei.com/cse/model"
	"code.huawei.com/cse/route/consumer"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/server"
	"net"
	"net/http"
)

const Name = "http"

func init() {
	server.InstallPlugin(Name, NewServer)
}

type ChassisInvoker struct {
}

func (c *ChassisInvoker) Invoke(p *model.InvokeOption) error {
	service.Handle(p)
	return nil
}

type HttpServer struct {
	Opts server.Options
}

func (h *HttpServer) Register(interface{}, ...server.RegisterOption) (string, error) {
	return "", nil
}
func (h *HttpServer) Start() error {
	_, _, err := net.SplitHostPort(h.Opts.Address)
	if err != nil {
		return err
	}
	handler := consumer.NewRouter(&ChassisInvoker{})
	go func() {
		s := http.Server{
			Addr:    h.Opts.Address,
			Handler: handler,
		}
		if err := s.ListenAndServe(); err != nil {
			server.ErrRuntime <- err
			return
		}
	}()
	lager.Logger.Warnf("Http server listening on: %s", h.Opts.Address)
	return nil
}
func (h *HttpServer) Stop() error {
	return nil
}
func (h *HttpServer) String() string {
	return Name
}

func NewServer(opts server.Options) server.ProtocolServer {
	return &HttpServer{
		Opts: opts,
	}
}
