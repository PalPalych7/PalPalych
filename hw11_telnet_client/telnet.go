package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type myTelnet struct {
	Adr     string
	Timeout time.Duration
	Conn    net.Conn
	In      io.ReadCloser
	Out     io.Writer
}

func (myT *myTelnet) Connect() error {
	myCon, err := net.DialTimeout("tcp", myT.Adr, myT.Timeout)
	myT.Conn = myCon
	return err
}

func (myT *myTelnet) Send() error {
	_, err := io.Copy(myT.Conn, myT.In)
	return err
}

func (myT *myTelnet) Receive() error {
	_, err := io.Copy(myT.Out, myT.Conn)
	return err
}

func (myT *myTelnet) Close() error {
	err := myT.Conn.Close()
	return err
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient /*myTelnet*/ {
	myTelnetClient := &myTelnet{
		Adr:     address,
		In:      in,
		Timeout: timeout,
		Out:     out,
	}
	return TelnetClient(myTelnetClient)
}
