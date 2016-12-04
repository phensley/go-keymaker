package keymaker

import (
	"crypto/tls"
	"encoding/pem"
	"fmt"
	"io"
	"net"
	"runtime"

	"regexp"

	"github.com/valyala/gorpc"
)

// Drone defines an RPC service that generates keys
type Drone struct {
	Config  *DroneConfig
	service *gorpc.Server

	// Pattern to match against client subject's CommonName
	regexpClientCN *regexp.Regexp
}

// NewDrone creates a Drone service
func NewDrone(config *DroneConfig) (*Drone, error) {
	d := &Drone{
		Config: config,
	}

	cc := config.Concurrency
	if cc <= 0 {
		cc = runtime.NumCPU()
	}

	if config.ClientCNRegexp != "" {
		r, err := regexp.Compile(config.ClientCNRegexp)
		if err != nil {
			return nil, err
		}
		d.regexpClientCN = r
	}

	tlsConfig, err := LoadDroneTLSConfig(config.Dir, config)
	if err != nil {
		return nil, err
	}

	d.service = &gorpc.Server{
		Addr:      config.Address,
		Handler:   d.handle,
		OnConnect: d.onConnect,
		Listener: &listener{
			F: func(addr string) (net.Listener, error) {
				return tls.Listen("tcp", config.Address, tlsConfig)
			},
		},
	}

	d.service.Concurrency = cc
	return d, nil
}

// Start up the drone service
func (d *Drone) Start() error {
	return d.service.Serve()
}

func (d *Drone) onConnect(remoteAddr string, rwc io.ReadWriteCloser) (io.ReadWriteCloser, error) {
	conn, ok := rwc.(*tls.Conn)
	if !ok {
		return nil, fmt.Errorf("Client %s did not connect over TLS", remoteAddr)
	}

	err := conn.Handshake()
	if err != nil {
		return nil, fmt.Errorf("Client %s TLS handshake failed: %s", remoteAddr, err)
	}

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return nil, fmt.Errorf("Client %s did not provide a certificate", remoteAddr)
	}

	cn := certs[0].Subject.CommonName
	if d.regexpClientCN == nil || d.regexpClientCN.MatchString(cn) {
		return rwc, nil
	}

	return nil, fmt.Errorf("Rejected client certificate as CN '%s' does not match regexp %s",
		cn, d.Config.ClientCNRegexp)
}

// Receives a single RPC request and generates a response
func (d *Drone) handle(addr string, req interface{}) interface{} {
	if keyType, ok := req.(string); ok {
		payload, err := generateKey(keyType)
		if err == nil {
			return &KeyResponse{payload, ErrOK, ""}
		}
		return &KeyResponse{nil, ErrKeyGen, err.Error()}
	}
	return &KeyResponse{nil, ErrBadRequest, fmt.Sprintf("bad request: %v", req)}
}

// Generates a private key and returns the PEM-encoded PKCS#8
func generateKey(keyType string) ([]byte, error) {
	key, err := GeneratePrivateKey(keyType)
	if err != nil {
		return nil, err
	}
	data, err := MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: data,
	}), nil
}

type listener struct {
	F func(addr string) (net.Listener, error)
	L net.Listener
}

func (n *listener) Init(addr string) (err error) {
	n.L, err = n.F(addr)
	return
}

func (n *listener) ListenAddr() net.Addr {
	if n.L != nil {
		return n.L.Addr()
	}
	return nil
}

func (n *listener) Accept() (conn io.ReadWriteCloser, clientAddr string, err error) {
	c, err := n.L.Accept()
	if err != nil {
		return nil, "", err
	}

	return c, c.RemoteAddr().String(), nil
}

func (n *listener) Close() error {
	return n.L.Close()
}
