package config

import "strconv"

type Collector struct {
	Host string
	Port int
}

func (c Collector) Addr() string {
	return c.Host + ":" + strconv.Itoa(c.Port)
}

func (man Manager) LoadCollector() Collector {
	return man.LoadConfig().CollectorConfig
}
