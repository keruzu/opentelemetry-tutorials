// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package snmptrapreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/snmptrapreceiver"

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/snmptrapreceiver/internal/metadata"
)

var errConfigNotSNMP = errors.New("config was not a SNMP trap receiver config")

// NewFactory creates a new receiver factory for SNMP
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		metadata.Type,
		createDefaultConfig,
		receiver.WithLogs(createLogsReceiver, metadata.LogsStability))
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
	// Add the schema prefix to the listen address if it doesn't contain one
	if !strings.Contains(cfg.ListenAddress, "://") {
		cfg.ListenAddress = "udp://" + cfg.ListenAddress
	}

	// Add default port to listen address if it doesn't contain one
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

// createLogsReceiver creates a logs receiver based on provided config.
func createLogsReceiver(
	_ context.Context,
	params receiver.CreateSettings,
	cfg component.Config,
	consumer consumer.Logs,
) (receiver.Logs, error) {

	snmpConfig, ok := cfg.(*Config)
	if !ok {
		return nil, errConfigNotSNMP
	}

	if err := addMissingConfigDefaults(snmpConfig); err != nil {
		return nil, fmt.Errorf("failed to validate added config defaults: %w", err)
	}

	// FIXME: add in sane code here
	return nil, nil
}
