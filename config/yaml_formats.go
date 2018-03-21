package config

type proxyConfig struct {
	Static []mapping
	Proxy  []mapping
}

type mapping struct {
	From   string
	Regexp string
	To     string
}
