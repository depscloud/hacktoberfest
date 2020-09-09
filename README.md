# deps.cloud Hacktoberfest tooling

![GitHub](https://img.shields.io/github/license/depscloud/hacktoberfest.svg)
![Google Analytics](https://www.google-analytics.com/collect?v=1&cid=555&t=pageview&ec=repo&ea=open&dp=hacktoberfest&dt=hacktoberfest&tid=UA-143087272-2)

[Hacktoberfest] is an event put on by [DigitalOcean] every year.
It encourages participation in the open source community which grow every year.
Some companies have a hard time encouraging open source contributions.
Many work hard to tie their contributions to open source back to their company. 

[Hacktoberfest]: https://hacktoberfest.digitalocean.com/
[DigitalOcean]: https://www.digitalocean.com/

[deps.cloud] is a tool built to help companies understand how projects relate to one another.
It does this by detecting dependencies defined in common manifest files.
Using this information, we’re able to construct a dependency graph.
As a result, it's able to reason about library usage at a company.
This includes both internal and open source libraries that support up your company.

[deps.cloud]: https://deps.cloud

## Identifying contribution candidates

`identify-contribution-candidates` identifies candidates for open source contributions.
It works by querying a [deployment](https:/deps.cloud/docs/deploy/) of deps.cloud.
Using the information in the graph, we score each library by the number of edges in the dependent sub-tree.
These edges must point to a company owned module to be considered as part of the score.

When a library returns a non-0 score,we look up it's source information using [LibrariesIO](https://libraries.io).
Finally, we write the sorted results to a `candidate.json` file.

**Getting Started**

1. Download the latest command from the [releases](https://github.com/depscloud/hacktoberfest/releases) tab.

1. Create a `config.yaml` file.
   We use this to identify your companies modules in the graph.
    ```yaml
    company_patterns:
      - ^.*depscloud.*$
    ```

2. Configure your [deps.cloud](https://deps.cloud/docs/deploy/) endpoint.
    ```bash
    export DEPSCLOUD_BASE_URL=http://localhost:8080
    ``` 

3. Obtain an API Key from libraries.io
    ```bash
    export LIBRARIESIO_API_KEY=123wxyz
    ```

4. Run `identify-contribution-candidates`

Once complete, you'll be left with a file that contains a JSON array.
Each entry will have two fields: a source url and a score.
The source URL is the location of the source code.
The score is the sum of all scores for libraries produced by that source URL.

```json
[
  {
    "repository_url": "https://github.com/mjpitz/go-gracefully",
    "score": 4
  }
]
```

## Feeding into Mariner

Indeed's [Mariner] project takes in a list of GitHub repositories and their scores.
Because Mariner only works with GitHub repositories, we will need to filter and format accordingly.
After feeding this information through Mariner, you will be left with a set of tickets.
These tickets are not only great for getting started in open source, but also have a direct impact on your company. 

[Mariner]: https://github.com/indeedeng/Mariner

# Support

Join our [mailing list] to get access to virtual events and ask any questions there.

We also have a [Slack] channel.

[mailing list]: https://groups.google.com/a/deps.cloud/forum/#!forum/community/join
[Slack]: https://depscloud.slack.com/join/shared_invite/zt-fd03dm8x-L5Vxh07smWr_vlK9Qg9q5A

## Branch Checks

![branch](https://github.com/depscloud/hacktoberfest/workflows/branch/badge.svg?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/depscloud/hacktoberfest)](https://goreportcard.com/report/github.com/depscloud/hacktoberfest)
