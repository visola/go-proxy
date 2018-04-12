package mapping

type yamlMapping struct {
	Static []mapping
	Proxy  []mapping
}

type mapping struct {
	From   string
	Inject Injection
	Regexp string
	To     string
}
