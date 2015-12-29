package client

import (
	"io"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/eliothedeman/bangarang/src/github.com/eliothedeman/bangarang/event"
	"github.com/eliothedeman/newman"
	"github.com/eliothedeman/randutil"
)

var (
	numEvents = 0
)

func newTestTcpClient() (*TcpClient, *newman.Conn) {
	// create a listener
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err.Error())
	}

	done := make(chan struct{})

	var rwc io.ReadWriteCloser
	go func() {
		conn, err := l.Accept()
		if err != nil {
			panic(err.Error())
		}

		rwc = conn
		done <- struct{}{}
	}()
	c, err := NewTcpClient(l.Addr().String())
	if err != nil {
		panic(err.Error())
	}

	<-done

	return c, newman.NewConn(rwc)
}

func newTestEvent() *event.Event {
	numEvents += 1
	e := event.NewEvent()
	e.Time = time.Now()
	e.Metric = rand.Float64()
	e.Tags.Set(randutil.AlphaString(10), randutil.AlphaString(10))

	return e
}

func TestSendTcp(t *testing.T) {
	c, r := newTestTcpClient()
	e := newTestEvent()

	err := c.Send(e)
	if err != nil {
		t.Fatal(err)
	}

	ne := event.NewEvent()

	err = r.Next(ne)
	if err != nil {
		t.Fatal(err)
	}

	if e.Metric != ne.Metric {
		t.Fail()
	}

	if e.Tags.String() != ne.Tags.String() {
		t.Fail()
	}

	if e.Time != ne.Time {
		t.Fail()
	}

	t.Log(*e, *ne)
}

func BenchmarkSendTcp(b *testing.B) {
	c, r := newTestTcpClient()
	go func() {
		g, _ := r.Generate(func() newman.Message {
			return event.NewEvent()
		})

		for {
			<-g
		}
	}()

	e := event.NewEvent()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c.Send(e)
	}
}
