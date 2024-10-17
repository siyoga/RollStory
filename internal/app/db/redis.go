package db

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
)

type RedisClient struct {
	DB *redis.Client
}

func NewRedisClient(DSN string, certLoc string) (*RedisClient, error) {
	var client *redis.Client
	if certLoc != "" {
		rootCertLoc := x509.NewCertPool()
		pem, err := os.ReadFile(certLoc)
		if err != nil {
			return nil, fmt.Errorf("error reading cert: %w", err)
		}

		rootCertLoc.AppendCertsFromPEM(pem)

		connCfg, err := redis.ParseURL(DSN)
		if err != nil {
			return nil, fmt.Errorf("error parsing redis url: %w", err)
		}
		connCfg.TLSConfig = &tls.Config{
			RootCAs:            rootCertLoc,
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS12,
		}

		client = redis.NewClient(connCfg)
	} else {
		connCfg, err := redis.ParseURL(DSN)
		if err != nil {
			return nil, fmt.Errorf("error parsing redis url: %w", err)
		}

		client = redis.NewClient(connCfg)
	}

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("error connecting to redis: %w", err)
	}

	return &RedisClient{DB: client}, nil
}

func (client *RedisClient) Close() error {
	if err := client.DB.Close(); err != nil {
		return err
	}

	return nil
}
