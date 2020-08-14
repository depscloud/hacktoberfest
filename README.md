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

Since deps.cloud knows how projects relate to one another, it also know about library usage.
This includes both internal libraries used across your company and open source libraries your company was built on.
This project currently includes two commands:

`identify-oss-libraries` uses a deps.cloud [deployment](https:/deps.cloud/docs/deploy/) to identify library use across your company.
This script will output a file containing a sorted list of libraries and a score (the number of edges in their dependent subtree.)

`translate-oss-libraries` takes the file output from the previous command and looks each library's source up using [LibrariesIO](https://libraries.io).
In the end, this outputs a file containing a sorted list of repositories and a score (the sum of all mapped library scores.)

This list of repositories file can then feed into [Indeed's Mariner](https://github.com/indeedeng/Mariner) project.
Mariner takes this list of repositories and interrogates them for issues marked for contributions.
This issues are then used to help new contributors to open source how their contributions impact the company.
