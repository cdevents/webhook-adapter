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
	"context"
	"github.com/cdevents/webhook-adapter/pkg/proto"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"net/http"
	"net/rpc"
)

// Handshake is a common handshake that is shared by plugins and host.
var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	"translator_grpc": &TranslatorGRPCPlugin{},
	"translator":      &TranslatorPlugin{},
}

// EventTranslator is the interface that we're exposing as a plugins.
type EventTranslator interface {
	TranslateEvent(event string, headers http.Header) (string, error)
}

// TranslatorPlugin is the implementation of plugins.Plugin so we can serve/consume this.
type TranslatorPlugin struct {
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl EventTranslator
}

func (p *TranslatorPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &RPCServer{Impl: p.Impl}, nil
}

func (*TranslatorPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RPCClient{client: c}, nil
}

// TranslatorGRPCPlugin is the implementation of plugin.GRPCPlugin so we can serve/consume this.
type TranslatorGRPCPlugin struct {
	// GRPCPlugin must still implement the Plugin interface
	plugin.Plugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl EventTranslator
}

func (p *TranslatorGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterEventTranslatorServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *TranslatorGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewEventTranslatorClient(c)}, nil
}
