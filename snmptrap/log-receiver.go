
package snmptrapreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/snmptrapreceiver"

import (
	"context"
//	"time"
        "net"

	"go.opentelemetry.io/collector/component"
//	"go.opentelemetry.io/collector/consumer"
//	"go.uber.org/zap"

	gosnmp "github.com/gosnmp/gosnmp"
)

type snmptrapReceiver struct{
	host component.Host
	cancel context.CancelFunc
	config *Config
	listener  *gosnmp.TrapListener
}


// Start will create a SNMP trap listener with a callback handler
func (snmptrapRcvr *snmptrapReceiver) Start(ctx context.Context, host component.Host) error {
	snmptrapRcvr.host = host
	ctx = context.Background()
	ctx, snmptrapRcvr.cancel = context.WithCancel(ctx)

	// A TrapListener defines parameters for running a SNMP Trap receiver
	// nil values will be replaced by default values.
	snmptrapRcvr.listener = gosnmp.NewTrapListener()

	// When the listener receives a trap, invoke this callback handler
	snmptrapRcvr.listener.OnNewTrap = trapCallback

	snmptrapRcvr.listener.Params = gosnmp.Default
	snmptrapRcvr.listener.Params.Community = snmptrapRcvr.config.Community

	return nil
}

// Shutdown will stop our listner
func (snmptrapRcvr *snmptrapReceiver) Shutdown(ctx context.Context) error {
	snmptrapRcvr.listener.Close()
	snmptrapRcvr.cancel()
	return nil
}

// trapCallback is the callback for handling traps received by the listener.
//
func trapCallback(rawtrap *gosnmp.SnmpPacket, addr *net.UDPAddr) {
}

