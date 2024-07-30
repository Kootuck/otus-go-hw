package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Parameters struct {
	Host    string
	Port    string
	Timeout time.Duration
}

var p Parameters

func main() {
	// 1. Programm arguments
	err := setParameters()
	if err != nil {
		log.Fatal("Error:", err)
	}

	// 2. Setting up client and tcp connection
	addr := net.JoinHostPort(p.Host, p.Port)
	telnet := NewTelnetClient(addr, p.Timeout, os.Stdin, os.Stdout)

	err = telnet.Connect()
	if err != nil {
		log.Fatal("Error:", err)
	}
	defer telnet.Close()

	// 3. Setting context:
	// a. Signals-responsive;
	// b. Cancel: make blocking i/o quit via connection deadline.
	ctxWCancel, cancel := context.WithCancel(context.Background())
	ctx, _ := signal.NotifyContext(ctxWCancel, syscall.SIGINT, syscall.SIGTERM)
	defer cancel() // all paths cancel the context to avoid context leak.

	telnet.SetDeadline(ctx)

	var wg sync.WaitGroup
	wg.Add(2)

	// 4. sending messages from client.
	go func() {
		defer func() {
			wg.Done()
			cancel()
		}()
		err := telnet.Send()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		}
	}()

	// 5. receiving messages.
	go func() {
		defer func() {
			wg.Done()
			cancel()
		}()
		err := telnet.Receive()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		}
	}()

	wg.Wait()
}

func setParameters() error {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout for connection closure if idle (e.g., 5s, 2m, 1h)")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		return fmt.Errorf("invalid arguments, usage: go-telnet [--timeout duration] <host> <port>")
	}

	p = Parameters{
		Host:    args[0],
		Port:    args[1],
		Timeout: timeout,
	}

	return nil
}
