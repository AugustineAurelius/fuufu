package config

type Server struct {
	Addr                  string
	AuthMiddlewareExclude []string
	GeoMiddlewareExclude  []string
}

func (man Manager) LoadServer() Server {
	return man.LoadConfig().SeverConfig
}
