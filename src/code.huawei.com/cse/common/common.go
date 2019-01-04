package common

const (
	ConsumerGoSDK  = "sdkat_consumer_gosdk"
	ConsumerMesher = "sdkat_consumer_mesher"
	ProviderGoSDK  = "sdkat_provider_gosdk"
	ProviderMesher = "sdkat_provider_mesher"
)

const (
	EnvPlatform      = "PLATFORM"
	EnvServiceName   = "SERVICE_NAME"
	EnvVersion       = "VERSION"
	EnvAppId         = "APP_ID"
	EnvInstanceAlias = "INSTANCE_ALIAS"

	// for test
	EnvConsumerGoSdkAddr  = "SDKAT_CONSUMER_GOSDK_ADDR"
	EnvConsumerMesherAddr = "SDKAT_CONSUMER_MESHER_ADDR"
)

const (
	PlatformGoSDK  = "gosdk"
	PlatformMesher = "mesher"
)

// rest path params
const (
	ParamProtocol = "protocol"
	ParamApp      = "app"
	ParamProvider = "provider"
	ParamVersion  = "version"
	ParamTimes    = "times"
)

const (
	DefaultTimes = 1
)

const (
	SchemeCse  = "cse"
	SchemeHttp = "http"
)

const (
	GoSdk2GoSdk   = "GoSDK -> GoSDK"
	GoSdk2Mesher  = "GoSDK -> Mesher"
	Mesher2GoSdk  = "Mesher -> GoSDK"
	Mesher2Mesher = "Mesher -> Mesher"
)

const (
	Version20 = "2.0"
	Version30 = "3.0"
)

const (
	AppSDKAT = "sdkat"
)

const (
	AddrDefaultGoSDKConsumer  = "http://127.0.0.1:8000"
	AddrDefaultMesherConsumer = "http://127.0.0.1:8080"
)

const (
	HeaderContentType = "Content-Type"
)

const (
	ContentTypeJson = "application/json"
)

const (
	InstanceNumDefault = 2
)

const (
	CallTimes2     = 2
	CallTimes20    = 20
	CallTimes2Str  = "2"
	CallTimes20Str = "20"

	FailNum2 = 2
)

const (
	ProtocolRest    = "rest"
	ProtocolHttp    = "http"
	ProtocolHighway = "highway"
)
