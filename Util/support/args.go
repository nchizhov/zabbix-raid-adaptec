package support

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
)

type ArgsList struct {
	Index      uint
	Update     bool
	GetInfo    bool
	Discovery  bool
	Info       bool
	Pd         bool
	Ld         bool
	DriveIndex string
	Field      string
}

func (args ArgsList) validateMainFields() error {
	if !args.Info && !args.GetInfo && !args.Discovery && !args.Update {
		return errors.New("one of update, info, get-info, discovery is required")
	}
	if args.Info && args.GetInfo && args.Discovery && args.Update {
		return errors.New("update, info, get-info and discovery flags cannot used at same time")
	}
	if (!args.Update && !args.Info && !args.GetInfo && args.Discovery) ||
		(!args.Update && !args.Info && !args.Discovery && args.GetInfo) ||
		(!args.Update && !args.GetInfo && !args.Discovery && args.Info) ||
		(!args.Info && !args.GetInfo && !args.Discovery && args.Update) {
		return nil
	}
	return errors.New("update, info, get-info or discovery flags cannot be combined")
}

func (args ArgsList) validateDriveFields() {
	if (!args.Pd && !args.Ld) || (args.Pd && args.Ld) {
		ZabbixErrors.Arguments.Show()
	}
}

func (args ArgsList) validateInfoFields() {
	if args.Ld || args.Pd {
		args.validateDriveFields()
		if args.DriveIndex == "" {
			ZabbixErrors.Arguments.Show()
		}
	} else {
		if args.DriveIndex != "" {
			ZabbixErrors.Arguments.Show()
		}
	}
	if args.Field == "" {
		ZabbixErrors.Arguments.Show()
	}
}

func (args ArgsList) validateArgs() {
	err := args.validateMainFields()
	if err != nil {
		if !args.GetInfo && !args.Update {
			ZabbixErrors.Arguments.Show()
		}
		log.Fatalln(err)
	}
	if args.Discovery {
		args.validateDriveFields()
		return
	}
	if args.Info {
		args.validateInfoFields()
	}
}

func showHelpBanner() {
	fmt.Println("Adaptec RAID-Controller Utility for Zabbix")
	fmt.Println()
	fmt.Println("(c) 2025 Chizhov Nikolay <nchizhov@inok.ru>")
	fmt.Println()
	flag.PrintDefaults()
}

func ParseArgs() ArgsList {
	var args ArgsList

	isHelp := flag.Bool("help", false, "Show help")

	flag.UintVar(&args.Index, "index", 1, "RAID-Controller index")
	flag.BoolVar(&args.Update, "update", false, "Self-update")
	flag.BoolVar(&args.GetInfo, "get-info", false, "Get info of RAID-Controller")
	flag.BoolVar(&args.Discovery, "discovery", false, "Zabbix discovery for Logical/Physical Drives")
	flag.BoolVar(&args.Info, "info", false, "RAID-Controller or Logical/Physical Drive field info for Zabbix")
	flag.BoolVar(&args.Pd, "pd", false, "Physical Drive")
	flag.BoolVar(&args.Ld, "ld", false, "Logical Drive")
	flag.StringVar(&args.DriveIndex, "drive-index", "", "Logical/Physical Drive index")
	flag.StringVar(&args.Field, "field", "", "Field name for Zabbix. Use \".\" for subfields")

	flag.Parse()

	if *isHelp {
		showHelpBanner()
		os.Exit(0)
	}

	args.validateArgs()
	return args
}
