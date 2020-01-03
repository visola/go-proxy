package upstream

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// ProxyEndpoint proxy requests to another HTTP/S server
type ProxyEndpoint struct {
	BaseEndpoint
	Headers map[string][]string
	To      string `json:"to"`
}

// Handle handles requests proxying the content to another HTTP server
func (p *ProxyEndpoint) Handle(req *http.Request, resp http.ResponseWriter) HandleResult {
	var newURL *url.URL
	var parseErr error

	if p.Regexp != "" {
		newURL, parseErr = url.Parse(replaceRegexp(req.URL.Path, p.To, p.ensureRegexp()))
		if parseErr != nil {
			return internalServerError(p.UpstreamName+":[error]", req, resp, parseErr)
		}
	} else {
		newURL, parseErr = url.Parse(p.To)
		if parseErr != nil {
			return internalServerError(p.UpstreamName+":[error]", req, resp, parseErr)
		}
		newURL.Path = newURL.Path + "/" + req.URL.Path[len(p.From):]
	}

	return proxyHandleResult(p, newURL, req, resp)
}

func createHTTPClient() *http.Client {
	return &http.Client{
		// Do not auto-follow redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

func proxyHandleResult(p *ProxyEndpoint, newURL *url.URL, req *http.Request, resp http.ResponseWriter) HandleResult {
	executedURL := p.UpstreamName + ":" + p.To

	defer req.Body.Close()
	bodyInBytes, readBodyErr := ioutil.ReadAll(req.Body)
	if readBodyErr != nil {
		return internalServerError(executedURL, req, resp, readBodyErr)
	}

	// Copy parameters
	newQuery := newURL.Query()
	for name, values := range req.URL.Query() {
		for _, value := range values {
			newQuery.Add(name, value)
		}
	}
	newURL.RawQuery = newQuery.Encode()

	proxiedReq, proxiedReqErr := http.NewRequest(req.Method, newURL.String(), bytes.NewBuffer(bodyInBytes))
	if proxiedReqErr != nil {
		return internalServerError(executedURL, req, resp, proxiedReqErr)
	}

	// Copy request headers
	for name, values := range req.Header {
		for _, value := range values {
			proxiedReq.Header.Add(name, value)
		}
	}

	// Inject headers
	for name, values := range p.Headers {
		for _, value := range values {
			proxiedReq.Header.Add(name, value)
		}
	}

	proxiedResp, respErr := createHTTPClient().Do(proxiedReq)
	if respErr != nil {
		return internalServerError(executedURL, req, resp, respErr)
	}

	// Copy response headers
	for name, values := range proxiedResp.Header {
		for _, value := range values {

			// Fix location headers to point to proxy
			if strings.ToLower(name) == "location" {
				if strings.HasPrefix(value, p.To) {
					value = req.URL.Host + value[len(p.To):]
				}
			}
			resp.Header().Add(name, value)
		}
	}

	// Copy status
	resp.WriteHeader(proxiedResp.StatusCode)

	handleResult := handleReadCloser(proxiedResp.Body, executedURL, req, resp)
	handleResult.ResponseCode = proxiedResp.StatusCode
	handleResult.ResponseHeaders = resp.Header()

	return handleResult
}
