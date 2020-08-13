package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"time"

	"github.com/depscloud/api/v1alpha/schema"
)

type LibrariesIOResult struct {
	HomePage string `json:"homepage"`
	RepositoryURL string `json:"repository_url"`
}

type result struct {
	Module *schema.Module `json:"module"`
	Score  int            `json:"score"`
}

type output struct {
	URL string `json:"url"`
	Score int `json:"score"`
}

func formatName(module *schema.Module) (string, string) {
	switch module.GetLanguage() {
	case "java":
		return "maven", module.GetOrganization() + ":" + module.GetModule()
	case "node":
		name := module.GetModule()
		if org := module.GetOrganization(); org != "" {
			name = "@" + org + "/" + name
		}
		return "npm", name
	}
	return "", ""
}

func fetchRepoInfo(apiKey, platform, name string) *LibrariesIOResult {
	uri := fmt.Sprintf("https://libraries.io/api/%s/%s?api_key=%s",
		url.QueryEscape(platform), url.QueryEscape(name), url.QueryEscape(apiKey))

	resp, err := http.Get(uri)
	if err != nil {
		log.Println(err)
		return nil
	}


	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}

	result := &LibrariesIOResult{}
	if err := json.Unmarshal(data, result); err != nil {
		log.Println(string(data))
		return nil
	}

	return result
}

func main() {
	inputFile := "libraries.json"
	outputFile := "repositories.json"
	apiKeyLibrariesIO := os.Getenv("LIBRARIESIO_API_KEY")

	rawInputData, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	inputData := make([]*result, 0)
	if err := json.Unmarshal(rawInputData, &inputData); err != nil {
		log.Fatal(err)
	}

	idx := make(map[string]int)
	for _, input := range inputData {
		start := time.Now()

		platform, name := formatName(input.Module)
		if name == "" {
			continue
		}

		log.Print("processing ", name)

		info := fetchRepoInfo(apiKeyLibrariesIO, platform, name)
		if info == nil {
			continue
		}


		repoURL := info.RepositoryURL
		if repoURL == "" {
			repoURL = info.HomePage
		}

		if _, ok := idx[repoURL]; !ok {
			idx[repoURL] = 0
		}

		idx[repoURL] += input.Score

		// throttle
		remaining := time.Second - time.Now().Sub(start)
		if remaining > 0 {
			time.Sleep(remaining)
		}
	}

	scores := make([]output, len(idx))
	i := 0

	for repoURL, score := range idx {
		scores[i] = output{
			URL: repoURL,
			Score: score,
		}
		i++
	}

	sort.SliceStable(scores, func(i, j int) bool {
		return scores[i].Score >= scores[j].Score
	})

	data, err := json.MarshalIndent(scores, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(outputFile, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
