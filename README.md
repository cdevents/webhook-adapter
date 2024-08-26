# CDEvents Webhook Adapter

A CDEvents Webhook Adapter can receive events via Webhooks and translate them into CDEvents, 
with the help of various translators which supports publishing events via Webhooks.

This adapter provides client and server implementation over an RPC connection and exposes an interface for plugin implementation,
````go
// EventTranslator is the interface that we're exposing as a plugin.
type EventTranslator interface {
	TranslateEvent(event string, headers http.Header) (string, error)
}

````
Various translators can be developed as plugins by implementing this interface method and serve using Hashicorp's [go-plugin](https://github.com/hashicorp/go-plugin/?tab=readme-ov-file#usage)

## Plugins implemented using this adapter
| Plugin Name  | Plugin URL  |
| :------------ |:-------------------|
| gerrit-translator-cdevents| https://github.com/cdevents/gerrit-translator |
|  jira-translator-cdevents   | https://github.com/cdevents/jira-translator    |