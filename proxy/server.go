package proxy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/visola/go-proxy/configuration"
	myhttp "github.com/visola/go-proxy/http"
	"github.com/visola/go-proxy/mapping"
	"github.com/visola/go-proxy/statistics"
	"github.com/visola/go-proxy/variables"
)

// StartProxyServer starts the proxy server
func StartProxyServer() error {
	certFile := configuration.GetEnvironment().CertificateFile
	keyFile := configuration.GetEnvironment().KeyFile
	proxyPort := configuration.GetEnvironment().ProxyPort

	proxyServer := http.NewServeMux()
	proxyServer.HandleFunc("/", requestHandler)

	if certFile == "" || keyFile == "" {
		fmt.Printf("Starting proxy at: %s\n", GetURL())
		return http.ListenAndServe(fmt.Sprintf(":%d", proxyPort), proxyServer)
	}

	isSSL = true
	fmt.Printf("Starting proxy at: %s\n", GetURL())
	return http.ListenAndServeTLS(fmt.Sprintf(":%d", proxyPort), certFile, keyFile, proxyServer)
}

func requestHandler(w http.ResponseWriter, req *http.Request) {
	proxiedRequest := initializeHandling(req)

	mappings, mappingError := loadMappingsWithReplacedVariables()
	if mappingError != nil {
		finalizeWithError(req, w, proxiedRequest, mappingError)
		return
	}

	match := matchMapping(req, mappings)

	if match == nil {
		finalizeWithNotFound(req, w, proxiedRequest, "")
		return
	}

	response, handlingError := handleRequest(req, w, match)

	if handlingError == os.ErrNotExist {
		finalizeWithNotFound(req, w, proxiedRequest, response.executedURL)
		return
	}

	if handlingError != nil {
		finalizeWithError(req, w, proxiedRequest, handlingError)
		return
	}

	for name, value := range match.Mapping.Inject.Headers {
		proxiedRequest.RequestData.Headers[name] = []string{value}
	}

	proxiedRequest.ResponseData = statistics.HTTPData{
		Body:    response.body,
		Headers: response.headers,
	}
	proxiedRequest.ExecutedURL = response.executedURL
	proxiedRequest.ResponseCode = response.responseCode
	proxiedRequest.EndTime = getTime()
	statistics.AddProxiedRequest(proxiedRequest)
}

func finalizeWithError(req *http.Request, w http.ResponseWriter, proxiedRequest statistics.ProxiedRequest, err error) {
	myhttp.InternalError(req, w, err)
	proxiedRequest.EndTime = getTime()
	proxiedRequest.Error = err.Error()
	proxiedRequest.ResponseCode = http.StatusInternalServerError
	statistics.AddProxiedRequest(proxiedRequest)
	return
}

func finalizeWithNotFound(req *http.Request, w http.ResponseWriter, proxiedRequest statistics.ProxiedRequest, executedURL string) {
	myhttp.NotFound(req, w, executedURL)

	proxiedRequest.EndTime = getTime()
	proxiedRequest.ExecutedURL = executedURL
	proxiedRequest.ResponseCode = http.StatusNotFound
	statistics.AddProxiedRequest(proxiedRequest)
}

func getData(req *http.Request) (statistics.HTTPData, error) {
	result := statistics.HTTPData{}
	result.Headers = req.Header

	defer req.Body.Close()
	body, bodyErr := ioutil.ReadAll(req.Body)
	if bodyErr != nil {
		return result, bodyErr
	}

	// Rewind the body
	req.Body = closeableByteBuffer{bytes.NewBuffer(body)}

	if len(body) != 0 {
		if myhttp.IsText(req.Header["Content-Type"]...) {
			result.Body = string(body)
		} else {
			result.Body = "Binary"
		}
	}

	return result, nil
}

func getTime() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func handleRequest(req *http.Request, w http.ResponseWriter, match *mapping.MatchResult) (*proxyResponse, error) {
	if match.Mapping.Proxy {
		return proxyRequest(req, w, match)
	}

	return serveStaticFile(req, w, match)
}

func initializeHandling(req *http.Request) statistics.ProxiedRequest {
	reqData, _ := getData(req)

	return statistics.ProxiedRequest{
		Method:       req.Method,
		RequestedURL: req.URL.String(),
		RequestData:  reqData,
		StartTime:    getTime(),
	}
}

func loadMappingsWithReplacedVariables() ([]mapping.Mapping, error) {
	loadedMappings, mappingError := mapping.GetMappings()

	if mappingError != nil {
		return nil, mappingError
	}

	variables, varError := variables.GetSelectedValues()
	if varError != nil {
		return nil, varError
	}

	mappings := make([]mapping.Mapping, len(loadedMappings))
	for i, m := range loadedMappings {
		mappings[i] = m.WithReplacedVariables(variables)
		// Not all variables were replaced
		vars := mappings[i].GetVariables()
		if len(vars) > 0 {
			names := make([]string, len(vars))
			for i, v := range vars {
				names[i] = v.Name
			}
			return nil, fmt.Errorf("Mapping doesn't have all variables replaced. Missing variables: %s", names)
		}
	}

	return mappings, nil
}
