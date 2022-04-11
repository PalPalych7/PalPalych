package main

import (
	"context"
	"flag"
	"fmt"
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

const minParms = 3

func telnetRead(ctx context.Context, client TelnetClient, stop context.CancelFunc) {
	for {
		select {
		case <-ctx.Done():
			stop()
			return
		default:
			if myErr := client.Receive(); myErr == nil {
				fmt.Fprint(os.Stderr, "...Connection was closed by peer")
			} else {
				fmt.Fprint(os.Stderr, "Reading error: ", myErr)
			}
			stop()
		}
	}
}

func telnetWrite(ctx context.Context, client TelnetClient, stop context.CancelFunc) {
	for {
		select {
		case <-ctx.Done():
			stop()
			return
		default:
			out := "...EOF"
			if myErr := client.Send(); myErr == nil {
				out = "Writing error: " + myErr.Error()
			}
			fmt.Fprint(os.Stderr, out)
			stop()
		}
	}
}

func main() {
	flag.DurationVar(&myTimeout, "timeout", time.Second*10, "time out")
	flag.Parse()
	myArgs := os.Args
	if len(myArgs) < minParms {
		fmt.Println("not enough input params")
		return
	}
	if strings.Contains(myArgs[1], "--timeout") {
		address = net.JoinHostPort(myArgs[2], myArgs[3])
	} else {
		address = net.JoinHostPort(myArgs[1], myArgs[2])
	}
	in := os.Stdin
	out := os.Stdout
	client := NewTelnetClient(address, myTimeout, in, out)
	client.Connect()
	defer client.Close()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	go telnetRead(ctx, client, stop)
	go telnetWrite(ctx, client, stop)
	<-ctx.Done()
}
