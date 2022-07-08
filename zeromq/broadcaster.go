package zeromq

import (
	"context"
	"sync"

	"github.com/koykov/fastconv"
	"github.com/koykov/traceID/broadcaster"
	"github.com/pebbe/zmq4"
)

type Broadcaster struct {
	broadcaster.Base
	once  sync.Once
	ctx   *zmq4.Context
	sock  *zmq4.Socket
	topic []byte
	err   error
}

func (b *Broadcaster) Broadcast(_ context.Context, p []byte) (n int, err error) {
	b.once.Do(func() {
		conf := b.GetConfig()
		if len(conf.Topic) == 0 {
			conf.Topic = TopicNative
		}
		b.topic = []byte(conf.Topic)

		if b.ctx, b.err = zmq4.NewContext(); b.err != nil {
			return
		}
		if b.sock, b.err = b.ctx.NewSocket(zmq4.PUB); b.err != nil {
			return
		}
		if conf.HWM == 0 {
			conf.HWM = DefaultHWM
		}
		if b.err = b.sock.SetSndhwm(int(conf.HWM)); b.err != nil {
			return
		}
		if b.err = b.sock.Connect(conf.Addr); b.err != nil {
			return
		}
	})

	if b.err != nil {
		err = b.err
		return
	}

	var n1 int
	defer func() { n = n1 }()

	conf := b.GetConfig()
	if conf.Ping == 0 {
		conf.Ping = DefaultPing
	}
	for i := uint(0); i < conf.Ping; i++ {
		if n1, err = b.sock.SendBytes(fastconv.S2B(TopicService), zmq4.SNDMORE); err != nil {
			return
		}
		if n1, err = b.sock.SendBytes(svcPing, 0); err != nil {
			return
		}
	}

	if n1, err = b.sock.SendBytes(b.topic, zmq4.SNDMORE); err != nil {
		return
	}
	if n1, err = b.sock.SendBytes(p, 0); err != nil {
		return
	}

	return
}

var (
	svcPing = []byte("ping")
)
