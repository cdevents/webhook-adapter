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
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"

	cdevents "github.com/cdevents/sdk-go/pkg/api"
	cdevents04 "github.com/cdevents/sdk-go/pkg/api/v04"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"gopkg.in/yaml.v3"
)

type Plugin struct {
	Name          string `yaml:"name"`
	PluginURL     string `yaml:"pluginURL"`
	MessageBroker string `yaml:"messageBroker"`
}

type Translator struct {
	Path    string   `yaml:"path"`
	Plugins []Plugin `yaml:"plugins"`
}

type TranslatorPlugins struct {
	Translator Translator `yaml:"translator"`
}

func LoadConfig(fileName string) (*TranslatorPlugins, error) {
	var translator TranslatorPlugins

	file, err := os.ReadFile(fileName)
	if err != nil {
		log.Printf("Error Reading configuration file: %v", err)
		return nil, err
	}

	err = yaml.Unmarshal(file, &translator)
	if err != nil {
		log.Printf("Error Unmarshal configuration: %v", err)
		return nil, err
	}

	return &translator, nil
}

func ValidateURL(URL string) (string, error) {
	parsedURL, err := url.Parse(URL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return "", err
	}
	return parsedURL.String(), nil
}

func SendCDEvent(event string, messageBrokerURL string) error {
	fmt.Println("IN SendCDEvent with event " + event)
	cdEvent, err := cdevents04.NewFromJsonString(event)
	if err != nil {
		log.Printf("failed to create CDEvent from Json string, %v", err)
		return err
	}
	ce, err := cdevents.AsCloudEvent(cdEvent)
	if err != nil {
		log.Printf("failed to create CDEvent as CloudEvent, %v", err)
		return err
	}

	ctx := cloudevents.ContextWithTarget(context.Background(), messageBrokerURL)
	ctx = cloudevents.WithEncodingBinary(ctx)

	c, err := cloudevents.NewClientHTTP()
	if err != nil {
		log.Printf("failed to create client, %v", err)
		return err
	}
	result := c.Send(ctx, *ce)
	if cloudevents.IsNACK(result) || cloudevents.IsUndelivered(result) {
		log.Printf("Failed to send CDEvent, %v", result)
		return errors.New("failed to send CDEvent to target message-broker URL: " + messageBrokerURL)
	}

	log.Printf("Sent CDEvent to target message-broker URL: %v", result)
	return nil
}
