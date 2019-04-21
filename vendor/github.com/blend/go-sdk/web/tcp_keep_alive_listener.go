package web

import (
	"net"
	"time"
)

const (
	// DefaultTCPKeepAliveListenerPeriod is the default keep alive period for the tcp listener.
	DefaultTCPKeepAliveListenerPeriod = 3 * time.Minute
)

// TCPKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
// Specifically it allows us to bind to any available port (i.e. `BIND_ADDR=127.0.0.1:`
// and eventually determine that port from the listener.
type TCPKeepAliveListener struct {
	*net.TCPListener
}

// Accept accepts the connection.
func (ln TCPKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(DefaultTCPKeepAliveListenerPeriod)
	return tc, nil
}
