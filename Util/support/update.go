package support

import (
	"bytes"
	"crypto"
	"crypto/md5"
	"elog"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/minio/selfupdate"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

type ReleaseInfo struct {
	Name string
	Url  string
	MD5  string
}

func (releaseInfo *ReleaseInfo) fillReleaseInfo(data []byte, prefix string) error {
	var dataInfo map[string]any
	err := json.Unmarshal(data, &dataInfo)
	if err != nil {
		return err
	}
	releaseInfo.Name = getRequiredReleaseName(dataInfo["tag_name"].(string), prefix)
	assets, ok := dataInfo["assets"]
	if !ok {
		return errors.New("assets in release not found")
	}
	err = releaseInfo.parseReleaseInfo(assets.([]any))
	if err != nil {
		return err
	}
	return nil
}

func (releaseInfo *ReleaseInfo) parseReleaseInfo(assetInfo []any) error {
	hasMD5 := false
	hasRelease := false
	for _, baseAsset := range assetInfo {
		asset := baseAsset.(map[string]any)
		assetName := asset["name"].(string)
		if !strings.HasPrefix(assetName, releaseInfo.Name) {
			continue
		}
		url := asset["browser_download_url"].(string)
		if strings.HasSuffix(assetName, ".md5") {
			md5Hash, err := downloadFile(url)
			if err != nil {
				return err
			}
			releaseInfo.MD5 = strings.TrimSpace(string(md5Hash))
			hasMD5 = true
			continue
		}
		releaseInfo.Url = url
		hasRelease = true
	}
	if !hasRelease {
		return errors.New("release file not exists")
	}
	if !hasMD5 {
		return errors.New("md5 checksum not exists")
	}
	return nil
}

func getRequiredReleaseName(tag string, prefix string) string {
	postfix := ""
	if runtime.GOOS == "windows" {
		postfix = ".exe"
	}
	return fmt.Sprintf("%s-%s-%s-%s%s", prefix, tag, runtime.GOOS, runtime.GOARCH, postfix)
}

func getGithubData(url string) ([]byte, error) {
	elog.Info(fmt.Sprintf("Check new release at %s", url))

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := readDownloadedData(resp, "incorrect status code for latest release info: %d")
	if err != nil {
		return nil, err
	}
	return body, nil
}

func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := readDownloadedData(resp, "incorrect status code for MD5: %d")
	if err != nil {
		return nil, err
	}
	return body, nil
}

func readDownloadedData(resp *http.Response, errFormat string) ([]byte, error) {
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(errFormat, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func doSelfUpdate(url string, md5Hash []byte) error {
	data, err := downloadFile(url)
	if err != nil {
		return err
	}
	return selfupdate.Apply(bytes.NewReader(data), selfupdate.Options{
		Hash:     crypto.MD5,
		Checksum: md5Hash,
	})
}

func GetUpdateInfo(url string, prefix string) (ReleaseInfo, error) {
	var releaseInfo ReleaseInfo

	data, err := getGithubData(url)
	if err != nil {
		return releaseInfo, err
	}

	elog.Info("Parse release info")
	err = releaseInfo.fillReleaseInfo(data, prefix)
	if err != nil {
		return releaseInfo, err
	}
	elog.Info(fmt.Sprintf("Latest release info: %s", releaseInfo.Name))
	return releaseInfo, nil
}

func GetProgramChecksum() (string, error) {
	fl, err := os.ReadFile(os.Args[0])
	if err != nil {
		return "", err
	}
	md5SumTmp := md5.Sum(fl)
	md5Sum := hex.EncodeToString(md5SumTmp[:])
	elog.Info(fmt.Sprintf("Current program MD5-sum: %s", md5Sum))
	return md5Sum, nil
}

func SelfUpdate(releaseInfo ReleaseInfo) error {
	md5Hash, err := hex.DecodeString(releaseInfo.MD5)
	if err != nil {
		return err
	}
	return doSelfUpdate(releaseInfo.Url, md5Hash)
}
