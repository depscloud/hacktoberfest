# deps.cloud Hacktoberfest tooling 2021

This repository contains tooling to help support organizations participating in [Hacktoberfest].

[Hacktoberfest]: https://hacktoberfest.digitalocean.com/

![GitHub](https://img.shields.io/github/license/depscloud/hacktoberfest.svg)
![Google Analytics](https://www.google-analytics.com/collect?v=1&cid=555&t=pageview&ec=repo&ea=open&dp=hacktoberfest&dt=hacktoberfest&tid=UA-143087272-2)

deps.cloud is a tool to help understand how projects relate to one another.
It works by detecting dependencies defined in common [manifest files] (`pom.xml`, `package.json`, `go.mod`, etc).
Using this information, we’re able to answer questions about project dependencies.

* What versions of _k8s.io/client-go_ do we depend on?
* Which projects use _eslint_ as a non-dev dependency?
* What open source libraries do we use the most?

[manifest files]: https://deps.cloud/docs/concepts/manifests/

## Prerequisites

In order to use this repository, you will need a [deployment] of deps.cloud running.

[deployment]: https:/deps.cloud/docs/deploy/

## Identifying contribution candidates

`identify-contribution-candidates` looks for open source library use across your company.
It scores each project by how often it's used in your projects.
We then look up its source using [LibrariesIO](https://libraries.io).

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

That's it!
At the end of it's run, it writes out a `candidate.json` file.

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
