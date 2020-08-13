package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/depscloud/api/v1alpha/schema"
	"github.com/depscloud/api/v1alpha/tracker"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	VariableBaseURL = "DEPSCLOUD_BASE_URL"

	DefaultBaseURL = "https://api.deps.cloud"
)

func translateBaseURL(baseURL string) (bool, string) {
	tls := false
	uri, _ := url.Parse(baseURL)

	if uri.Scheme == "https" {
		tls = true
	}

	host := uri.Host
	if !strings.Contains(host, ":") {
		if tls {
			host = host + ":443"
		} else {
			host = host + ":80"
		}
	}

	return tls, host
}

type result struct {
	Module schema.Module `json:"module"`
	Score  int           `json:"score"`
}

func main() {
	// todo: populate
	outputFile := "libraries.json"
	languages := map[string]bool{}
	corporateOrganizations := map[string]bool{}

	baseURL := DefaultBaseURL
	if val := os.Getenv(VariableBaseURL); val != "" {
		baseURL = val
	}

	isSecure, target := translateBaseURL(baseURL)

	options := make([]grpc.DialOption, 0)
	if isSecure {
		creds := credentials.NewTLS(&tls.Config{})
		options = append(options, grpc.WithTransportCredentials(creds))
	} else {
		options = append(options, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(target, options...)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	moduleService := tracker.NewModuleServiceClient(conn)
	searchService := tracker.NewSearchServiceClient(conn)

	ctx := context.Background()
	count := 100

	openSourceLibraries := make([]*schema.Module, 0)

	for page := int32(1); true; page++ {
		resp, err := moduleService.List(ctx, &tracker.ListRequest{
			Page:  page,
			Count: int32(count),
		})

		if err != nil {
			log.Fatal(err)
		}

		for _, module := range resp.GetModules() {
			lang := module.GetLanguage()
			if _, ok := languages[lang]; !ok {
				continue
			}

			if _, ok := corporateOrganizations[module.GetOrganization()]; !ok {
				openSourceLibraries = append(openSourceLibraries, module)
			}
		}

		if len(resp.GetModules()) < count {
			// no more left
			break
		}
	}

	log.Printf("discovered %d open source libraries in use", len(openSourceLibraries))

	scores := make([]result, len(openSourceLibraries))
	for i, ossLibrary := range openSourceLibraries {
		log.Printf("processing (%s)\n", ossLibrary)

		call, err := searchService.BreadthFirstSearch(ctx)
		if err != nil {
			log.Fatal(err)
		}

		err = call.Send(&tracker.SearchRequest{
			DependentsOf: &tracker.DependencyRequest{
				Language:     ossLibrary.GetLanguage(),
				Organization: ossLibrary.GetOrganization(),
				Module:       ossLibrary.GetModule(),
			},
		})

		score := 0
		for msg, err := call.Recv(); err == nil; msg, err = call.Recv() {
			for _, dependentModule := range msg.GetDependents() {
				if _, ok := corporateOrganizations[dependentModule.GetModule().GetOrganization()]; ok {
					score++
				}
			}
		}

		scores[i] = result{
			Module: *ossLibrary,
			Score:  score,
		}
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
