package zeromq

import (
	"bytes"
	"context"

	"github.com/koykov/fastconv"
	"github.com/koykov/traceID/listener"
	"github.com/pebbe/zmq4"
)

type Listener struct {
	listener.Base
}

func (l Listener) Listen(ctx context.Context, out chan []byte) (err error) {
	conf := l.GetConfig()
	if len(conf.Topic) == 0 {
		conf.Topic = TopicNative
	}

	var (
		ztx *zmq4.Context
		zsk *zmq4.Socket
	)
	if ztx, err = zmq4.NewContext(); err != nil {
		return
	}
	if zsk, err = ztx.NewSocket(zmq4.SUB); err != nil {
		return
	}
	if conf.HWM == 0 {
		conf.HWM = DefaultHWM
	}
	if err = zsk.SetSndhwm(int(conf.HWM)); err != nil {
		return
	}
	if err = zsk.Connect(conf.Addr); err != nil {
		return
	}
	if err = zsk.SetSubscribe(conf.Topic); err != nil {
		return
	}

	for {
		select {
		case <-ctx.Done():
			_ = zsk.Close()
			return
		default:
			var p []byte
			for {
				if p, err = zsk.RecvBytes(0); err != nil || len(p) == 0 {
					continue
				}
				switch p {
				case fastconv.S2B(TopicNative), fastconv.S2B(TopicProtobuf):
					continue
				case fastconv.S2B(TopicService):
					var svc []byte
					if svc, err = zsk.RecvBytes(0); err != nil || len(p) == 0 {
						continue
					}
					switch svc {
					case svcPing:
						// do noting
					}
				}
				break
			}
			out <- p
		}
	}
}

func (l Listener) isTopic(p []byte) bool {
	return bytes.Equal(p, fastconv.S2B(TopicNative)) || bytes.Equal(p, fastconv.S2B(TopicProtobuf))
}

func (l Listener) isService(p []byte) bool {
	return bytes.Equal(p, fastconv.S2B(TopicService))
}
