package config

type Server struct {
	Addr string
}

func (man Manager) LoadServer() Server {
	return man.LoadConfig().SeverConfig
}
