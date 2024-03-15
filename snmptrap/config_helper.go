// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package snmptrapreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/snmptrapreceiver"

// import (
// )

// configHelper contains many of the functions required to get various info from the SNMP config
type configHelper struct {
	cfg                         *Config
}

// newConfigHelper returns a new configHelper with various pieces of static info saved for easy access
func newConfigHelper(cfg *Config) *configHelper {
	ch := configHelper{
		cfg:                         cfg,
	}

	return &ch
}

