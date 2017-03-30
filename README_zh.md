# Ego e
<!--<img align="right" src="https://raw.githubusercontent.com/go-ego/ego/master/logo.jpg">-->
<!--[![Build Status](https://travis-ci.org/go-ego/ego.svg)](https://travis-ci.org/go-ego/ego)
[![codecov](https://codecov.io/gh/go-ego/ego/branch/master/graph/badge.svg)](https://codecov.io/gh/go-ego/ego)-->
<!--<a href="https://circleci.com/gh/go-ego/ego/tree/dev"><img src="https://img.shields.io/circleci/project/go-ego/ego/dev.svg" alt="Build Status"></a>-->
[![CircleCI Status](https://circleci.com/gh/go-ego/e.svg?style=shield)](https://circleci.com/gh/go-ego/e)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-ego/e)](https://goreportcard.com/report/github.com/go-ego/e)
[![GoDoc](https://godoc.org/github.com/go-ego/e?status.svg)](https://godoc.org/github.com/go-ego/e)
[![Release](https://github-release-version.herokuapp.com/github/go-ego/e/release.svg?style=flat)](https://github.com/go-ego/e/releases/latest)
[![Join the chat at https://gitter.im/go-ego/ego](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/go-ego/ego?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
<!--<a href="https://github.com/go-ego/ego/releases"><img src="https://img.shields.io/badge/%20version%20-%206.0.0%20-blue.svg?style=flat-square" alt="Releases"></a>-->
  
  >E 是协助 ego 框架进行开发的命令行工具. Thanks [bee](https://github.com/beego/bee) for providing power.

这是一项正在完善的工作.

## Contents
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
- [Plans](#plans)
- [Donate](#donate)
- [Contributing](#contributing)
- [License](#license)

## Requirements:

- Go version >= 1.3

## Installation:
```
go get -u github.com/go-ego/e 
```
## Usage

```sh

USAGE
    ego command [arguments]

AVAILABLE COMMANDS

    new         auto-generate code for the ego app, Creates a ego API application
    run         Run the application by starting a local development server
    pack        Compresses a Ego application into a single file
    bale        Transforms non-Go files to Go source files
    version     Prints the current E version
    migrate     Runs database migrations
    fix         fix the ego application to make it compatible with ego 1.0

Use ego help [command] for more information about a command.

```

### e new 

To create a new Ego web application

### e run

To run the application we just created, you can navigate to the application folder and execute:
```
$ cd my-webapp && e run
```
Or from anywhere in your machine:
```
$ e run github.com/user/my-webapp
```
For more information on the usage, run e help run.

## Plans
- generate code and docs
- generating a dockerfile
- help with debugging your application

## Donate
- Supporting ego, [buy me a coffee](https://github.com/go-vgo/buy-me-a-coffee).
## Contributing

- To contribute to e, please see [Contribution Guidelines](https://github.com/go-ego/e/blob/master/CONTRIBUTING.md).

- See [contributors page](https://github.com/go-ego/e/graphs/contributors) for full list of contributors.

## License

E is primarily distributed under the terms of both the MIT license and the Apache License (Version 2.0).

See [LICENSE-APACHE](http://www.apache.org/licenses/LICENSE-2.0), [LICENSE-MIT](https://github.com/go-ego/ego/blob/master/LICENSE), and COPYRIGHT for details.