package resolvers

import (
	"context"

	"github.com/depscloud/api/v1alpha/schema"
)

// OSS fetches additional metadata for a module
type OSS struct {
	Scorer *Score
	Lookup *URL
}

// Resolve attempts to get additional information for a module
func (o *OSS) Resolve(ctx context.Context, module *schema.Module) (int, string) {
	score := o.Scorer.Resolve(ctx, module)
	if score == 0 {
		return 0, ""
	}

	url := o.Lookup.Resolve(ctx, module)
	if url == "" {
		return 0, ""
	}

	return score, url
}
