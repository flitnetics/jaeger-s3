// Copyright (c) 2020 The Jaeger Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shared

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	_ "github.com/jaegertracing/jaeger/pkg/gogocodec" // force gogo codec registration
	"github.com/jaegertracing/jaeger/proto-gen/storage_v1"
)

// Ensure plugin.GRPCPlugin API match.
var _ plugin.GRPCPlugin = (*StorageGRPCPlugin)(nil)

// StorageGRPCPlugin is the implementation of plugin.GRPCPlugin.
type StorageGRPCPlugin struct {
	plugin.Plugin

	// Concrete implementation, This is only used for plugins that are written in Go.
	Impl        StoragePlugin
	ArchiveImpl ArchiveStoragePlugin
	StreamImpl  StreamingSpanWriterPlugin
}

// RegisterHandlers registers the plugin with the server
func (p *StorageGRPCPlugin) RegisterHandlers(s *grpc.Server) error {
	server := &grpcServer{
		Impl:        p.Impl,
		ArchiveImpl: p.ArchiveImpl,
		StreamImpl:  p.StreamImpl,
	}
	storage_v1.RegisterSpanReaderPluginServer(s, server)
	storage_v1.RegisterSpanWriterPluginServer(s, server)
	storage_v1.RegisterArchiveSpanReaderPluginServer(s, server)
	storage_v1.RegisterArchiveSpanWriterPluginServer(s, server)
	storage_v1.RegisterPluginCapabilitiesServer(s, server)
	storage_v1.RegisterDependenciesReaderPluginServer(s, server)
	storage_v1.RegisterStreamingSpanWriterPluginServer(s, server)
	return nil
}

// GRPCServer implements plugin.GRPCPlugin. It is used by go-plugin to create a grpc plugin server.
func (p *StorageGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	return p.RegisterHandlers(s)
}

// GRPCClient implements plugin.GRPCPlugin. It is used by go-plugin to create a grpc plugin client.
func (*StorageGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return NewGRPCClient(c), nil
}
