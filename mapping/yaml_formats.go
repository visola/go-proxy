package mapping

type yamlMapping struct {
	Static []mapping
	Proxy  []mapping
}

type mapping struct {
	From   string
	Regexp string
	To     string
}
