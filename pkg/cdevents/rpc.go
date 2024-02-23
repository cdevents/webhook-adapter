package cdevents

import (
	"net/rpc"
)

// RPCClient is an implementation of EventTranslator that talks over RPC.
type RPCClient struct{ client *rpc.Client }

func (m *RPCClient) TranslateEvent(event string) (string, error) {
	var resp string
	err := m.client.Call("Plugin.TranslateEvent", event, &resp)
	return resp, err
}

// Here is the RPC server that RPCClient talks to, conforming to
// the requirements of net/rpc
type RPCServer struct {
	// This is the real implementation
	Impl EventTranslator
}

func (m *RPCServer) TranslateEvent(event string, resp *string) error {
	v, err := m.Impl.TranslateEvent(event)
	*resp = v
	return err
}
