package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	rabbithole "github.com/michaelklishin/rabbit-hole"
	"github.com/reconquest/karma-go"
)

func parseDSN(rawDSN string) string {
	DSN := strings.TrimSpace(rawDSN)

	if !strings.HasPrefix(DSN, "http://") &&
		!strings.HasPrefix(DSN, "https://") {

		return fmt.Sprintf("http://%s", DSN)
	}

	return DSN
}

func makeRabbitMQClient(
	dsn string,
	username string,
	password string,
	caPath string,
) (*rabbithole.Client, error) {
	destiny := karma.Describe(
		"method", "makeRabbitMQClient",
	)

	var (
		rmqc *rabbithole.Client
		err  error
	)

	if !strings.HasPrefix(
		dsn, "https://",
	) {
		rmqc, err = rabbithole.NewClient(
			dsn,
			username,
			password,
		)
		if err != nil {
			return nil, destiny.Describe(
				"RabbitMQ DSN", dsn,
			).Describe(
				"error", err,
			).Reason(
				"can't create RabbitMQ client",
			)
		}

		return rmqc, nil
	}

	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		return nil, destiny.Describe(
			"error", err,
		).Reason(
			"can't obtain root CA pool",
		)
	}

	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	if caPath != noneValue {
		cert, err := ioutil.ReadFile(caPath)
		if err != nil {
			return nil, destiny.Describe(
				"CA path", caPath,
			).Describe(
				"error", err,
			).Reason(
				"can't read CA",
			)
		}

		if hasAppended := rootCAs.AppendCertsFromPEM(cert); !hasAppended {
			return nil, destiny.Reason("cert CA has't appended")
		}
	}

	tlsConfig := &tls.Config{
		RootCAs: rootCAs,
	}

	rmqc, err = rabbithole.NewTLSClient(
		dsn,
		username,
		password,
		&http.Transport{TLSClientConfig: tlsConfig},
	)

	if err != nil {

		return nil, destiny.Describe(
			"RabbitMQ DSN", dsn,
		).Describe(
			"error", err,
		).Reason(
			"can't create RabbitMQ TLS client",
		)

	}

	return rmqc, nil
}
