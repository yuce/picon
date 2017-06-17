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
	"encoding/json"
	"fmt"
	"time"

	pj "github.com/hokaccha/go-prettyjson"
)

type promptInfo struct {
	address string
	index   string
}

func printError(err error) {
	fmt.Println(colorString(fgRed, err.Error()))
}

func printWarning(msg string) {
	fmt.Println(colorString(fgRed, msg))
}

func colorString(color Ansi, msg string) string {
	return fmt.Sprintf("%s%s%s", color, msg, attrReset)
}

func autoSessionName() string {
	return time.Now().Format("2006-01-02_15-04-05")
}

func tryPrettifyJSON(text []byte) []byte {
	jsonResponse := make(map[string]interface{})
	err := json.Unmarshal(text, &jsonResponse)
	if err != nil {
		return text
	}
	prettyText, err := pj.Marshal(jsonResponse)
	if err != nil {
		return text
	}
	return prettyText
}
