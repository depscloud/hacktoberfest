![GitHub](https://img.shields.io/github/license/depscloud/hacktoberfest.svg)
![branch](https://github.com/depscloud/hacktoberfest/workflows/branch/badge.svg?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/depscloud/hacktoberfest)](https://goreportcard.com/report/github.com/depscloud/hacktoberfest)
![Google Analytics](https://www.google-analytics.com/collect?v=1&cid=555&t=pageview&ec=repo&ea=open&dp=hacktoberfest&dt=hacktoberfest&tid=UA-143087272-2)

**WORK IN PROGRESS**

# deps.cloud

[deps.cloud](https://deps.cloud/) is system that helps track and manage library usage across an organization.
Unlike many alternatives, it was built with portability in mind making easy for anyone to get started.

For more information on how to get involved take a look at our [project board](https://github.com/orgs/depscloud/projects/1).

## hacktoberfest

Some tooling that facilitates insights into open source library usage using deps.cloud infrastructure.

Since deps.cloud knows how projects relate to one another, it also able to reason about library usage.
This includes both internal and open source libraries that support up your company.

`identify-contribution-candidates` is a purpose built command line tool for identifying high-impact candidates for open source contributions.
It works by using a [deployment](https:/deps.cloud/docs/deploy/) of deps.cloud to query for open source library usage across your company.
We score each library by the number of edges in the dependent sub-tree that point to a company's module.
When a library returns a non-0 score, we look up it's source location using [LibrariesIO](https://libraries.io).
Finally, the sorted results and written to a `candidate.json` file.

This list of repositories file can then feed into [Indeed's Mariner](https://github.com/indeedeng/Mariner) project.
Mariner takes this list of repositories and interrogates them for issues marked for contributions.
These issues help inform new contributors to open source how their contributions impact the company.

### Getting Started

1. Download the latest command from the [releases](https://github.com/depscloud/hacktoberfest/releases) tab.

1. Create a `config.yaml` file:
    ```yaml
    company_patterns:
      - ^.*depscloud.*$
    ```

2. Configure [deps.cloud](https://deps.cloud/docs/deploy/) endpoint.
    ```bash
    export DEPSCLOUD_BASE_URL=http://localhost:8080
    ``` 

3. Obtain an API Key from [libraries.io](https://libraries.io/)
    ```bash
    export LIBRARIESIO_API_KEY=123wxyz
    ```

4. Run `identify-contribution-candidates`
