package shared

import "github.com/hashicorp/go-plugin"

var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "TASK_PLUGIN",
	MagicCookieValue: "aurora",
}
