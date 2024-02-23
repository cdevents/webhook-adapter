# CDEvents Webhook Adapter

A CDEvents Webhook Adapter can receive events over Webhook and translate them to CDEvents, 
with the help of various translators which supports publishing events over Webhooks.

This adapter provides client and server implementation over RPC connection and exposes the interface to implement the plugins,  
Various translators can be implemented as plugins using Hashicorp's [go-plugin](https://github.com/hashicorp/go-plugin/)
