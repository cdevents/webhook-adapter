# CDEvents Webhook Adapter

A CDEvents Webhook Adapter can receive events via Webhooks and translate them into CDEvents, 
with the help of various translators which supports publishing events via Webhooks.
The translated CDEvents can be sent to configured `messageBroker` URL 

This adapter provides client and server implementation over an RPC connection and exposes an interface for plugin implementation,
````go
// EventTranslator is the interface that we're exposing as a plugin.
type EventTranslator interface {
	TranslateEvent(event string, headers http.Header) (string, error)
}

````
Various translators can be developed as plugins by implementing this interface method and serve using Hashicorp's [go-plugin](https://github.com/hashicorp/go-plugin/?tab=readme-ov-file#usage)

## Translator's configuration
Once the plugin implemented [translator-plugins.yaml](./translator-plugins.yaml) needs an update with the details of Plugin name, pluginURL and messageBroker
````yaml
translator:
  path: "./plugins"
  plugins:
    - name: "<pluginName>"
      pluginURL: "<plugin's binary URL>"
      messageBroker: "<message broker URL>"
````

`DOWNLOAD_PLUGIN` is an environment variable that can be set to `true` or `false` to download the plugin's binary from the translator's `pluginURL` when launching this adapter. Alternatively, the plugin's binary can also be downloaded manually and placed under the `./plugins` directory (create it if it does not exist).

## Plugins implemented using this adapter
| Plugin Name  | Plugin Repo  |
| :------------ |:-------------------|
| gerrit-translator-cdevents| https://github.com/cdevents/gerrit-translator |
|  jira-translator-cdevents   | https://github.com/cdevents/jira-translator    |