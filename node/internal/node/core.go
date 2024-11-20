package node

import (
	"context"
	"time"

	"github.com/armon/go-metrics"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/rs/zerolog"

	"github.com/blocklessnetwork/b7s/host"
	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/node/internal/syncmap"
	"github.com/blocklessnetwork/b7s/telemetry/tracing"
)

type Core interface {
	Logger
	Network
	Telemetry
	NodeOps
}

type Logger interface {
	Log() *zerolog.Logger
}

type Network interface {
	Host() *host.Host // TODO Further abstract away Host()
	Messaging
}

type Messaging interface {
	Connected(peer.ID) bool

	Send(context.Context, peer.ID, blockless.Message) error
	SendToMany(context.Context, []peer.ID, blockless.Message, bool) error

	JoinTopic(string) error
	Subscribe(context.Context, string) error
	Publish(context.Context, blockless.Message) error
	PublishToTopic(context.Context, string, blockless.Message) error
}

type Telemetry interface {
	Tracer() *tracing.Tracer
	Metrics() *metrics.Metrics
}

type NodeOps interface {
	Run(context.Context)
	SetupRunFunc(healthInterval time.Duration, concurrency uint, process ProcessFunc) RunFunc
}

type core struct {
	log  zerolog.Logger
	host *host.Host

	topics *syncmap.Map[string, topicInfo]

	// Telemetry
	tracer  *tracing.Tracer
	metrics *metrics.Metrics
}

func NewCore(log zerolog.Logger, host *host.Host) *core {

	core := &core{
		log:  log,
		host: host,
		// tracer:   tracing.NewTracer(tracerName),
		metrics: metrics.Default(),
		topics:  syncmap.New[string, topicInfo](),
	}

	return core
}

func (c *core) Log() *zerolog.Logger {
	return &c.log
}

func (c *core) Host() *host.Host {
	return c.host
}

func (c *core) Tracer() *tracing.Tracer {
	return c.tracer
}

func (c *core) Metrics() *metrics.Metrics {
	return c.metrics
}

func (c *core) Network() {
	c.host.Network()
}
