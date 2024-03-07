
package snmptrapreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/snmptrapreceiver"

import (
	"context"
	"time"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"
)

type snmptrapReceiver struct{
	host component.Host
	cancel context.CancelFunc
	logger *zap.Logger
	nextConsumer consumer.Traces
	config *Config
}

func (snmptrapRcvr *snmptrapReceiver) Start(ctx context.Context, host component.Host) error {
	snmptrapRcvr.host = host
	ctx = context.Background()
	ctx, snmptrapRcvr.cancel = context.WithCancel(ctx)

	interval, _ := time.ParseDuration("1m")
	go func() {
		ticker:= time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <- ticker.C:
				snmptrapRcvr.logger.Info("I should start processing logs now!")
			case <- ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (snmptrapRcvr *snmptrapReceiver) Shutdown(ctx context.Context) error {
	snmptrapRcvr.cancel()
	return nil
}

