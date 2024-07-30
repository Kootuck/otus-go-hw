package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

var ErrFailedToConnect = errors.New("failed to establish a connection")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
	SetDeadline(context.Context)
}
type TelnetClientImpl struct {
	address string
	timeout time.Duration
	conn    net.Conn
	in      io.Reader
	out     io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.Reader, out io.Writer) TelnetClient {
	return &TelnetClientImpl{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (t *TelnetClientImpl) Connect() error {
	var err error
	t.conn, err = net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		wrapped := fmt.Errorf("net: %w", err)
		return fmt.Errorf("%w: %v", ErrFailedToConnect, wrapped) //nolint:errorlint
	}
	return nil
}

func (t *TelnetClientImpl) Close() error {
	if t.conn != nil {
		err := t.conn.Close()
		if err != nil {
			return fmt.Errorf("failed to close connection: %w", err)
		}
		fmt.Println("Connection closed (client)")
	}
	return nil
}

func (t *TelnetClientImpl) Send() error {
	_, err := io.Copy(t.conn, t.in)
	if err != nil {
		return fmt.Errorf("error copying from stdin to conn: %w", err)
	}
	return nil
}

func (t *TelnetClientImpl) Receive() error {
	_, err := io.Copy(t.out, t.conn)
	if err != nil {
		return fmt.Errorf("error copying from conn to stdout: %w", err)
	}
	return nil
}

func (t *TelnetClientImpl) SetDeadline(ctx context.Context) {
	go func() {
		<-ctx.Done()
		if t.conn != nil {
			t.conn.SetDeadline(time.Now()) // Forces io.Copy to return.
		}
	}()
}
