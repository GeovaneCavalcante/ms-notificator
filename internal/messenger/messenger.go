package notification

import (
	"context"
	"crypto/tls"
	"errors"
)

var (
	TimeoutError = errors.New("Error timeout during broker request")
)

type Request interface {
	Send(ctx context.Context, configuration Configuration, r interface{}) (Response, error)
	ConvertCertificateToPEM(clientCert, certPassword string) (tls.Certificate, error)
}

type Configuration struct {
	Topic   string
	Message interface{}
}

type Response struct {
	Status  int
	Body    []byte
	Headers []byte
}
