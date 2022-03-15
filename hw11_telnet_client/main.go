package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"
)

var (
	myTimeout time.Duration
	address   string
)

func telnetRead(ctx context.Context, client TelnetClient, stop context.CancelFunc, out *bytes.Buffer) {
OUTER:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context done on read.", ctx.Err())
			stop()
			break OUTER
		default:
			myErr := client.Receive()
			if myErr != nil {
				stop()
			}
			fmt.Print(out.String())
		}
	}
}

func telnetWrite(ctx context.Context, client TelnetClient, stop context.CancelFunc, in *bytes.Buffer) {
	scanner := bufio.NewScanner(os.Stdin)
OUTER:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context done on write. ", ctx.Err())
			stop()
			break OUTER
		default:
			scanner.Scan()
			in.WriteString(scanner.Text() + "\n")
			myErr := client.Send()
			if myErr != nil {
				stop()
			}
		}
	}
}

func main() {
	flag.DurationVar(&myTimeout, "timeout", time.Second*10, "time out")
	flag.Parse()
	myArgs := os.Args
	if len(myArgs) < 3 {
		fmt.Println("not enough input params")
		return
	}
	if strings.Contains(myArgs[1], "--timeout") {
		address = net.JoinHostPort(myArgs[2], myArgs[3])
	} else {
		address = net.JoinHostPort(myArgs[1], myArgs[2])
	}

	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	client := NewTelnetClient(address, myTimeout, ioutil.NopCloser(in), out)
	client.Connect()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	go telnetRead(ctx, client, stop, out)
	go telnetWrite(ctx, client, stop, in)
	<-ctx.Done()
	client.Close()
	fmt.Println("finish")
}
