package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sort"

	"github.com/depscloud/api/v1alpha/tracker"
	"github.com/depscloud/hacktoberfest/internal/config"
	"github.com/depscloud/hacktoberfest/internal/depscloud"
	"github.com/depscloud/hacktoberfest/internal/librariesio"
	"github.com/depscloud/hacktoberfest/internal/resolvers"
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

func scoreAllModules(moduleService tracker.ModuleServiceClient, cfg *config.Config, oss *resolvers.OSS) map[string]int {
	ctx := context.Background()
	count := 100

	scores := make(map[string]int)
	for page := int32(1); true; page++ {
		resp, err := moduleService.List(ctx, &tracker.ListRequest{
			Page:  page,
			Count: int32(count),
		})

		if err != nil {
			log.Println(err)
			break
		}

		for _, module := range resp.GetModules() {
			if cfg.IsCompanyModule(module) {
				log.Println("filtering company module", module)
				continue
			}

			score, url := oss.Resolve(ctx, module)
			if url == "" {
				continue
			}

			log.Println("updating", url)

			if _, ok := scores[url]; !ok {
				scores[url] = 0
			}

			scores[url] += score
		}
	}

	return scores
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

	moduleService := tracker.NewModuleServiceClient(conn)

	oss := &resolvers.OSS{
		Scorer: &resolvers.Score{
			SearchService: tracker.NewSearchServiceClient(conn),
			Config:        cfg,
		},
		Lookup: &resolvers.URL{
			LibrariesIO: librariesio.NewClient(apiKeyLibrariesIO),
		},
	}

	scores := scoreAllModules(moduleService, cfg, oss)

	rankedScores := make([]*scoredRepository, 0, len(scores))
	for url, score := range scores {
		rankedScores = append(rankedScores, &scoredRepository{
			RepositoryURL: url,
			Score:         score,
		})
	}

	sort.SliceStable(rankedScores, func(i, j int) bool {
		return rankedScores[i].Score >= rankedScores[j].Score
	})

	data, err := json.MarshalIndent(rankedScores, "", "  ")
	fatal(err)

	err = ioutil.WriteFile(outputFile, data, 0644)
	fatal(err)
}
