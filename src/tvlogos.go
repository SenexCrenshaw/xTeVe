package src

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type GitFileFolder struct {
	FileFolder string `json:"name"`
	Url        string `json:"url"`
	IsType     string `json:"type"`
}

// We don't want to store in our JSON any files or folders with these substrings in the name
func ignoredFiles(fileFolderName string) bool {
	var ignoredFilesMap = map[int]string{
		0:  ".md",
		1:  "license",
		2:  "demo",
		3:  "how-to",
		4:  "cname",
		5:  "github",
		6:  "new-sky-sports-logos",
		7:  "paypal",
		8:  "sponser-button",
		9:  "terrorcon-",
		10: "repository-open",
		11: ".yml",
		12: "foo.txt",
		13: "sponsor-button",
		14: "%CE%A9",
	}
	for _, s := range ignoredFilesMap {
		if strings.Contains(strings.ToLower(fileFolderName), strings.ToLower(s)) {
			return true
		}
	}
	return false
}

func contains(elems []html.Attribute, v string) bool {
	for _, s := range elems {
		if strings.Contains(strings.ToLower(s.Val), v) {
			return true
		}
	}
	return false
}

func gitFileFolders(doc *html.Node) map[int]GitFileFolder {
	i := 0
	var crawler func(*html.Node)
	var list = map[int]GitFileFolder{}
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			var isTypeOf string
			if contains(node.Attr, "#repo-content-pjax-container") && contains(node.Attr, "js-navigation-open link--primary") && !ignoredFiles(node.Attr[1].Val) {
				if strings.Contains(node.Attr[1].Val, ".png") || strings.Contains(node.Attr[1].Val, ".jpg") {
					isTypeOf = "Image"
				} else {
					isTypeOf = "Folder"
				}
				gitFileFolder := GitFileFolder{FileFolder: node.Attr[1].Val, Url: node.Attr[4].Val, IsType: isTypeOf}
				list[i] = gitFileFolder
				i++
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)
	return list
}

func fetchRepositoryData(u string) (string, error) {
	response, err := http.Get(u)
	if err != nil {
		errMsg := errors.New("unable to fetch git repository URL: " + err.Error())
		ShowError(errMsg, 000)
		return "", errMsg
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		errMsg := errors.New("error in reading git repository response body: " + err.Error())
		ShowError(errMsg, 000)
		return "", errMsg
	}
	data := string(body)
	if strings.Contains(data, "#repo-content-pjax-container") {
		return data, nil
	}
	return "", errors.New("did not find repository file or folder data: " + u)
}

func buildLogos() {

	var logoJsonExists bool

	fileStats, err := os.Stat(System.File.TVLogos)

	if err == nil {
		logoJsonExists = true
	} else {
		logoJsonExists = false
	}

	showInfo("TV Logos:tvlogos.json file exists " + strconv.FormatBool(logoJsonExists))
	// This should only be updated once a week or if the file does not exist
	if !logoJsonExists || (logoJsonExists && time.Now().After(fileStats.ModTime().Add(time.Hour*168))) {
		showInfo("TV Logos:File does not exist or time for update")
		var baseUri string
		var gitCrawl func(string)
		masterList := make(map[string][]string)

		baseUri = "https://github.com/"

		gitCrawl = func(uri string) {
			data, err := fetchRepositoryData(uri)
			path := strings.TrimPrefix(uri, "https://github.com/")
			logos := []string{}
			if err != nil {
				return
			}
			if strings.Contains(data, "#repo-content-pjax-container") {
				doc, _ := html.Parse(strings.NewReader(data))
				paths := gitFileFolders(doc)
				for _, f := range paths {
					if f.IsType == "Folder" && !ignoredFiles(f.FileFolder) {
						gitCrawl(baseUri + f.Url)
					} else if !ignoredFiles(f.FileFolder) {
						logos = append(logos, f.FileFolder)
						masterList[path] = logos
					}
				}
			}
		}

		data, err := fetchRepositoryData(baseUri + "Tapiosinn/tv-logos/tree/master/")

		if err != nil {
			ShowError(err, 0)
		}
		if strings.Contains(data, "#repo-content-pjax-container") {
			logos := []string{}
			doc, _ := html.Parse(strings.NewReader(data))
			paths := gitFileFolders(doc)
			for _, f := range paths {
				if f.IsType == "Folder" && !ignoredFiles(f.FileFolder) {
					gitCrawl(baseUri + f.Url)
				} else if !ignoredFiles(f.FileFolder) {
					logos = append(logos, f.FileFolder)
					masterList["Tapiosinn/tv-logos/tree/master/"] = logos
				}
			}
		}
		jsonData, _ := json.MarshalIndent(masterList, "", "\t")
		file, err := os.Create(getPlatformFile(System.Folder.Config + "tvlogos.json"))
		if err != nil {
			ShowError(err, 000)
			return
		}
		file.Close()
		err = ioutil.WriteFile(getPlatformFile(System.Folder.Config+"tvlogos.json"), jsonData, 0644)
		if err != nil {
			ShowError(err, 000)
			return
		}

		Data.TVLogos.Files = make(map[string]interface{})

		Data.TVLogos.Files, err = loadJSONFileToMap(System.File.TVLogos)
		showInfo("TV Logos:Done updating TV logos file")
		if err != nil {
			ShowError(err, 000)
		}
	} else if logoJsonExists {
		showInfo("TV Logos:Loading TV logos file")
		Data.TVLogos.Files = make(map[string]interface{})

		Data.TVLogos.Files, err = loadJSONFileToMap(System.File.TVLogos)

		if err != nil {
			ShowError(err, 000)
		}
	}
}
