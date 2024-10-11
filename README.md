# Aurora
[![GitHub License](https://img.shields.io/github/license/mr5356/aurora)](https://github.com/mr5356/aurora?tab=Apache-2.0-1-ov-file)

[![go-test](https://github.com/MR5356/aurora/workflows/Go%20Test/badge.svg?query=branch%3Amaster)](https://github.com/MR5356/aurora/actions?query=branch%3Amaster)
[![docker-image](https://github.com/MR5356/aurora/workflows/Docker%20Build/badge.svg?query=branch%3Amaster)](https://hub.docker.com/r/toodo/aurora/tags)
[![go-report](https://goreportcard.com/badge/github.com/MR5356/aurora)](https://goreportcard.com/report/github.com/MR5356/aurora)

[![Docker Pulls](https://img.shields.io/docker/pulls/toodo/aurora)](https://hub.docker.com/r/toodo/aurora)
[![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/mr5356/aurora/total)](https://github.com/MR5356/aurora/releases)
[![GitHub Release Date](https://img.shields.io/github/release-date/mr5356/aurora)](https://github.com/MR5356/aurora/releases)
[![release](https://img.shields.io/github/v/release/MR5356/aurora)](https://github.com/MR5356/aurora/releases)

<img src="./logo/logo.svg" width="100" alt="logo">

----

Aurora is an open source system for DevOps, consisting of:
* convenient **Admin dashboard UI**
* **Host** and **scheduled task** management
* **Health check** etc.
* and simple **REST-ish API**

**For documentation and examples, please visit [https://aurora.docker.ac.cn](https://aurora.docker.ac.cn).**

----

## To start developing Aurora

The [repository](/) hosts all information about 
building Aurora from source, how to contribute code 
and documentation, who to contact about what, etc.

If you want to build Aurora right away there are two options:

##### You have a working Go environment

```shell
git clone https://github.com/MR5356/aurora.git
cd aurora
make init build
```

##### You have a working Docker environment

```shell
git clone https://github.com/MR5356/aurora.git
cd aurora
make docker
```

## Support

If you have questions, reach out to us one way or another.

## Acknowledgements
Grateful acknowledgement to [JetBrains](https://www.jetbrains.com) for supporting this project through their Open Source License Program and providing exceptional development tools.