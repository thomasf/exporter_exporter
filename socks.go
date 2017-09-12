package main

import (
	"context"
	"net"

	socks5 "github.com/armon/go-socks5"
	"github.com/davecgh/go-spew/spew"
	"github.com/golang/glog"
	"github.com/valyala/fasthttp/fasthttputil"
)

// Resolver .
type Resolver struct {
}

func (d Resolver) Resolve(ctx context.Context, name string) (context.Context, net.IP, error) {
	glog.Warningln("resolve", name)
	return ctx, net.IPv4(0, 0, 0, 0), nil
}

//  .
type Conn struct {
	net.Conn
	WhatToGet string
	ready     chan bool
}

func NewSL() *socksListener {
	return &socksListener{
		connC: make(chan *Conn, 1),
	}
}

// socksListener .
type socksListener struct {
	connC chan *Conn
}

// Accept waits for and returns the next connection to the listener.
func (s socksListener) Accept() (net.Conn, error) {

	c := <-s.connC
	glog.Infoln("accept")
	spew.Dump(c)
	c.ready <- true
	return c, nil

}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (s socksListener) Close() error {
	return nil
}

// Addr returns the listener's network address.
func (s socksListener) Addr() net.Addr {
	return &net.IPAddr{net.IPv4zero, ""}
}

// Addr returns the listener's network address.
func (s socksListener) HandleConn(ctx context.Context, r *socks5.Request, conn net.Conn) error {
	ready := make(chan bool, 1)

	cr := &Conn{
		Conn:      conn,
		ready:     ready,
		WhatToGet: r.DestAddr.String(),
	}
	s.connC <- cr

	<-ready
	spew.Dump(r.DestAddr)
	spew.Dump(r.RemoteAddr)
	return nil
}

func socksp(sl *fasthttputil.InmemoryListener) {
	conf := &socks5.Config{
		Resolver: Resolver{},
		Dial: func(ctx context.Context, net_, addr string) (net.Conn, error) {
			return sl.Dial()
		},
	}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}
	if err := server.ListenAndServe("tcp", ":3000"); err != nil {
		panic(err)
	}
}
