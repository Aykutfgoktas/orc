# ORC - Organization Repo Cloner

[![test-img]][test-url]
[![lint-img]][lint-url]
[![version-img]][version-url]

![alt text](./orc.png)

## Introduction

This is a Command Line Interface (CLI) application written in Golang that allows users to clone repositories from a specified organization on GitHub. Additionally, the application can list the organizations that the user has entered and list the repositories of the selected organization.

## Installation

Before using this application, you need to have Golang installed on your system. You also need to have a GitHub account and generate an access token from your account settings.

```sh
go install github.com/Aykutfgoktas/orc
```

## Usage

Without using with any flag, it will list the repositories from the default organization.

|   Flag   | Flag Long |         Description          |
| :------: | :-------: | :--------------------------: |
|    -a    |   -add    |     add the organization     |
|    -l    |  --list   |    list the organization     |
|    -s    |   --set   | set the default organization |

Configuration will be stored in `$HOME/.orc.conf.json` file.

## Linting

- Install [golangci-lint](https://github.com/golangci/golangci-lint)

- You can run golang-lint local environment before committing the code.
- It will read the rules from the [.golangci.yml](https://github.com/Aykutfgoktas/orc/blob/master/.golangci.yml) file.

```sh
  make lint
```

## Testing

- [Ginkgo](https://onsi.github.io/ginkgo/) is used for the golang bdd testing framework.
- You can trigger commands on the makefile for easy execution of unit.

```sh
 make test
```

## Contributing

Contributions to this project are welcome. To contribute, please fork this repository, make your changes, and submit a pull request.

## Bug Reports and Support

If you encounter any bugs or issues with this application, please open an issue on the GitHub repository page.

## Backstory

When I started a new job, I found myself frequently needing to clone Git repositories to my local machine. However, this involved opening the repository's webpage, copying the clone URL, and then pasting it into my terminal. This process was time-consuming and error-prone, as I sometimes copied the wrong URL or forgot to include the ".git" extension.

To address this problem, I decided to create a command-line interface (CLI) tool that would automate the process of cloning Git repositories. My goal was to create a simple, easy-to-use tool that would save me time.

To be honest, I do not know if there are some other applications that solve the same problem. However, I just want to create my own to solve my problem and any feedback will be appreciated thank you.

In the future, I plan to continue refining the tool and adding new features, such as support for other version control systems and integration with code review tools. I also hope to share the tool with other developers and contribute to the open-source community.

[test-img]: https://github.com/Aykutfgoktas/orc/workflows/go-test/badge.svg
[test-url]: https://github.com/Aykutfgoktas/orc/workflows/go-test/badge.svg
[lint-img]: https://github.com/Aykutfgoktas/orc/workflows/golangci-lint/badge.svg
[lint-url]: https://github.com/Aykutfgoktas/orc/workflows/golangci-lint/badge.svg
[version-img]: https://img.shields.io/github/v/release/Aykutfgoktas/orc
[version-url]: https://github.com/Aykutfgoktas/orc/releases