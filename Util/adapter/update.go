package adapter

import (
	"elog"
	"fmt"
	"log"
	"os"
	"support"
)

func Update() {
	releaseURL := "https://api.github.com/repos/nchizhov/zabbix-raid-adaptec/releases/latest"
	programName := "adapter"

	log.SetFlags(log.Ldate | log.Ltime)
	elog.Info("Start check update")
	releaseInfo, err := support.GetUpdateInfo(releaseURL, programName)
	if err != nil {
		elog.Fatal(err)
	}

	programMD5, err := support.GetProgramChecksum()
	if err != nil {
		elog.Fatal(err)
	}
	if programMD5 == releaseInfo.MD5 {
		elog.Info("No new release was found")
	} else {
		err = support.SelfUpdate(releaseInfo)
		if err != nil {
			elog.Fatal(err)
		}
		elog.Info(fmt.Sprintf("Successfully updated: %s", releaseInfo.Name))
	}
	elog.Info("Finish check update")
	os.Exit(1)
}
