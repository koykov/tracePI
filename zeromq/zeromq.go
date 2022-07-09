package zeromq

const (
	TopicNative   = "TRACE_NATIVE"
	TopicProtobuf = "TRACE_PROTOBUF"
	TopicService  = "TRACE_SERVICE"
	DefaultHWM    = 1000
)

var btNative, btProtobuf, btService, bsPing []byte

func init() {
	btNative = []byte(TopicNative)
	btProtobuf = []byte(TopicProtobuf)
	btService = []byte(TopicService)
	bsPing = []byte("ping")
}
