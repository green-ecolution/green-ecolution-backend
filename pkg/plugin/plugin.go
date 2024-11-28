package plugin

import "net/url"

// Plugin represents a plugin that can be registered with the plugin host.
type Plugin struct {
	Slug           string
	Name           string
	Version        string
	Description    string
	PluginHostPath *url.URL
}

// PluginOption is a functional option for configuring a Plugin.
type PluginOption func(*Plugin)

// WithName sets the name of the plugin.
func WithName(name string) PluginOption {
	return func(p *Plugin) {
		p.Name = name
	}
}

// WithSlug sets the slug of the plugin.
func WithSlug(slug string) PluginOption {
	return func(p *Plugin) {
		p.Slug = slug
	}
}

// WithVersion sets the version of the plugin.
func WithVersion(version string) PluginOption {
	return func(p *Plugin) {
		p.Version = version
	}
}

// WithDescription sets the description of the plugin.
func WithDescription(description string) PluginOption {
	return func(p *Plugin) {
		p.Description = description
	}
}

// WithHostPath sets the host path of the plugin.
func WithHostPath(hostPath *url.URL) PluginOption {
	return func(p *Plugin) {
		p.PluginHostPath = hostPath
	}
}

var defaultPlugin = Plugin{
	Version: "develop",
}

// NewPlugin creates a new Plugin with the provided options. If no options are provided, the default values are used.
func NewPlugin(opts ...PluginOption) *Plugin {
	p := defaultPlugin
	for _, opt := range opts {
		opt(&p)
	}
	return &p
}
