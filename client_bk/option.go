package client1

import (
	"crypto/tls"
	"time"
)

//Option ...
type Option struct {
	Retries        int
	TLSConfig      *tls.Config
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration

	Heartbeat         bool
	HeartbeatInterval time.Duration
}

// DefaultOption is a common option configuration for client1.
var DefaultOption = Option{
	Retries:        3,
	ConnectTimeout: 10 * time.Second,
}