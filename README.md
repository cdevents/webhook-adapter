# CDEvents Webhook Adapter

A CDEvents Webhook Adapter can receive events via Webhooks and translate them into CDEvents, 
with the help of various translators which supports publishing events via Webhooks.

This adapter provides client and server implementation over an RPC connection and exposes an interface for plugin implementation,  
Various translators can be developed as plugins using Hashicorp's [go-plugin](https://github.com/hashicorp/go-plugin/)
