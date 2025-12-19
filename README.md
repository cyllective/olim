<h1 align="center">onetim3</h1>

<p align="center">
  <i>A simple, encrypted, one-time secret sharing app</i>
</p>

<p align="center">
  <a href="https://github.com/cyllective/onetim3/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/cyllective/onetim3" alt="LICENSE">
  </a>
  <a href="https://github.com/cyllective/onetim3/actions">
    <img src="https://img.shields.io/github/actions/workflow/status/cyllective/onetim3/ghcr.yaml" alt="Build Status" />
  </a>
</p>

## Setup

This is intended to be used as a Docker container. You can either use the image `ghcr.io/cyllective/onetim3` or build the image yourself with `make docker-image`. As an example on how to deploy it, see [docker-compose.yaml](./docker-compose.yaml).

> [!IMPORTANT]
> **HTTPS is required!** The Web Crypto API will only work on localhost or HTTPS connections.

### Options

The app is configured with environment variables.

| Variable  | Affects                                  | Default            |
| --------- | ---------------------------------------- | ------------------ |
| `DEBUG`   | If debug output should be printed        | -                  |
| `DB_PATH` | Where the SQLite3 database will be saved | `./onetim3.sqlite` |

## Development

For development, use [air](https://github.com/air-verse/air). There is a configuration file already present. It opens a proxy on port 8081.