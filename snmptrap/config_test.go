// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package snmptrapreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/snmptrapreceiver"

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap/confmaptest"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/snmptrapreceiver/internal/metadata"
)

func TestLoadConfigConnectionConfigs(t *testing.T) {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)

	factory := NewFactory()

	type testCase struct {
		name        string
		nameVal     string
		expectedCfg *Config
		expectedErr string
	}

	expectedConfigSimple := factory.CreateDefaultConfig().(*Config)

	expectedConfigInvalidListenAddress := factory.CreateDefaultConfig().(*Config)
	expectedConfigInvalidListenAddress.ListenAddress = "udp://a:a:a:a:a:a"

	expectedConfigNoPort := factory.CreateDefaultConfig().(*Config)
	expectedConfigNoPort.ListenAddress = "udp://localhost"

	expectedConfigNoPortTrailingColon := factory.CreateDefaultConfig().(*Config)
	expectedConfigNoPortTrailingColon.ListenAddress = "udp://localhost:"

	expectedConfigBadListenAddressScheme := factory.CreateDefaultConfig().(*Config)
	expectedConfigBadListenAddressScheme.ListenAddress = "http://localhost:162"

	expectedConfigNoListenAddressScheme := factory.CreateDefaultConfig().(*Config)
	expectedConfigNoListenAddressScheme.ListenAddress = "localhost:162"

	expectedConfigBadVersion := factory.CreateDefaultConfig().(*Config)
	expectedConfigBadVersion.Version = "9999"

	expectedConfigV3NoUser := factory.CreateDefaultConfig().(*Config)
	expectedConfigV3NoUser.Version = "v3"
	expectedConfigV3NoUser.SecurityLevel = "no_auth_no_priv"

	expectedConfigV3NoSecurityLevel := factory.CreateDefaultConfig().(*Config)
	expectedConfigV3NoSecurityLevel.Version = "v3"
	expectedConfigV3NoSecurityLevel.User = "u"

	expectedConfigV3BadSecurityLevel := factory.CreateDefaultConfig().(*Config)
	expectedConfigV3BadSecurityLevel.Version = "v3"
	expectedConfigV3BadSecurityLevel.SecurityLevel = "super"
	expectedConfigV3BadSecurityLevel.User = "u"

	expectedConfigV3NoAuthType := factory.CreateDefaultConfig().(*Config)
	expectedConfigV3NoAuthType.Version = "v3"
	expectedConfigV3NoAuthType.User = "u"
	expectedConfigV3NoAuthType.SecurityLevel = "auth_no_priv"
	expectedConfigV3NoAuthType.AuthPassword = "p"

	expectedConfigV3BadAuthType := factory.CreateDefaultConfig().(*Config)
	expectedConfigV3BadAuthType.Version = "v3"
	expectedConfigV3BadAuthType.User = "u"
	expectedConfigV3BadAuthType.SecurityLevel = "auth_no_priv"
	expectedConfigV3BadAuthType.AuthType = "super"
	expectedConfigV3BadAuthType.AuthPassword = "p"

	expectedConfigV3NoAuthPassword := factory.CreateDefaultConfig().(*Config)
	expectedConfigV3NoAuthPassword.Version = "v3"
	expectedConfigV3NoAuthPassword.User = "u"
	expectedConfigV3NoAuthPassword.SecurityLevel = "auth_no_priv"

	expectedConfigV3Simple := factory.CreateDefaultConfig().(*Config)
	expectedConfigV3Simple.Version = "v3"
	expectedConfigV3Simple.User = "u"
	expectedConfigV3Simple.SecurityLevel = "auth_priv"
	expectedConfigV3Simple.AuthPassword = "p"
	expectedConfigV3Simple.PrivacyPassword = "pp"

	expectedConfigV3BadPrivacyType := factory.CreateDefaultConfig().(*Config)
	expectedConfigV3BadPrivacyType.Version = "v3"
	expectedConfigV3BadPrivacyType.User = "u"
	expectedConfigV3BadPrivacyType.SecurityLevel = "auth_priv"
	expectedConfigV3BadPrivacyType.AuthPassword = "p"
	expectedConfigV3BadPrivacyType.PrivacyType = "super"
	expectedConfigV3BadPrivacyType.PrivacyPassword = "pp"

	expectedConfigV3NoPrivacyPassword := factory.CreateDefaultConfig().(*Config)
	expectedConfigV3NoPrivacyPassword.Version = "v3"
	expectedConfigV3NoPrivacyPassword.User = "u"
	expectedConfigV3NoPrivacyPassword.SecurityLevel = "auth_priv"
	expectedConfigV3NoPrivacyPassword.AuthPassword = "p"

	testCases := []testCase{
		{
			name:        "NoListenAddressUsesDefault",
			nameVal:     "no_endpoint",
			expectedCfg: expectedConfigSimple,
			expectedErr: "",
		},
		{
			name:        "InvalidListenAddressErrors",
			nameVal:     "invalid_endpoint",
			expectedCfg: expectedConfigInvalidListenAddress,
			expectedErr: fmt.Sprintf(errMsgInvalidListenAddress[:len(errMsgInvalidListenAddress)-2], "udp://a:a:a:a:a:a"),
		},
		{
			name:        "NoPortErrors",
			nameVal:     "no_port",
			expectedCfg: expectedConfigNoPort,
			expectedErr: fmt.Sprintf(errMsgInvalidListenAddress[:len(errMsgInvalidListenAddress)-2], "udp://localhost"),
		},
		{
			name:        "NoPortTrailingColonErrors",
			nameVal:     "no_port_trailing_colon",
			expectedCfg: expectedConfigNoPortTrailingColon,
			expectedErr: fmt.Sprintf(errMsgInvalidListenAddress[:len(errMsgInvalidListenAddress)-2], "udp://localhost:"),
		},
		{
			name:        "BadListenAddressSchemeErrors",
			nameVal:     "bad_endpoint_scheme",
			expectedCfg: expectedConfigBadListenAddressScheme,
			expectedErr: errListenAddressBadScheme.Error(),
		},
		{
			name:        "NoListenAddressSchemeErrors",
			nameVal:     "no_endpoint_scheme",
			expectedCfg: expectedConfigNoListenAddressScheme,
			expectedErr: fmt.Sprintf(errMsgInvalidListenAddress[:len(errMsgInvalidListenAddress)-2], "localhost:162"),
		},
		{
			name:        "NoVersionUsesDefault",
			nameVal:     "no_version",
			expectedCfg: expectedConfigSimple,
			expectedErr: "",
		},
		{
			name:        "BadVersionErrors",
			nameVal:     "bad_version",
			expectedCfg: expectedConfigBadVersion,
			expectedErr: errBadVersion.Error(),
		},
		{
			name:        "V3NoUserErrors",
			nameVal:     "v3_no_user",
			expectedCfg: expectedConfigV3NoUser,
			expectedErr: errEmptyUser.Error(),
		},
		{
			name:        "V3NoSecurityLevelUsesDefault",
			nameVal:     "v3_no_security_level",
			expectedCfg: expectedConfigV3NoSecurityLevel,
			expectedErr: "",
		},
		{
			name:        "V3BadSecurityLevelErrors",
			nameVal:     "v3_bad_security_level",
			expectedCfg: expectedConfigV3BadSecurityLevel,
			expectedErr: errBadSecurityLevel.Error(),
		},
		{
			name:        "V3NoAuthTypeUsesDefault",
			nameVal:     "v3_no_auth_type",
			expectedCfg: expectedConfigV3NoAuthType,
			expectedErr: "",
		},
		{
			name:        "V3BadAuthTypeErrors",
			nameVal:     "v3_bad_auth_type",
			expectedCfg: expectedConfigV3BadAuthType,
			expectedErr: errBadAuthType.Error(),
		},
		{
			name:        "V3NoAuthPasswordErrors",
			nameVal:     "v3_no_auth_password",
			expectedCfg: expectedConfigV3NoAuthPassword,
			expectedErr: errEmptyAuthPassword.Error(),
		},
		{
			name:        "V3NoPrivacyTypeUsesDefault",
			nameVal:     "v3_no_privacy_type",
			expectedCfg: expectedConfigV3Simple,
			expectedErr: "",
		},
		{
			name:        "V3BadPrivacyTypeErrors",
			nameVal:     "v3_bad_privacy_type",
			expectedCfg: expectedConfigV3BadPrivacyType,
			expectedErr: errBadPrivacyType.Error(),
		},
		{
			name:        "V3NoPrivacyPasswordErrors",
			nameVal:     "v3_no_privacy_password",
			expectedCfg: expectedConfigV3NoPrivacyPassword,
			expectedErr: errEmptyPrivacyPassword.Error(),
		},
		{
			name:        "GoodV2CConnectionNoErrors",
			nameVal:     "v2c_connection_good",
			expectedCfg: expectedConfigSimple,
			expectedErr: "",
		},
		{
			name:        "GoodV3ConnectionNoErrors",
			nameVal:     "v3_connection_good",
			expectedCfg: expectedConfigV3Simple,
			expectedErr: "",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			sub, err := cm.Sub(component.NewIDWithName(metadata.Type, test.nameVal).String())
			require.NoError(t, err)

			cfg := factory.CreateDefaultConfig()
			require.NoError(t, component.UnmarshalConfig(sub, cfg))
			if test.expectedErr == "" {
				require.NoError(t, component.ValidateConfig(cfg))
			} else {
				require.ErrorContains(t, component.ValidateConfig(cfg), test.expectedErr)
			}

			require.Equal(t, test.expectedCfg, cfg)
		})
	}
}


func getBaseAttrConfig(attrType string) map[string]*AttributeConfig {
	switch attrType {
	case "oid":
		return map[string]*AttributeConfig{
			"a2": {
				OID: "1",
			},
		}
	case "prefix":
		return map[string]*AttributeConfig{
			"a2": {
				IndexedValuePrefix: "p",
			},
		}
	default:
		return map[string]*AttributeConfig{
			"a2": {
				Enum: []string{"val1", "val2"},
			},
		}
	}
}

func getBaseResourceAttrConfig(attrType string) map[string]*ResourceAttributeConfig {
	switch attrType {
	case "oid":
		return map[string]*ResourceAttributeConfig{
			"ra1": {
				OID: "2",
			},
		}
	case "scalar_oid":
		return map[string]*ResourceAttributeConfig{
			"ra1": {
				ScalarOID: "0",
			},
		}
	default:
		return map[string]*ResourceAttributeConfig{
			"ra1": {
				IndexedValuePrefix: "p",
			},
		}
	}
}


// Testing Validate directly to test that missing data errors when no defaults are provided
func TestValidate(t *testing.T) {
	type testCase struct {
		name        string
		cfg         *Config
		expectedErr string
	}

	testCases := []testCase{
		{
			name: "NoListenAddressErrors",
			cfg: &Config{
				Version:   "v2c",
				Community: "public",
			},
			expectedErr: errEmptyListenAddress.Error(),
		},
		{
			name: "NoVersionErrors",
			cfg: &Config{
				ListenAddress:  "udp://localhost:162",
				Community: "public",
			},
			expectedErr: errEmptyVersion.Error(),
		},
		{
			name: "V3NoSecurityLevelErrors",
			cfg: &Config{
				ListenAddress: "udp://localhost:162",
				Version:  "v3",
				User:     "u",
			},
			expectedErr: errEmptySecurityLevel.Error(),
		},
		{
			name: "V3NoAuthTypeErrors",
			cfg: &Config{
				ListenAddress:      "udp://localhost:162",
				Version:       "v3",
				SecurityLevel: "auth_no_priv",
				User:          "u",
				AuthPassword:  "p",
			},
			expectedErr: errEmptyAuthType.Error(),
		},
		{
			name: "V3NoPrivacyTypeErrors",
			cfg: &Config{
				ListenAddress:        "udp://localhost:162",
				Version:         "v3",
				SecurityLevel:   "auth_priv",
				User:            "u",
				AuthType:        "md5",
				AuthPassword:    "p",
				PrivacyPassword: "pp",
			},
			expectedErr: errEmptyPrivacyType.Error(),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			err := test.cfg.Validate()
			assert.ErrorContains(t, err, test.expectedErr)
		})
	}
}
