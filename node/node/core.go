package node

import (
	"github.com/armon/go-metrics"
	"github.com/rs/zerolog"

	"github.com/blocklessnetwork/b7s/host"
	"github.com/blocklessnetwork/b7s/telemetry/tracing"
)

type Core struct {
	Log  zerolog.Logger
	Host *host.Host

	// Telemetry
	Tracer  *tracing.Tracer
	Metrics *metrics.Metrics
}

func NewCore(log zerolog.Logger, host *host.Host) *Core {

	core := &Core{
		Log:  log,
		Host: host,
		// tracer:   tracing.NewTracer(tracerName),
		Metrics: metrics.Default(),
	}

	return core
}
