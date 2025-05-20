package adapter

import (
	"fmt"
	"os"
	"reflect"
	"slices"
	"strings"
	"support"
)

var allowedTypes = []string{"string", "int", "uint", "int8", "uint8", "int16", "uint16",
	"int32", "uint32", "int64", "uint64", "byte", "rune", "float32", "float64", "complex64", "complex128", "bool"}

func prepareDeviceInfo(deviceFullInfo map[string]any, args support.ArgsList) map[string]any {
	if !args.Ld && !args.Pd {
		return deviceFullInfo
	}
	deviceInfo, ok := deviceFullInfo[args.DriveIndex].(map[string]any)
	if !ok {
		support.ZabbixErrors.ParamNotFound.Show()
	}
	return deviceInfo
}

func getInfoData(data map[string]any, param string) {
	sParam := strings.Split(param, ".")
	if len(sParam) == 1 {
		showParamValue(data, sParam[0])
	}
	sParamLength := len(sParam)
	mainParam := ""
	lastIdx := -1
	for idx, tParam := range sParam {
		mainParam = strings.Join([]string{mainParam, tParam}, "")
		if support.InMap(data, mainParam) {
			lastIdx = idx
			break
		}
		mainParam = strings.Join([]string{mainParam, "."}, "")
		if support.InMap(data, mainParam) {
			lastIdx = idx
		}
	}
	if lastIdx == -1 {
		support.ZabbixErrors.ParamNotFound.Show()
	}
	lastIdx += 1
	if lastIdx == sParamLength {
		showParamValue(data, mainParam)
	}
	test := reflect.TypeOf(data[mainParam])
	if test.Kind() == reflect.Map {
		secondParam := strings.Join(sParam[lastIdx:], ".")
		showParamValue(data[mainParam].(map[string]any), secondParam)
	}
	support.ZabbixErrors.ParamNotFound.Show()
}

func showParamValue(data map[string]any, param string) {
	if checkInfoData(data, param) {
		fmt.Print(data[param])
		os.Exit(0)
	}
	support.ZabbixErrors.ParamNotFound.Show()
}

func checkInfoData(info map[string]any, param string) bool {
	tmpInfo, ok := info[param]
	if !ok {
		return false
	}
	return slices.Contains(allowedTypes, reflect.TypeOf(tmpInfo).Name())
}

func Info(args support.ArgsList) {
	var deviceType support.DeviceInfoFields
	if args.Ld {
		deviceType = support.Devices.Ld
	} else if args.Pd {
		deviceType = support.Devices.Pd
	} else {
		deviceType = support.Devices.Controller
	}
	controllerInfo, ok := support.GetAdapterInfo(args.Index)
	if !ok {
		support.ZabbixErrors.ControllerNotFound.Show()
	}
	deviceFullInfo, ok := controllerInfo[deviceType.Name]
	if !ok {
		support.ZabbixErrors.ParamNotFound.Show()
	}
	getInfoData(prepareDeviceInfo(deviceFullInfo.(map[string]any), args), args.Field)
}
