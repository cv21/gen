package pkg

import (
	"encoding/gob"
	"net/rpc"

	plugin "github.com/hashicorp/go-plugin"
	"github.com/vetcher/go-astra/types"
)

// It is default command for each generator plugin.
const cmdPluginGenerate = "Plugin.Generate"

// DefaultHandshakeConfig useful for plugin compatibility specification.
var DefaultHandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "HOST",
	MagicCookieValue: "GEN",
}

// It is a result of generator plugin work.
// PluginGenerateResult helps to make solid division between generator error (which is normal)
// and plugin error (which is outstanding).
type PluginGenerateResult struct {
	GenerateResult *GenerateResult
	Error          error
}

// It is a plugin client.
type Client struct {
	client *rpc.Client
}

// Generate method implements Generator interface.
// It calls plugin generation in special net/rpc format.
func (g *Client) Generate(params *GenerateParams) (*GenerateResult, error) {
	result := &PluginGenerateResult{}
	err := g.client.Call(cmdPluginGenerate, params, &result)
	if err != nil {
		panic(err)
	}

	return result.GenerateResult, result.Error
}

// Server allows to serve a request to plugin.
type Server struct {
	Impl Generator
}

// Generate runs underlying implementation of code generator.
func (g *Server) Generate(args *GenerateParams, resp *PluginGenerateResult) error {
	resp.GenerateResult, resp.Error = g.Impl.Generate(args)
	return nil
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

// Registers specific types for gob encoding/decoding.
// This work properly this function MUST be called by both client and server.
func RegisterGobTypes() {
	gob.Register(types.TInterface{})
	gob.Register(types.TMap{})
	gob.Register(types.TName{})
	gob.Register(types.TPointer{})
	gob.Register(types.TArray{})
	gob.Register(types.TImport{})
	gob.Register(types.TEllipsis{})
	gob.Register(types.TChan{})
}
