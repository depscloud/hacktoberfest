package resolvers

import (
	"context"
	"log"

	"github.com/depscloud/api/v1alpha/schema"
	"github.com/depscloud/api/v1alpha/tracker"
	"github.com/depscloud/hacktoberfest/internal/config"
)

// Score attempts to resolve a numeric score for a module
type Score struct {
	SearchService tracker.SearchServiceClient
	Config        *config.Config
}

// Resolve attempts to get additional information for a module
func (s *Score) Resolve(ctx context.Context, module *schema.Module) int {
	call, err := s.SearchService.BreadthFirstSearch(ctx)
	if err != nil {
		log.Print(err)
		return 0
	}

	err = call.Send(&tracker.SearchRequest{
		DependentsOf: &tracker.DependencyRequest{
			Language:     module.GetLanguage(),
			Organization: module.GetOrganization(),
			Module:       module.GetModule(),
		},
	})
	if err != nil {
		log.Print(err)
		return 0
	}

	score := 0
	for msg, err := call.Recv(); err == nil; msg, err = call.Recv() {
		for _, dependentModule := range msg.GetDependents() {
			if s.Config.IsCompanyModule(dependentModule.GetModule()) {
				score++
			}
		}
	}
	return score
}
