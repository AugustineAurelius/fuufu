package config

type Logger struct {
	Debug bool `yaml:"debug"`
	JSON  bool `yaml:"json"`
}

func (man Manager) LoadLogging() Logger {
	return man.LoadConfig().LoggerConfig
}
