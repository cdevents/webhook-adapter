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

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/cdevents/webhook-adapter/pkg/cdevents"
	"github.com/hashicorp/go-plugin"
)

func run() error {
	translator, err := cdevents.LoadConfig("translator-plugins.yaml")
	if err != nil {
		log.Printf("Error loading translator plugins Config: %s\n", err)
		return err
	}
	log.Printf("Loaded translator plugins Config: %s\n", translator)

	port := 8080
	log.Printf("### Starting CDEvents Webhook Adapter ###\n")
	log.Printf("Server listening on :%d\n", port)
	for _, translatorPlugin := range translator.Translator.Plugins {
		pluginPath := path.Join(translator.Translator.Path, translatorPlugin.Name)
		err := downloadPlugin(translatorPlugin, pluginPath, translator)
		if err != nil {
			return err
		}
		endpoint := "/translate/" + translatorPlugin.Name
		http.HandleFunc(endpoint, makeHandler(translatorPlugin, pluginPath))
		log.Printf("Serving Translator endpoint for %s plugins : %s\n", translatorPlugin.Name, endpoint)
	}

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Printf("Error starting HTTP server: %v\n", err)
		return err
	}
	return nil
}

func downloadPlugin(translatorPlugin cdevents.Plugin, pluginPath string, translator *cdevents.TranslatorPlugins) error {
	if translatorPlugin.PluginURL != "" && strings.ToLower(os.Getenv("DOWNLOAD_PLUGIN")) == "true" {
		pluginURL, err := cdevents.ValidateURL(translatorPlugin.PluginURL)
		if err != nil {
			log.Printf("Error validating translator plugins URL : %s\n", pluginURL)
			return err
		}
		err = cdevents.Download(pluginURL, pluginPath)
		if err != nil {
			log.Printf("Error downloading translator plugins : %s\n", err)
			return err
		}

	} else if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		log.Printf("Plugin %s not downloaded under : %s\n", translatorPlugin.Name, translator.Translator.Path)
		log.Fatalf("Please download the plugins or update the pluginURL and set the env variable DOWNLOAD_PLUGIN=true for : %s\n", translatorPlugin.Name)
		return err
	}
	return nil
}

func makeHandler(translatorPlugin cdevents.Plugin, pluginPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		log.Println("Received event payload : " + string(body))

		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: cdevents.Handshake,
			Plugins:         cdevents.PluginMap,
			Cmd:             exec.Command("sh", "-c", pluginPath),
			AllowedProtocols: []plugin.Protocol{
				plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
		})
		defer client.Kill()

		rpcClient, err := client.Client()
		if err != nil {
			log.Printf("Error connecting RPC client: %v", err)
			http.Error(w, "Error connecting to RPC client", http.StatusInternalServerError)
			return
		}
		log.Printf("RPC client created for plugins %s\n", translatorPlugin.Name)

		raw, err := rpcClient.Dispense("translator_grpc")
		if err != nil {
			log.Printf("Error requesting the GRPC translator plugins: %v", err)
			http.Error(w, "Error connecting to RPC client", http.StatusInternalServerError)
			return
		}

		eventTranslator := raw.(cdevents.EventTranslator)

		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		event, err := eventTranslator.TranslateEvent(string(body), r.Header)
		if err != nil {
			http.Error(w, "Error translating event", http.StatusInternalServerError)
			return
		}
		log.Println("Event translated : " + event)
		err = cdevents.SendCDEvent(event, translatorPlugin.MessageBroker)
		if err != nil {
			http.Error(w, "Error sending CDEvent", http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error exit main(): %+v\n", err)
	}
}
