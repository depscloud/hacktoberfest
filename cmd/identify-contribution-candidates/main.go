package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/depscloud/api/v1alpha/schema"
	"github.com/depscloud/api/v1alpha/tracker"
	"github.com/depscloud/hacktoberfest/internal/config"
	"github.com/depscloud/hacktoberfest/internal/depscloud"
	"github.com/depscloud/hacktoberfest/internal/librariesio"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getEnvOrDefault(name, def string) string {
	if val := os.Getenv(name); val != "" {
		return val
	}
	return def
}

type scoredRepository struct {
	RepositoryURL string `json:"repository_url"`
	Score         int    `json:"score"`
}

func key(module *schema.Module) string {
	lang := module.GetLanguage()
	if n := module.GetName(); n != "" {
		return fmt.Sprintf("%s/%s", lang, n)
	}
	return fmt.Sprintf("%s/%s--%s", lang, module.GetOrganization(), module.GetModule())
}

func scoreTree(root string, edges map[string]map[string]bool, counts map[string]int) int {
	seen := map[string]bool{root: true}
	sum := counts[root]
	tier := []string{root}

	for length := len(tier); length > 0; length = len(tier) {
		next := make([]string, 0)

		for i := 0; i < length; i++ {
			current := tier[i]

			for edge := range edges[current] {
				if !seen[edge] {
					seen[edge] = true
					sum += counts[edge]
					next = append(next, edge)
				}
			}
		}

		tier = next
	}

	return sum
}

func languageToPlatform(language string) string {
	switch language {
	case "java":
		return "maven"
	case "node":
		return "npm"
	case "php":
		return "packagist"
	default:
		return language
	}
}

func main() {
	configFile := getEnvOrDefault("CONFIG_FILE", "config.yaml")
	apiKeyLibrariesIO := getEnvOrDefault("LIBRARIESIO_API_KEY", "")
	outputFile := getEnvOrDefault("OUTPUT_FILE", "candidate.json")

	cfg, err := config.Load(configFile)
	fatal(err)

	conn, err := depscloud.Connect()
	fatal(err)
	defer conn.Close()

	librariesioClient := librariesio.NewClient(apiKeyLibrariesIO)

	moduleService := tracker.NewModuleServiceClient(conn)
	dependencyService := tracker.NewDependencyServiceClient(conn)

	ctx := context.Background()
	count := 100

	index := make(map[string]*schema.Module)
	counts := make(map[string]int)
	edges := make(map[string]map[string]bool)

	for page := int32(1); true; page++ {
		resp, err := moduleService.List(ctx, &tracker.ListRequest{
			Page:  page,
			Count: int32(count),
		})

		if err != nil {
			log.Println(err)
			break
		}

		modules := resp.GetModules()
		for _, module := range modules {
			k1 := key(module)
			index[k1] = module

			// process all modules
			log.Println("processing", module)

			resp, err := dependencyService.ListDependents(ctx, &tracker.DependencyRequest{
				Language:     module.GetLanguage(),
				Organization: module.GetOrganization(),
				Module:       module.GetModule(),
				Name:         module.GetName(),
			})
			if err != nil {
				log.Println("error", err)
				continue
			}

			dependents := resp.GetDependents()
			edges[k1] = make(map[string]bool)
			for _, dependent := range dependents {
				dependentModule := dependent.GetModule()
				k2 := key(dependentModule)

				if cfg.IsCompanyModule(dependentModule) {
					counts[k1]++
				}
				edges[k1][k2] = true
			}
		}

		if len(modules) < count {
			break
		}
	}

	scores := make(map[string]int)
	for key := range edges {
		module := index[key]

		// don't score company modules
		if cfg.IsCompanyModule(module) {
			continue
		}

		log.Println("computing subtree", module)

		score := scoreTree(key, edges, counts)

		if score == 0 {
			continue
		}

		scores[key] = score
	}

	resultsIndex := make(map[string]int)

	// query libraries io @ 1qps
	for key, score := range scores {
		module := index[key]

		log.Println("lookup", module)

		platform := languageToPlatform(module.GetLanguage())
		for attempt := 0; attempt < 5; attempt++ {
			result, err := librariesioClient.LookUp(platform, module.GetName())
			if err != nil {
				log.Println("retrying", err)
				time.Sleep(5 * time.Second)
				continue
			} else if result.RepositoryURL != "" {
				resultsIndex[result.RepositoryURL] += score
				break
			}
		}

		time.Sleep(time.Second)
	}

	results := make([]*scoredRepository, 0, len(scores))
	for url, score := range resultsIndex {
		log.Println("sum", url, score)

		results = append(results, &scoredRepository{
			RepositoryURL: url,
			Score:         score,
		})
	}

	data, err := json.MarshalIndent(results, "", "  ")
	fatal(err)

	err = ioutil.WriteFile(outputFile, data, 0644)
	fatal(err)
}
