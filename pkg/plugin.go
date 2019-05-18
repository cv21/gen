package pkg

import (
	"net/rpc"

	plugin "github.com/hashicorp/go-plugin"
)

// It is default command for each generator plugin.
const cmdPluginGenerate = "Plugin.Generate"

// DefaultHandshakeConfig useful for plugin compatibility specification.
var DefaultHandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "HOST",
	MagicCookieValue: "GEN",
}

// It is a plugin client.
type Client struct {
	client *rpc.Client
}

// Generate method implements Generator interface.
// It calls plugin generation in special net/rpc format.
func (g *Client) Generate(params *GenerateParams) (*GenerateResult, error) {
	result := &GenerateResult{}
	err := g.client.Call(cmdPluginGenerate, params, &result)
	if err != nil {
		panic(err)
	}

	return nil, nil
}

// Server allows to serve a request to plugin.
type Server struct {
	Impl Generator
}

// Generate runs underlying implementation of code generator.
func (g *Server) Generate(args *GenerateParams, resp *GenerateResult) (err error) {
	resp, err = g.Impl.Generate(args)
	return
}

// It is a worker for net/rpc.
type NetRPCWorker struct {
	Impl Generator
}

// Returns a server for net/rpc.
func (n *NetRPCWorker) Server(*plugin.MuxBroker) (interface{}, error) {
	return &Server{Impl: n.Impl}, nil
}

// Returns a client for net/rpc.
func (NetRPCWorker) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &Client{client: c}, nil
}
