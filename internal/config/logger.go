package config

type Logger struct {
	Debug bool `yaml:"debug"`
	Json  bool `yaml:"json"`
}

func (man Manager) LoadLogging() Logger {
	return man.LoadConfig().LoggerConfig
}
