package src

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

type TVLogoInformation struct {
	Sha  string `json:"sha"`
	URL  string `json:"url"`
	Tree []struct {
		Path string `json:"path"`
		Mode string `json:"mode"`
		Type string `json:"type"`
		Sha  string `json:"sha"`
		URL  string `json:"url"`
		Size int    `json:"size,omitempty"`
	} `json:"tree"`
	Truncated bool `json:"truncated"`
}

type LogoInformation struct {
	Country  string `json:"country"`
	URL      string `json:"url"`
	FileName string `json:"filename"`
}

func downloadLogoJSON() {

	info, err := os.Stat(System.File.TVLogos)
	logoJsonExists := err == nil

	//var logoInformation []LogoInformation

	if logoJsonExists && time.Now().Before(info.ModTime().Add(time.Hour*168)) {
		content, err := ioutil.ReadFile(System.File.TVLogos)
		if err != nil {
			ShowError(err, 0)
		}

		err = json.Unmarshal(content, &Data.Logos.logoInformation)
		if err != nil {
			log.Fatal("Error during Unmarshal(): ", err)
		}

	} else {
		providerURL := "https://api.github.com/repos/Tapiosinn/tv-logos/git/trees/master?recursive=1"

		resp, err := http.Get(providerURL)
		if err != nil {
			return
		}

		b, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return
		}

		var tvLogoInformation *TVLogoInformation
		err = json.Unmarshal(b, &tvLogoInformation)
		if err != nil {
			return
		}

		rcountry := regexp.MustCompile(`countries\/(?P<country>.*)\/(?P<filename>.*[png|jpg])$`)
		rcmisc := regexp.MustCompile(`misc\/(?P<country>.*)\/(?P<filename>.*[png|jpg])$`)
		rcmisc2 := regexp.MustCompile(`misc\/(.*[png|jpg])$`)

		var match [][]string

		for _, v := range tvLogoInformation.Tree {
			if !strings.HasSuffix(v.Path, "png") {
				continue
			}
			fmt.Printf("%s\n", v.Path)
			if strings.HasPrefix(v.Path, "countries") {
				match = rcountry.FindAllStringSubmatch(v.Path, -1)
				if match == nil {
					fmt.Printf("BAD match %s\n", v.Path)
					continue
				}
			} else {
				match = rcmisc.FindAllStringSubmatch(v.Path, -1)
				if match == nil {
					match = rcmisc2.FindAllStringSubmatch(v.Path, -1)
					if match == nil {
						fmt.Printf("BAD match %s\n", v.Path)
						continue
					}
				}
			}

			if match != nil && match[0] != nil && len(match[0]) != 3 {
				continue
			}

			Data.Logos.logoInformation = append(Data.Logos.logoInformation, LogoInformation{Country: match[0][1], URL: "https://raw.githubusercontent.com/Tapiosinn/tv-logos/master/" + v.Path, FileName: match[0][2]})
		}

		sort.Slice(Data.Logos.logoInformation, func(i, j int) bool {
			return Data.Logos.logoInformation[i].FileName < Data.Logos.logoInformation[j].FileName
		})

		file, _ := json.MarshalIndent(Data.Logos.logoInformation, "", " ")
		_ = ioutil.WriteFile(System.File.TVLogos, file, 0644)

	}
}
