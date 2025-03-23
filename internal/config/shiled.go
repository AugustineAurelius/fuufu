package config

type Shield struct {
	Secret string
}

func (man Manager) LoadShield() Shield {
	return man.LoadConfig().ShieldConfig
}

func (s Shield) GetSecret() string {
	return s.Secret
}
