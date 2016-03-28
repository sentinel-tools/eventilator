package config

type Rconfig struct {
	RedisAddress string
	RedisPort    int
	RedisAuth    string
}

type SensuConfig struct {
	Hostname  string
	User      string
	Token     string
	Port      int
	Enabled   bool
	TriggerOn []string
}

type SlackConfig struct {
	Token         string
	Enabled       bool
	Channel       string
	AuthorName    string
	Username      string
	AuthorSubname string
	TriggerOn     []string
}

type Evconfig struct {
	RedisAddress string
	RedisPort    int
	RedisAuth    string
	RedisEnabled bool
	Sensu        SensuConfig
	SensuJIT     SensuJITConfig
	Slack        SlackConfig
}

type SensuJITConfig struct {
	IP                   string
	Port                 int
	Enabled              bool
	TriggerOn            []string
	HandlerNoGoodSlave   string
	HandlerPromotedSlave string
}
