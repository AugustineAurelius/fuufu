package config

type Yuki struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (man Manager) LoadYuki() Yuki {
	return man.LoadConfig().YukiConfig
}
