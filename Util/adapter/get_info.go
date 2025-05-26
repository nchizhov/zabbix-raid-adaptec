package adapter

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"support"
)

func adapterData(configUtil string, index uint) (string, error) {
	return support.RunCommand(configUtil, "GETCONFIG", strconv.FormatUint(uint64(index), 10), "AL")
}

func parseData(data []string) map[string]any {
	titleSplit := "----------------------------------------------------------------------"
	subtitleSplit := "--------------------------------------------------------"
	logicalDeviceRe, _ := regexp.Compile(`^Logical device number (\d+)$`)
	physicalDeviceRe, _ := regexp.Compile(`^Device #(\d+)$`)

	controllerData := make(map[string]any)
	currentSection := support.NewNil()
	currentSubsection := support.NewNil()
	currentIndex := support.NewNil()
	for idx, line := range data {
		if idx < 3 {
			continue
		}
		line = strings.TrimSpace(line)
		if line == titleSplit && strings.TrimSpace(data[idx-2]) == titleSplit {
			currentSection = support.NewData(strings.TrimSpace(data[idx-1]))
			currentSubsection = support.NewNil()
			currentIndex = support.NewNil()
			support.AddDeepNested(controllerData, []string{currentSection.Value().(string)}, make(map[string]any))
			continue
		}
		if line == subtitleSplit && strings.TrimSpace(data[idx-2]) == subtitleSplit {
			currentSubsection = support.NewData(strings.TrimSpace(data[idx-1]))
			sectionValue := currentSection.Value().(string)
			subsectionValue := currentSubsection.Value().(string)
			if currentIndex.Value() == nil {
				support.AddDeepNested(controllerData,
					[]string{sectionValue, subsectionValue}, make(map[string]any))
			} else {
				support.AddDeepNested(controllerData,
					[]string{sectionValue, currentIndex.Value().(string), subsectionValue}, make(map[string]any))
			}
			continue
		}
		sectionValue := currentSection.Value().(string)
		logicalDevice := logicalDeviceRe.FindSubmatch([]byte(line))
		if logicalDevice != nil {
			indexValue := string(logicalDevice[1])
			currentIndex = support.NewData(indexValue)
			currentSubsection = support.NewNil()
			support.AddDeepNested(controllerData, []string{sectionValue, indexValue}, make(map[string]any))
			continue
		}
		physicalDevice := physicalDeviceRe.FindSubmatch([]byte(line))
		if physicalDevice != nil {
			indexValue := string(physicalDevice[1])
			currentIndex = support.NewData(indexValue)
			currentSubsection = support.NewNil()
			support.AddDeepNested(controllerData,
				[]string{sectionValue, indexValue, "Device Info"}, strings.TrimSpace(data[idx+1]))
			continue
		}
		sLine := strings.SplitN(line, ": ", 2)
		if len(sLine) != 2 {
			continue
		}
		field := strings.TrimSpace(sLine[0])
		fieldVal := strings.TrimSpace(sLine[1])
		if currentIndex.Value() == nil {
			if currentSubsection.Value() == nil {
				support.AddDeepNested(controllerData,
					[]string{sectionValue, field}, fieldVal)
			} else {
				support.AddDeepNested(controllerData,
					[]string{sectionValue, currentSubsection.Value().(string), field}, fieldVal)
			}
			continue
		}
		indexValue := currentIndex.Value().(string)
		if currentSubsection.Value() == nil {
			support.AddDeepNested(controllerData,
				[]string{sectionValue, indexValue, field}, fieldVal)
		} else {
			support.AddDeepNested(controllerData,
				[]string{sectionValue, indexValue, currentSubsection.Value().(string), field}, fieldVal)
		}
	}
	return controllerData
}

func saveInfo(index uint, data map[string]any) {
	infoPath := path.Join(os.TempDir(), fmt.Sprintf("adaptec-%d.json", index))
	if support.FileExists(infoPath) {
		err := os.Remove(infoPath)
		if err != nil {
			return
		}
	}
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return
	}
	err1 := os.WriteFile(infoPath, jsonStr, 0644)
	if err1 != nil {
		return
	}
}

func GetInfo(args support.ArgsList) {
	config, err := support.ReadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	aData, err := adapterData(config.ArcconfPath, args.Index)
	if err != nil {
		log.Fatalln(err)
	}

	data := strings.Split(aData, "\n")
	controllerData := parseData(data)

	saveInfo(args.Index, controllerData)
	os.Exit(0)
}
