package config

type Plugin struct {
	Name   string
	Params interface{}
}

type Config struct {
	Plugins []Plugin
}

func ParseConfig(path string) *Config {
	return &Config{
		Plugins: []Plugin{
			{
				Name:   "mock",
				Params: map[string]string{},
			},
		},
	}
}
