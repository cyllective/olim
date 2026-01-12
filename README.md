# olim

_A simple, encrypted, one-time secret sharing app_

<p>
  <a href="https://github.com/cyllective/olim/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/cyllective/olim" alt="LICENSE">
  </a>
  <a href="https://github.com/cyllective/olim/actions">
    <img src="https://img.shields.io/github/actions/workflow/status/cyllective/olim/ghcr.yaml" alt="Build Status" />
  </a>
</p>

## Setup

This is intended to be used as a Docker container. You can either use the image `ghcr.io/cyllective/olim` or build the image yourself with `make docker-image`. As an example on how to deploy it, see [docker-compose.yaml](./docker-compose.yaml).

> [!IMPORTANT]
> **HTTPS is required!** The Web Crypto API will only work on localhost or HTTPS connections.

### Options

The app is configured with environment variables.

| Variable  | Affects                                  | Default            |
| --------- | ---------------------------------------- | ------------------ |
| `DEBUG`   | If debug output should be printed        | -                  |
| `DB_PATH` | Where the SQLite3 database will be saved | `./olim.sqlite` |

## Development

For development, use [air](https://github.com/air-verse/air). There is a configuration file already present. It opens a proxy on port 8081.

## How does this work?

[The Web Crypto API](https://developer.mozilla.org/en-US/docs/Web/API/Web_Crypto_API) is used to generate AES-GCM keys. Those encrypt the text (or the contents of the file). This encrypted blob is then sent to the server. The keys are added as a [URL fragment (#)](https://developer.mozilla.org/en-US/docs/Web/URI/Reference/Fragment) to the share URL. Since those are not sent to the HTTP server when opening them in a browser, the server never gets the keys necessary for decryption. 

## Name

The word "olim" is Latin for "once upon a time".

## Kudos

This was inspired by [transfer.pw](https://transfer.pw/).