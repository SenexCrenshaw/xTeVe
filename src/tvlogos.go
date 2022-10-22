package src

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

type TVLogoInformation struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Path string `json:"path"`
	Mode string `json:"mode"`
}

type LogoInformation struct {
	Country  string `json:"country"`
	Path     string `json:"path"`
	FileName string `json:"filename"`
}

func downloadLogoJSON() {

	info, err := os.Stat(System.File.TVLogos)
	logoJsonExists := err == nil

	if logoJsonExists && time.Now().Before(info.ModTime().Add(time.Hour*168)) {
		content, err := os.ReadFile(System.File.TVLogos)
		if err != nil {
			ShowError(err, 0)
		}

		err = json.Unmarshal(content, &Data.Logos)
		if err != nil {
			log.Fatal("Error during Unmarshal(): ", err)
		}

	} else {
		page := "1"

		Data.Logos.URL = "https://gitlab.com/tapiosinn/tv-logos/-/raw/main/"

		var tvLogoInformation *[]TVLogoInformation
		rcountry := regexp.MustCompile(`^(countries|misc)\/(?P<country>[a-zA-Z0-9-]+)\/(.*\/)?(?P<filename>.*[png|jpg])$`)
		var match [][]string

		fmt.Println("Building logo information")
		for {

			providerURL := "https://gitlab.com/api/v4/projects/40347226/repository/tree?recursive=true&per_page=100&page=" + page

			resp, err := http.Get(providerURL)
			if err != nil {
				break
			}

			b, err := io.ReadAll(resp.Body)
			defer resp.Body.Close()
			if err != nil {
				break
			}

			err = json.Unmarshal(b, &tvLogoInformation)
			if err != nil {
				break
			}

			for _, v := range *tvLogoInformation {
				if v.Type != "blob" || !strings.HasSuffix(v.Path, "png") {
					continue
				}

				match = rcountry.FindAllStringSubmatch(v.Path, -1)
				if match == nil {
					fmt.Printf("BAD match %s\n", v.Path)
					continue

				}

				if match != nil && match[0] != nil && len(match[0]) != 5 {
					continue
				}

				Data.Logos.LogoInformation = append(Data.Logos.LogoInformation, LogoInformation{Country: match[0][2], Path: v.Path, FileName: match[0][4]})

			}

			if page = resp.Header.Get("X-Next-Page"); page == "" {
				break
			}

			fmt.Println("Page", page, "/", resp.Header.Get("X-Total-Pages"))

		}

		sort.Slice(Data.Logos.LogoInformation, func(i, j int) bool {
			return Data.Logos.LogoInformation[i].FileName < Data.Logos.LogoInformation[j].FileName
		})

		file, _ := json.MarshalIndent(Data.Logos, "", " ")
		_ = os.WriteFile(System.File.TVLogos, file, 0644)

	}
}
