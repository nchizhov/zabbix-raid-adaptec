package support

type DeviceInfoFields struct {
	Name  string
	Field string
}

type deviceInfo struct {
	Ld         DeviceInfoFields
	Pd         DeviceInfoFields
	Controller DeviceInfoFields
}

var Devices = deviceInfo{
	Ld: DeviceInfoFields{
		Name:  "Logical device information",
		Field: "Logical device name",
	},
	Pd: DeviceInfoFields{
		Name:  "Physical Device information",
		Field: "Reported Location",
	},
	Controller: DeviceInfoFields{
		Name:  "Controller information",
		Field: "",
	},
}
