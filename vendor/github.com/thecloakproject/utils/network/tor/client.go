// Steve Phillips / elimisteve
// 2012.11.04

package tor

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	// vars instead of consts so they can be change from the calling
	// program. TODO: Consider making these parameters of
	// `NewProxiedRequest`
	PROXY_URL  = "http://localhost:8118"
	USER_AGENT = "Mozilla/5.0 (Windows NT 6.1; rv:10.0) Gecko/20100101 Firefox/10.0"
)

// NewProxiedRequest does an HTTP request using the given method to
// remoteAddr, using r as the payload, over the local Tor
// connection. Assumes Polipo (Tor proxy) is running on
// http://localhost:8118.
func NewProxiedRequest(method, remoteAddr string, r io.Reader) (respBody []byte, err error) {
	proxyURL, err := url.Parse(PROXY_URL)
	if err != nil {
		return nil, fmt.Errorf("Error parsing proxy URL: %v", err)
	}

	transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	client := &http.Client{Transport: transport}

	// Create request
	request, err := http.NewRequest(method, remoteAddr, r)
	if err != nil {
		return nil, fmt.Errorf("Error creating %s request: %v", method, err)
	}

	// Change User Agent string (default: "Go http package")
	request.Header["User-Agent"] = []string{USER_AGENT}

	// Make HTTP request via proxy
	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Error performing %s request: %v", method, err)
	}

	// Read response body
	// TODO: Consider not doing this, or create a different function
	// that doesn't, so that (i) this function is a replacement for
	// 'http.NewRequest()' and (ii) the caller isn't prevented from
	// streaming the body instead of grabbing it all at once, which
	// will be problematic for large downloads.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %v", err)
	}
	defer resp.Body.Close()

	return body, nil
}
