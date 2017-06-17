/*
Copyright 2017 Yuce Tekol

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

1. Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its
contributors may be used to endorse or promote products derived
from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND
CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR
CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH
DAMAGE.
*/

package picon

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	pilosa "github.com/pilosa/go-pilosa"
)

const ConnectTimeout = 10 * time.Second
const SocketTimeout = 100 * time.Second

type Client struct {
	URI        *pilosa.URI
	httpClient *http.Client
}

type HttpResponse struct {
	Body       []byte
	Type       string
	StatusCode int
}

func NewClient(addr string) (*Client, error) {
	uri, err := pilosa.NewURIFromAddress(addr)
	if err != nil {
		return nil, err
	}
	return &Client{
		URI:        uri,
		httpClient: newHTTPClient(),
	}, nil
}

func (c *Client) query(index string, text string) ([]byte, error) {
	path := "/index/" + index + "/query"
	response, err := c.httpRequest("POST", path, []byte(text))
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}

func (c *Client) serverVersion() (string, error) {
	response, err := c.httpGet("/version")
	if err != nil {
		return "", err
	}
	jsonResponse := make(map[string]interface{})
	json.Unmarshal(response.Body, &jsonResponse)
	if version, ok := jsonResponse["version"]; ok {
		return version.(string), nil
	}
	return "", errors.New("Version not detected")

}

func (c *Client) httpGet(path string) (*HttpResponse, error) {
	return c.httpRequest("GET", path, []byte{})
}

func (c *Client) httpRequest(method string, path string, data []byte) (*HttpResponse, error) {
	path = c.URI.Normalize() + path
	request, err := http.NewRequest(method, path, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("%s: %s", response.Status, buf)
	}
	return &HttpResponse{
		Body:       buf,
		Type:       response.Header.Get("content-type"),
		StatusCode: response.StatusCode,
	}, nil
}

func newHTTPClient() *http.Client {
	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: ConnectTimeout,
		}).Dial,
	}
	return &http.Client{
		Transport: transport,
		Timeout:   SocketTimeout,
	}
}
