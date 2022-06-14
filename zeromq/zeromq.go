package zeromq

const (
	TopicNative   = "TRACE_NATIVE"
	TopicProtobuf = "TRACE_PROTOBUF"
	TopicService  = "TRACE_SERVICE"
	DefaultHWM    = 1000
)

var _, _, _, _ = TopicNative, TopicProtobuf, TopicService, DefaultHWM
