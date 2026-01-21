package provider

import (
	"os"
	"runtime/debug"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-provider-openstack/terraform-provider-openstack/v3/openstack/version"
)

func getStringValue(configValue types.String, envVars []string, defaultValue string) string {
	if !configValue.IsNull() && !configValue.IsUnknown() {
		return configValue.ValueString()
	}

	for _, envVar := range envVars {
		if value := os.Getenv(envVar); value != "" {
			return value
		}
	}

	return defaultValue
}

func getBoolValue(configValue types.Bool, envVars []string, defaultValue bool) bool {
	if !configValue.IsNull() && !configValue.IsUnknown() {
		return configValue.ValueBool()
	}

	for _, envVar := range envVars {
		if value := os.Getenv(envVar); value != "" {
			return value == "true" || value == "1"
		}
	}

	return defaultValue
}

func getIntValue(configValue types.Int64, envVars []string, defaultValue int) int {
	if !configValue.IsNull() && !configValue.IsUnknown() {
		return int(configValue.ValueInt64())
	}

	for _, envVar := range envVars {
		if value := os.Getenv(envVar); value != "" {
			if intValue, err := strconv.Atoi(value); err == nil {
				return intValue
			}
		}
	}

	return defaultValue
}

func getSDKVersion() string {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return version.Version
	}

	for _, v := range buildInfo.Deps {
		if v.Path == "github.com/hashicorp/terraform-plugin-framework" {
			return v.Version
		}
	}

	return version.Version
}
