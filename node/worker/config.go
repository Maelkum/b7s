package worker

import (
	"errors"
	"path/filepath"
	"time"

	"github.com/hashicorp/go-multierror"

	"github.com/blocklessnetwork/b7s/metadata"
	"github.com/blocklessnetwork/b7s/models/blockless"
)

// Option can be used to set Node configuration options.
type Option func(*Config)

// DefaultConfig represents the default settings for the node.
var DefaultConfig = Config{
	Topics:           []string{blockless.DefaultTopic},
	HealthInterval:   blockless.DefaultHealthInterval,
	Concurrency:      blockless.DefaultConcurrency,
	LoadAttributes:   DefaultAttributeLoadingSetting,
	MetadataProvider: metadata.NewNoopProvider(),
}

// Config represents the Node configuration.
type Config struct {
	Topics           []string          // Topics to subscribe to.
	HealthInterval   time.Duration     // How often should we emit the health ping.
	Concurrency      uint              // How many requests should the node process in parallel.
	Workspace        string            // Directory where we can store files needed for execution.
	LoadAttributes   bool              // Node should try to load its attributes from IPFS.
	MetadataProvider metadata.Provider // Metadata provider for the node
}

// Validate checks if the given configuration is correct.
func (c Config) Valid() error {

	var err *multierror.Error

	if len(c.Topics) == 0 {
		err = multierror.Append(err, errors.New("topics cannot be empty"))
	}

	if !filepath.IsAbs(c.Workspace) {
		err = multierror.Append(err, errors.New("workspace must be an absolute path"))
	}

	return err.ErrorOrNil()
}

// WithTopics specifies the p2p topics to which node should subscribe.
func WithTopics(topics []string) Option {
	return func(cfg *Config) {
		cfg.Topics = topics
	}
}

// WithHealthInterval specifies how often we should emit the health signal.
func WithHealthInterval(d time.Duration) Option {
	return func(cfg *Config) {
		cfg.HealthInterval = d
	}
}

// WithConcurrency specifies how many requests the node should process in parallel.
func WithConcurrency(n uint) Option {
	return func(cfg *Config) {
		cfg.Concurrency = n
	}
}

// WithWorkspace specifies the workspace that the node can use for file storage.
func WithWorkspace(path string) Option {
	return func(cfg *Config) {
		cfg.Workspace = path
	}
}

// WithAttributeLoading specifies whether node should try to load its attributes data from IPFS.
func WithAttributeLoading(b bool) Option {
	return func(cfg *Config) {
		cfg.LoadAttributes = b
	}
}

// WithMetadataProvider sets the metadata provider for the node.
func WithMetadataProvider(p metadata.Provider) Option {
	return func(cfg *Config) {
		cfg.MetadataProvider = p
	}
}
