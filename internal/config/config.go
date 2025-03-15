package config

type Fuufu struct {
	PostgresConfig  Postgres
	PostgresSlave   Postgres
	LoggerConfig    Logger
	CollectorConfig Collector
}
