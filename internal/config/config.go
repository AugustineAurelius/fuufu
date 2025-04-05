package config

type Fuufu struct {
	PostgresConfig  Postgres
	PostgresSlave   Postgres
	LoggerConfig    Logger
	CollectorConfig Collector
	MinioConfig     Minio
	ShieldConfig    Shield
	SeverConfig     Server
	YukiConfig      Yuki
}
