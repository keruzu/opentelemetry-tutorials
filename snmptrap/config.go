// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
package snmptrapreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/snmptrapreceiver"

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"go.opentelemetry.io/collector/config/configopaque"
)

// Config Defaults
const (
	defaultTimeout            = 5 * time.Second  // In seconds
	defaultListenAddress           = "udp://localhost:162"
	defaultVersion            = "v2c"
	defaultCommunity          = "public"
	defaultSecurityLevel      = "no_auth_no_priv"
	defaultAuthType           = "MD5"
	defaultPrivacyType        = "DES"
)

var (
	// Config error messages
	errMsgInvalidListenAddressWError                     = `invalid endpoint '%s': must be in '[scheme]://[host]:[port]' format: %w`
	errMsgInvalidListenAddress                           = `invalid endpoint '%s': must be in '[scheme]://[host]:[port]' format`

	// Config errors
	errEmptyListenAddress        = errors.New("endpoint must be specified")
	errListenAddressBadScheme    = errors.New("endpoint scheme must be either tcp, tcp4, tcp6, udp, udp4, or udp6")
	errEmptyVersion         = errors.New("version must specified")
	errBadVersion           = errors.New("version must be either v1, v2c, or v3")
	errEmptyUser            = errors.New("user must be specified when version is v3")
	errEmptySecurityLevel   = errors.New("security_level must be specified when version is v3")
	errBadSecurityLevel     = errors.New("security_level must be either no_auth_no_priv, auth_no_priv, or auth_priv")
	errEmptyAuthType        = errors.New("auth_type must be specified when security_level is auth_no_priv or auth_priv")
	errBadAuthType          = errors.New("auth_type must be either MD5, SHA, SHA224, SHA256, SHA384, SHA512")
	errEmptyAuthPassword    = errors.New("auth_password must be specified when security_level is auth_no_priv or auth_priv")
	errEmptyPrivacyType     = errors.New("privacy_type must be specified when security_level is auth_priv")
	errBadPrivacyType       = errors.New("privacy_type must be either DES, AES, AES192, AES192C, AES256, AES256C")
	errEmptyPrivacyPassword = errors.New("privacy_password must be specified when security_level is auth_priv")
)

// Config defines the configuration for the various elements of the receiver.
type Config struct {
	// ListenAddress is the host (IP or hostname) + port to listen on. Must be formatted as [udp|tcp|][4|6|]://{host}:{port}.
	// Default: udp://localhost:162
	// If no scheme is given, udp4 is assumed.
	// If no port is given, 162 is assumed.
	ListenAddress string `mapstructure:"listen_address"`

	// Version is the version of SNMP to use for this connection.
	// Valid options: v1, v2c, v3.
	// Default: v2c
	Version string `mapstructure:"version"`

	// Community is the SNMP community string to use.
	// Only valid for versions "v1" and "v2c"
	// Default: public
	Community string `mapstructure:"community"`

	// User is the SNMP User for this connection.
	// Only valid for version “v3”
	User string `mapstructure:"user"`

	// SecurityLevel is the security level to use for this SNMP connection.
	// Only valid for version “v3”
	// Valid options: “no_auth_no_priv”, “auth_no_priv”, “auth_priv”
	// Default: "no_auth_no_priv"
	SecurityLevel string `mapstructure:"security_level"`

	// AuthType is the type of authentication protocol to use for this SNMP connection.
	// Only valid for version “v3” and if “no_auth_no_priv” is not selected for SecurityLevel
	// Valid options: “md5”, “sha”, “sha224”, “sha256”, “sha384”, “sha512”
	// Default: "md5"
	AuthType string `mapstructure:"auth_type"`

	// AuthPassword is the authentication password used for this SNMP connection.
	// Only valid for version "v3" and if "no_auth_no_priv" is not selected for SecurityLevel
	AuthPassword configopaque.String `mapstructure:"auth_password"`

	// PrivacyType is the type of privacy protocol to use for this SNMP connection.
	// Only valid for version “v3” and if "auth_priv" is selected for SecurityLevel
	// Valid options: “des”, “aes”, “aes192”, “aes256”, “aes192c”, “aes256c”
	// Default: "des"
	PrivacyType string `mapstructure:"privacy_type"`

	// PrivacyPassword is the authentication password used for this SNMP connection.
	// Only valid for version “v3” and if "auth_priv" is selected for SecurityLevel
	PrivacyPassword configopaque.String `mapstructure:"privacy_password"`

	// CloseTimeout is the max wait time for the socket to gracefully signal its closure.
	CloseTimeout time.Duration `mapstructure:"listener_close_timeout"`

}

// FIXME: require a representation of a PDU

// ResourceAttributeConfig contains config info about all of the resource attributes that will be used by this receiver.
type ResourceAttributeConfig struct {
	// Description is optional and describes what the resource attribute represents
	Description string `mapstructure:"description"`
	// OID is required only if ScalarOID or IndexedValuePrefix is not set.
	// This is the column OID which will provide indexed values to be used for this resource attribute. These indexed values
	// will ultimately each be associated with a different "resource" as an attribute on that resource. Indexed metric values
	// will then be used to associate metric datapoints to the matching "resource" (based on matching indexes).
	OID string `mapstructure:"oid"`
	// ScalarOID is required only if OID or IndexedValuePrefix is not set.
	// This is the scalar OID which will provide a value to be used for this resource attribute.
	// Single or indexed metrics can then be associated with the resource. (Indexed metrics also need an indexed attribute or resource attribute to associate with a scalar metric resource attribute)
	ScalarOID string `mapstructure:"scalar_oid"`
	// IndexedValuePrefix is required only if OID or ScalarOID is not set.
	// This will be used alongside indexed metric values for this resource attribute. The prefix value concatenated with
	// specific indexes of metric indexed values (Ex: prefix.1.2) will ultimately each be associated with a different "resource"
	// as an attribute on that resource. The related indexed metric values will then be used to associate metric datapoints to
	// those resources.
	IndexedValuePrefix string `mapstructure:"indexed_value_prefix"` // required and valid if no oid or scalar_oid field
}

// AttributeConfig contains config info about all of the metric attributes that will be used by this receiver.
type AttributeConfig struct {
	// Value is optional, and will allow for a different attribute key other than the attribute name
	Value string `mapstructure:"value"`
	// Description is optional and describes what the attribute represents
	Description string `mapstructure:"description"`
	// Enum is required only if OID and IndexedValuePrefix are not defined.
	// This contains a list of possible values that can be associated with this attribute
	Enum []string `mapstructure:"enum"`
	// OID is required only if Enum and IndexedValuePrefix are not defined.
	// This is the column OID which will provide indexed values to be uased for this attribute (alongside a metric with ColumnOIDs)
	OID string `mapstructure:"oid"`
	// IndexedValuePrefix is required only if Enum and OID are not defined.
	// This is used alongside metrics with ColumnOIDs to assign attribute values using this prefix + the OID index of the metric value
	IndexedValuePrefix string `mapstructure:"indexed_value_prefix"`
}


// Attribute is a connection between a metric configuration and an AttributeConfig
type Attribute struct {
	// Name is required and should match the key for an AttributeConfig
	Name string `mapstructure:"name"`
	// Value is optional and is only needed for a matched AttributeConfig's with enum value.
	// Value should match one of the AttributeConfig's enum values in this case
	Value string `mapstructure:"value"`
}

// Validate validates the given config, returning an error specifying any issues with the config.
func (cfg *Config) Validate() error {
	var combinedErr error

	combinedErr = errors.Join(combinedErr, validateListenAddress(cfg))
	combinedErr = errors.Join(combinedErr, validateVersion(cfg))
	if strings.ToUpper(cfg.Version) == "V3" {
		combinedErr = errors.Join(combinedErr, validateSecurity(cfg))
	}

	return combinedErr
}

// validateListenAddress validates the ListenAddress
func validateListenAddress(cfg *Config) error {
	if cfg.ListenAddress == "" {
		return errEmptyListenAddress
	}

	// Ensure valid endpoint
	u, err := url.Parse(cfg.ListenAddress)
	if err != nil {
		return fmt.Errorf(errMsgInvalidListenAddressWError, cfg.ListenAddress, err)
	}
	if u.Host == "" || u.Port() == "" {
		return fmt.Errorf(errMsgInvalidListenAddress, cfg.ListenAddress)
	}

	// Ensure valid scheme
	switch strings.ToUpper(u.Scheme) {
	case "TCP", "TCP4", "TCP6", "UDP", "UDP4", "UDP6": // ok
	default:
		return errListenAddressBadScheme
	}

	return nil
}

// validateVersion validates the Version
func validateVersion(cfg *Config) error {
	if cfg.Version == "" {
		return errEmptyVersion
	}

	// Ensure valid version
	switch strings.ToUpper(cfg.Version) {
	case "V1", "V2C", "V3": // ok
	default:
		return errBadVersion
	}

	return nil
}

// validateSecurity validates all v3 related security configs
func validateSecurity(cfg *Config) error {
	var combinedErr error

	// Ensure valid user
	if cfg.User == "" {
		combinedErr = errors.Join(combinedErr, errEmptyUser)
	}

	if cfg.SecurityLevel == "" {
		return errors.Join(combinedErr, errEmptySecurityLevel)
	}

	// Ensure valid security level
	switch strings.ToUpper(cfg.SecurityLevel) {
	case "NO_AUTH_NO_PRIV":
		return combinedErr
	case "AUTH_NO_PRIV":
		// Ensure valid auth configs
		return errors.Join(combinedErr, validateAuth(cfg))
	case "AUTH_PRIV": // ok
		// Ensure valid auth and privacy configs
		combinedErr = errors.Join(combinedErr, validateAuth(cfg))
		return errors.Join(combinedErr, validatePrivacy(cfg))
	default:
		return errors.Join(combinedErr, errBadSecurityLevel)
	}
}

// validateAuth validates the AuthType and AuthPassword
func validateAuth(cfg *Config) error {
	var combinedErr error

	// Ensure valid auth password
	if cfg.AuthPassword == "" {
		combinedErr = errors.Join(combinedErr, errEmptyAuthPassword)
	}

	// Ensure valid auth type
	if cfg.AuthType == "" {
		return errors.Join(combinedErr, errEmptyAuthType)
	}

	switch strings.ToUpper(cfg.AuthType) {
	case "MD5", "SHA", "SHA224", "SHA256", "SHA384", "SHA512": // ok
	default:
		combinedErr = errors.Join(combinedErr, errBadAuthType)
	}

	return combinedErr
}

// validatePrivacy validates the PrivacyType and PrivacyPassword
func validatePrivacy(cfg *Config) error {
	var combinedErr error

	// Ensure valid privacy password
	if cfg.PrivacyPassword == "" {
		combinedErr = errors.Join(combinedErr, errEmptyPrivacyPassword)
	}

	// Ensure valid privacy type
	if cfg.PrivacyType == "" {
		return errors.Join(combinedErr, errEmptyPrivacyType)
	}

	switch strings.ToUpper(cfg.PrivacyType) {
	case "DES", "AES", "AES192", "AES192C", "AES256", "AES256C": // ok
	default:
		combinedErr = errors.Join(combinedErr, errBadPrivacyType)
	}

	return combinedErr
}



// contains checks if string slice contains a string value
func contains(elements []string, value string) bool {
	for _, element := range elements {
		if value == element {
			return true
		}
	}
	return false
}
