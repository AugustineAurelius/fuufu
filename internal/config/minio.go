package config

type Minio struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	Token           string
}

func (man Manager) LoadMinio() Minio {
	return man.LoadConfig().MinioConfig
}

func (m Minio) GetEndpoint() string {
	return m.Endpoint
}

func (m Minio) GetAccessKeyID() string {
	return m.AccessKeyID
}

func (m Minio) GetSecretAccessKey() string {
	return m.SecretAccessKey
}

func (m Minio) GetToken() string {
	return m.Token
}
