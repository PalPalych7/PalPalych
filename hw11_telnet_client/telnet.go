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
	forSend, err := io.ReadAll(myT.In)
	if err != nil {
		return err
	}
	_, err = myT.Conn.Write(forSend)
	return err
}

func (myT *myTelnet) Receive() error {
	request := make([]byte, 1024)
	n, err := myT.Conn.Read(request)
	if err != nil {
		return err
	}
	_, err = myT.Out.Write(request[:n])
	return err
}

func (myT *myTelnet) Close() error {
	err := myT.Conn.Close()
	return err
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient /*myTelnet*/ {
	myTelnetClient := new(myTelnet)
	myTelnetClient.Adr = address
	myTelnetClient.Timeout = timeout
	myTelnetClient.In = in
	myTelnetClient.Out = out
	return TelnetClient(myTelnetClient)
}
