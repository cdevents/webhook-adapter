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
	"net/http"
)

type GRPCClient struct {
	client proto.EventTranslatorClient
}

func (m *GRPCClient) TranslateEvent(event string, header http.Header) (string, error) {
	var entries []*proto.HeaderEntry
	var headerValues []*proto.HeaderValue
	for key, values := range header {
		for _, value := range values {
			headerValue := &proto.HeaderValue{Value: value}
			headerValues = append(headerValues, headerValue)
		}
		entry := &proto.HeaderEntry{
			Key:    key,
			Values: headerValues,
		}
		entries = append(entries, entry)
	}
	resp, err := m.client.TranslateEvent(context.Background(), &proto.TranslateEventRequest{
		Event:   event,
		Headers: entries,
	})
	if err != nil {
		return "", err
	}

	return resp.Event, nil
}

type GRPCServer struct {
	Impl EventTranslator
}

func (m *GRPCServer) TranslateEvent(ctx context.Context, req *proto.TranslateEventRequest) (*proto.TranslateEventResponse, error) {
	header := make(http.Header)
	for _, entry := range req.Headers {
		for _, value := range entry.Values {
			header.Add(entry.Key, value.Value)
		}
	}
	cdEvent, err := m.Impl.TranslateEvent(req.Event, header)
	return &proto.TranslateEventResponse{Event: cdEvent}, err
}
