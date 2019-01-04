package model

//微服务信息
//不用于注册，只用于provider向consumer反馈自身信息
type ServiceStruct struct {
	Application string `json:"application"`
	ServiceName string `json:"service_name"`
	Version     string `json:"version"`
}

//实例信息
//不用于注册，只用于provider向consumer反馈自身信息
type InstanceStruct struct {
	MicroService  *ServiceStruct `json:"micro_service"`            //所属微服务
	InstanceName  string         `json:"instance_name,omitempty"`  //实例名，hostname_ip_port
	InstanceAlias string         `json:"instance_alias,omitempty"` //由环境变量INSTANCE_ALIAS设置
}

// consumer调用参数
type InvokeOption struct {
	Times      int               `json:"times,omitempty"`       //调用次数
	Service    *ServiceStruct    `json:"service,omitempty"`     //要访问的微服务
	Protocol   string            `json:"protocol,omitempty"`    //协议
	HttpMethod string            `json:"http_method,omitempty"` //http
	URI        string            `json:"uri,omitempty"`         //http
	Body       string            `json:"body,omitempty"`        //http
	Header     map[string]string `json:"header,omitempty"`      //http
	Property   map[string]string `json:"property,omitempty"`    //其他
	Reply      interface{}
}

//consumer的调用有两种：
//1,测协议。需要测试协议client/server通信是否正常，需要关注不同情况下的response
//2,测治理。只用到最基本的通信能力，如rest GET请求，server端返回自身信息即可。

//测治理的response，provider只需要反馈自身信息即可
type InstanceInfoResponse struct {
	Result []*SingleInstanceInfoResponse
}

type SingleInstanceInfoResponse struct {
	Num        int             `json:"num,omitempty"`        //调用序号
	Time       string          `json:"time,omitempty"`       //response 时间
	StatusCode int             `json:"statusCode,omitempty"` //http status code
	Error      string          `json:"error,omitempty"`      //错误信息，调用出错测error不为空
	Provider   *InstanceStruct `json:"provider,omitempty"`   //表示哪个provider实例回复了本次请求
	Body       string          `json:"body,omitempty"`       //如果无法unmarshal，则将body放到这里
}
