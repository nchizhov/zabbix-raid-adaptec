package support

import (
	"log"
)

type zbxError string

func (zbxErr zbxError) Show() {
	log.Fatalf("ZBX_NOTSUPPORTED\000%s", zbxErr)
}

// Zabbix incorrect errors
type zbxErrors struct {
	Arguments          zbxError
	DeviceTypeNotFound zbxError
	ControllerNotFound zbxError
	ParamNotFound      zbxError
	Unknown            zbxError
}

var ZabbixErrors = zbxErrors{
	Arguments:          "Incorrect arguments",
	DeviceTypeNotFound: "Device type not found. Correct value: ld, pd",
	ControllerNotFound: "Controller device info not found",
	ParamNotFound:      "Param not exists",
	Unknown:            "Unknown error",
}
