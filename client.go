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
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	pj "github.com/hokaccha/go-prettyjson"
	pilosa "github.com/pilosa/go-pilosa"
)

const ConnectTimeout = 10 * time.Second
const SocketTimeout = 100 * time.Second

type Client struct {
	URI        *pilosa.URI
	httpClient *http.Client
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
	jsonResponse := make(map[string]interface{})
	json.Unmarshal(response, &jsonResponse)
	return pj.Marshal(jsonResponse)
}

func (c *Client) httpRequest(method string, path string, data []byte) ([]byte, error) {
	path = c.URI.Normalize() + path
	request, err := http.NewRequest("POST", path, bytes.NewReader(data))
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
	return buf, err
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
