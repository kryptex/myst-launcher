package myst

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/mysteriumnetwork/myst-launcher/gui"
)

func CheckUpdates(imageDigest string) {
	url := "https://registry.hub.docker.com/v2/repositories/mysteriumnetwork/myst/tags?page_size=10"
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		return
	}
	data, _ := ioutil.ReadAll(resp.Body)

	//results
	latestDigest := ""
	latestVersion := ""
	currentVersion := ""

	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		name, err := jsonparser.GetString(value, "name")
		if err != nil {
			return
		}

		jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			digest, err := jsonparser.GetString(value, "digest")
			if err != nil {
				return
			}
			if name == "latest" {
				latestDigest = digest
			}
			match, _ := regexp.MatchString(`^\d+\.\d+\.\d+.*$`, name)
			if match && latestDigest == digest {
				latestVersion = name
			}
			digestsMatch := strings.ToLower(digest) == strings.ToLower(imageDigest)
			if digestsMatch && match {
				currentVersion = name
			}
		}, "images")
	}, "results")

	gui.UI.VersionUpToDate = latestDigest == imageDigest
	gui.UI.VersionCurrent = currentVersion
	gui.UI.VersionLatest = latestVersion
	gui.UI.Update()
}
