package resolvers

import (
	"context"
	"log"

	"github.com/depscloud/api/v1alpha/schema"
	"github.com/depscloud/hacktoberfest/internal/librariesio"
)

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
	case "go":
		return "go", module.GetOrganization() + "/" + module.GetModule()
	}
	return "", ""
}

type URL struct {
	LibrariesIO *librariesio.Client
}

func (u *URL) Resolve(ctx context.Context, module *schema.Module) string {
	platform, name := formatName(module)
	if name == "" {
		return ""
	}

	result, err := u.LibrariesIO.LookUp(platform, name)
	if err != nil {
		log.Println(err)
		return ""
	}

	repoURL := result.RepositoryURL
	if repoURL == "" {
		repoURL = result.HomePage
	}

	return repoURL
}
