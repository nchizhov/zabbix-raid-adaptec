package adapter

import (
	"encoding/json"
	"fmt"
	"os"
	"support"
)

type discoveryInfo struct {
	Data []map[string]string `json:"data"`
}

func getDiscoveryDevices(deviceType support.DeviceInfoFields, info map[string]any) []map[string]string {
	var discoveryData []map[string]string
	deviceInfo, ok := info[deviceType.Name]
	if !ok {
		return discoveryData
	}
	for field, fValue := range deviceInfo.(map[string]any) {
		fieldInfo := fValue.(map[string]any)
		fieldData, ok := fieldInfo[deviceType.Field]
		if !ok {
			continue
		}
		tmpData := map[string]string{
			"{#DEVICEID}":   field,
			"{#DEVICENAME}": fieldData.(string),
		}
		discoveryData = append(discoveryData, tmpData)
	}
	return discoveryData
}

func Discovery(args support.ArgsList) {
	var deviceType support.DeviceInfoFields
	if args.Ld {
		deviceType = support.Devices.Ld
	} else if args.Pd {
		deviceType = support.Devices.Pd
	} else {
		support.ZabbixErrors.DeviceTypeNotFound.Show()
	}
	controllerInfo, ok := support.GetAdapterInfo(args.Index)
	if !ok {
		support.ZabbixErrors.ControllerNotFound.Show()
	}
	discovery := discoveryInfo{
		Data: getDiscoveryDevices(deviceType, controllerInfo),
	}
	jsonData, jErr := json.Marshal(discovery)
	if jErr != nil {
		support.ZabbixErrors.Unknown.Show()
	}
	fmt.Print(string(jsonData))
	os.Exit(0)
}
