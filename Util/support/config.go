package support

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"
)

type ConfigInfo struct {
	ArcconfPath string `json:"arcconf_path"`
}

func ReadConfig() (ConfigInfo, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ConfigInfo{}, err
	}
	configFile, err := os.ReadFile(path.Join(dir, "config.json"))
	if err != nil {
		return ConfigInfo{}, err
	}
	jsonDec := json.NewDecoder(bytes.NewReader(configFile))
	jsonDec.DisallowUnknownFields()
	for {
		var config ConfigInfo
		if err := jsonDec.Decode(&config); err == io.EOF {
			break
		} else if err != nil {
			return ConfigInfo{}, err
		}
		return config, nil
	}
	return ConfigInfo{}, err
}

func GetAdapterInfo(index uint) (map[string]any, bool) {
	infoPath := path.Join(os.TempDir(), fmt.Sprintf("adaptec-%d.json", index))
	var data map[string]any
	retries := 0
	for retries < 4 {
		retries++
		if !FileExists(infoPath) {
			time.Sleep(time.Second)
			continue
		}
		fData, err := os.ReadFile(infoPath)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		jsonErr := json.Unmarshal(fData, &data)
		if jsonErr != nil || data == nil {
			time.Sleep(time.Second)
			continue
		}
		break
	}
	return data, data != nil
}
