package egrpcinteceptor

import (
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

// Semantic conventions for attribute keys for gRPC.
const (
	// RPCNameKey Name of message transmitted or received.
	RPCNameKey = attribute.Key("name")

	// RPCMessageTypeKey Type of message transmitted or received.
	RPCMessageTypeKey = attribute.Key("message.type")

	// RPCMessageIDKey Identifier of message transmitted or received.
	RPCMessageIDKey = attribute.Key("message.id")

	// RPCMessageCompressedSizeKey The compressed size of the message transmitted or received in bytes.
	RPCMessageCompressedSizeKey = attribute.Key("message.compressed_size")

	// RPCMessageUncompressedSizeKey The uncompressed size of the message transmitted or received in
	// bytes.
	RPCMessageUncompressedSizeKey = attribute.Key("message.uncompressed_size")

	GRPCKindKey = attribute.Key("rpc.grpc.kind")
)

// Semantic conventions for common RPC attributes.
var (
	// Semantic convention for gRPC as the remoting system.
	RPCSystemGRPC = semconv.RPCSystemKey.String("grpc")

	// Semantic convention for a message named message.
	RPCNameMessage = RPCNameKey.String("message")

	// Semantic conventions for RPC message types.
	RPCMessageTypeSent     = RPCMessageTypeKey.String("SENT")
	RPCMessageTypeReceived = RPCMessageTypeKey.String("RECEIVED")

	GRPCKindUnary  = GRPCKindKey.String("unary")
	GRPCKindStream = GRPCKindKey.String("stream")
)
