package config

type Rconfig struct {
	RedisAddress string
	RedisPort    int
	RedisAuth    string
}

type SensuConfig struct {
	Hostname string
	User     string
	Token    string
	Port     int
	Enabled  bool
}

type SlackConfig struct {
	Token         string
	Enabled       bool
	Channel       string
	AuthorName    string
	AuthorSubname string
	TriggerOn     []string
}

type Evconfig struct {
	RedisAddress string
	RedisPort    int
	RedisAuth    string
	Sensu        SensuConfig
	Slack        SlackConfig
}
