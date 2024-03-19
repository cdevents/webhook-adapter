/*
Copyright (C) 2024 Nordix Foundation.
For a full list of individual contributors, please see the commit history.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
SPDX-License-Identifier: Apache-2.0
*/

package cdevents

import (
	"net/http"
	"net/rpc"
)

// RPCClient is an implementation of EventTranslator that talks over RPC.
type RPCClient struct{ client *rpc.Client }

func (m *RPCClient) TranslateEvent(event string, headers http.Header) (string, error) {
	var resp string
	err := m.client.Call("Plugin.TranslateEvent", map[string]interface{}{
		"event":   event,
		"headers": headers,
	}, &resp)
	return resp, err
}

// Here is the RPC server that RPCClient talks to, conforming to
// the requirements of net/rpc
type RPCServer struct {
	// This is the real implementation
	Impl EventTranslator
}

func (m *RPCServer) TranslateEvent(event string, headers http.Header, resp *string) error {
	v, err := m.Impl.TranslateEvent(event, headers)
	*resp = v
	return err
}
