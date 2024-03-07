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


// getAttributeConfigValue returns the value of an attribute config
func (h configHelper) getAttributeConfigValue(name string) string {
	attrConfig := h.cfg.Attributes[name]
	if attrConfig == nil {
		return ""
	}

	return attrConfig.Value
}

// getAttributeConfigIndexedValuePrefix returns the indexed value prefix of an attribute config
func (h configHelper) getAttributeConfigIndexedValuePrefix(name string) string {
	attrConfig := h.cfg.Attributes[name]
	if attrConfig == nil {
		return ""
	}

	return attrConfig.IndexedValuePrefix
}


