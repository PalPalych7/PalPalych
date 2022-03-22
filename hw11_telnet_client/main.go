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

func telnetRead(ctx context.Context, client TelnetClient, stop context.CancelFunc /*, out io.Writer*/) {
OUTER:
	for {
		select {
		case <-ctx.Done():
			stop()
			break OUTER
		default:
			myErr := client.Receive()
			if myErr == nil {
				fmt.Fprint(os.Stderr, "...Connection was closed by peer")
			} else {
				fmt.Fprint(os.Stderr, "Reading error: ", myErr)
			}
			stop()
			// fmt.Print(out.String())
		}
	}
}

func telnetWrite(ctx context.Context, client TelnetClient, stop context.CancelFunc /*, in io.ReadCloser*/) {
	//	scanner := bufio.NewScanner(os.Stdin)
OUTER:
	for {
		select {
		case <-ctx.Done():
			stop()
			break OUTER
		default:
			//	scanner.Scan()
			//	in.WriteString(scanner.Text() + "\n")
			myErr := client.Send()
			if myErr == nil {
				fmt.Fprint(os.Stderr, "...EOF")
			} else {
				fmt.Fprint(os.Stderr, "Writng error: ", myErr)
			}
			stop()
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
	// in := &bytes.Buffer{}
	// out := &bytes.Buffer{}
	// client := NewTelnetClient(address, myTimeout, ioutil.NopCloser(in), out)
	in := os.Stdin
	out := os.Stdout
	client := NewTelnetClient(address, myTimeout, in, out)
	client.Connect()
	defer client.Close()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	go telnetRead(ctx, client, stop /*, out*/)
	go telnetWrite(ctx, client, stop /*, in*/)
	<-ctx.Done()
	//	fmt.Println("finish")
}
