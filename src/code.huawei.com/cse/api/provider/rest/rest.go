// package rest gives rest APIs info for sdk/mesher providers
package rest

// path parameters
const (
	Id           = "id"
	Ms           = "ms"
	Service      = "service"
	InstanceName = "instanceName"
	StatusCode   = "StatusCode"
)

// api path
const (
	Hello    = "/hello"
	SayHello = "/sayhello/{id}"
	Svc      = "/svc"               //return provider instance info
	Fail     = "/fail/{StatusCode}" //return specified status code
	//fail twice and succeed once;
	//return instance info if succeed, return specified status code if fail
	FailTwice = "/failTwice/{StatusCode}"
	//fail if instance name is equal to the {instanceName};
	//return instance info if succeed, return specified status code if fail
	FailInstance = "/failInstance/{instanceName}/{StatusCode}"
	//delay and return instance info
	Delay = "/delay/{ms}"
	//delay only if instance name is equal to the {instanceName}
	// return instance info after delay
	DelayInstance = "/delayInstance/{instanceName}/{ms}"
	FailV3        = "/failV3"
	Provider      = "/provider"
	ProxyTo       = "/proxyTo/{service}"
)
