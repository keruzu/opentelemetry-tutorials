// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package snmptrapreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/snmptrapreceiver"

import (
	//"context"
	"errors"
	"net/url"
	"strings"

	"go.opentelemetry.io/collector/component"
	//"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/snmptrapreceiver/internal/metadata"
)

var errConfigNotSNMP = errors.New("config was not a SNMP trap receiver config")

// NewFactory creates a new receiver factory for SNMP
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		metadata.Type,
		createDefaultConfig,
		receiver.WithMetrics(createMetricsReceiver, metadata.MetricsStability))
}

// createDefaultConfig creates a config for SNMP with as many default values as possible
func createDefaultConfig() component.Config {
	return &Config{
		ListenAddress:      defaultListenAddress,
		Version:       defaultVersion,
		Community:     defaultCommunity,
		SecurityLevel: defaultSecurityLevel,
		AuthType:      defaultAuthType,
		PrivacyType:   defaultPrivacyType,
	}
}

// addMissingConfigDefaults adds any missing config parameters that have defaults
func addMissingConfigDefaults(cfg *Config) error {
	// Add the schema prefix to the endpoint if it doesn't contain one
	if !strings.Contains(cfg.ListenAddress, "://") {
		cfg.ListenAddress = "udp://" + cfg.ListenAddress
	}

	// Add default port to endpoint if it doesn't contain one
	u, err := url.Parse(cfg.ListenAddress)
	if err == nil && u.Port() == "" {
		portSuffix := "162"
		if cfg.ListenAddress[len(cfg.ListenAddress)-1:] != ":" {
			portSuffix = ":" + portSuffix
		}
		cfg.ListenAddress += portSuffix
	}

	return component.ValidateConfig(cfg)
}
