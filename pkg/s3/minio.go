package s3

import (
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Configurator interface {
	GetEndpoint() string
	GetAccessKeyID() string
	GetSecretAccessKey() string
	GetToken() string
}

func NewMinio(config Configurator) (*minio.Client, error) {
	client, err := minio.New(config.GetEndpoint(), &minio.Options{
		Creds: credentials.NewStaticV4(config.GetAccessKeyID(), config.GetSecretAccessKey(), config.GetToken()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start minio connection: %w", err)
	}

	return client, nil
}
