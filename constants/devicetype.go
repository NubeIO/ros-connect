package constants

import "fmt"

const Permission = 0755

// If something occurs unusual we do the mappings here
var appNameToServiceNameMap = map[string]string{}

var appNameToRepoNameMap = map[string]string{}

var appNameToDataDirNameMap = map[string]string{}

func GetServiceNameFromAppName(appName string) string {
	if value, found := appNameToServiceNameMap[appName]; found {
		return value
	}
	return fmt.Sprintf("nubeio-%s.service", appName)
}

func GetAppNameFromRepoName(repoName string) string {
	for k := range appNameToRepoNameMap {
		if appNameToRepoNameMap[k] == repoName {
			return k
		}
	}
	return repoName
}

func GetRepoNameFromAppName(appName string) string {
	if value, found := appNameToRepoNameMap[appName]; found {
		return value
	}
	return appName
}

func GetDataDirNameFromAppName(appName string) string {
	if value, found := appNameToDataDirNameMap[appName]; found {
		return value
	}
	return appName
}

type DeviceType int64

const (
	Cloud DeviceType = iota
	Edge28
	RubixCompute
	RubixComputeVPN
	RubixComputeLoRaWAN
	RubixComputeLoRaWANVPN
	RubixComputeIO
)

func ValidRubixCompute(s string) bool {
	switch s {
	case RubixCompute.String():
		return true
	case RubixComputeVPN.String():
		return true
	case RubixComputeLoRaWAN.String():
		return true
	case RubixComputeLoRaWANVPN.String():
		return true
	case RubixComputeIO.String():
		return true
	}
	return false
}

func (s DeviceType) String() string {
	switch s {
	case Cloud:
		return "cloud"
	case Edge28:
		return "edge-28"
	case RubixCompute:
		return "rubix-compute"
	case RubixComputeVPN:
		return "rubix-compute-vpn"
	case RubixComputeLoRaWAN:
		return "rubix-compute-lorawan"
	case RubixComputeLoRaWANVPN:
		return "rubix-compute--lorawan-vpn"
	case RubixComputeIO:
		return "rubix-compute-io"
	}
	return "unknown"
}
