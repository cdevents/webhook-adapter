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
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadTranslator_PluginsConfig(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	testdata := filepath.Join(wd, "testdata")
	if err := os.Chdir(testdata); err != nil {
		panic(err)
	}
	defer func(dir string) {
		err := os.Chdir(dir)
		if err != nil {
			panic(err)
		}
	}(wd)
	translator, err := LoadConfig("translator-plugins-test.yaml")
	expectedPath := "./test_plugins_path"
	if translator.Translator.Path != expectedPath {
		t.Errorf("Translator path Expected %s, but got %s", expectedPath, translator.Translator.Path)
	}
	actualLen := len(translator.Translator.Plugins)
	if actualLen != 3 {
		t.Errorf("Translator Plugins count Expected 3, but got %d", actualLen)
	}
	assert.Nil(t, err, "Expected error nil")
}

func TestLoadTranslator_PluginsConfig_Invalid(t *testing.T) {
	translator, err := LoadConfig("translator-plugins-invalid.yaml")

	assert.Nil(t, translator, "Translator Plugins Config should be nil when there's an error")
	assert.Error(t, err, "Expected an error")

}

func TestValidate_ValidURL(t *testing.T) {
	expectedValidURL := "http://example.com/translator/version/plugin"
	actualUrl, err := ValidateURL(expectedValidURL)

	if actualUrl != expectedValidURL {
		t.Errorf("Failed to test with valid URL, %s", err)
	}
}

func TestValidate_InValidURL(t *testing.T) {
	inValidURL := "http://example.com/!@#$%^&*"
	expectedURL := ""
	actualUrl, err := ValidateURL(inValidURL)
	if actualUrl != expectedURL {
		t.Errorf("Failed to test with Invalid URL, %s", err)
	}
}

func TestSendCDEvent(t *testing.T) {
	eventToSend := "{\n  \"context\": {\n    \"version\": \"0.3.0\",\n    \"id\": \"eb175ff7-2fda-44c5-bdb4-17b1c7342fc0\",\n    \"source\": \"http://dev.cdevents\",\n    \"type\": \"dev.cdevents.pipelinerun.finished.0.1.1\",\n    \"timestamp\": \"2024-02-29T15:23:09Z\"\n  },\n  \"subject\": {\n    \"id\": \"/dev/pipeline/run/subject\",\n    \"source\": \"/dev/pipeline/run/subject\",\n    \"type\": \"pipelineRun\",\n    \"content\": {\n      \"pipelineName\": \"Name-pipeline\",\n      \"url\": \"http://dev/pipeline/url\",\n      \"outcome\": \"success\",\n      \"errors\": \"errors to place\"\n    }\n  },\n  \"customData\": {\n  },\n  \"customDataContentType\": \"application/json\"\n}"

	err := SendCDEvent(eventToSend, "http://cdevents.message.com/default/events-broker")
	assert.Error(t, err, "Expected an error")
	assert.Contains(t, err.Error(), "failed to send CDEvent", "Expected log message not found")
}
