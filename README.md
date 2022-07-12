<p align="center">
  <img alt="GoReleaser Logo" src="https://www.blinkops.com/favicon.ico" height="140" />
  <h3 align="center">Blink CLI</h3>
</p>

---
The Blink CLI helps you build and manage your Blink operations right from the terminal.

**With the CLI, you can:**

- Create, retrieve, update, or delete API objects.
- Run automations.

## Installation

Blink CLI is available for macOS, Windows, and Linux for distros like Ubuntu, Debian, RedHat and CentOS.

### macOS

Blink CLI is available on macOS via [Homebrew](https://brew.sh/):

```sh
brew tap blinkops/blink-go-cli https://github.com/blinkops/blink-go-cli
brew install blink
```

### Linux (Under Development)
WIP

### Windows (Under Development)

Blink CLI is available on Windows via the [Scoop](https://scoop.sh/) package manager:
WIP 

### Docker (Under Development)

The CLI is also available as a Docker image: [`blinkops/blink-cli`](https://hub.docker.com/r/blinkops/blink-cli).

```sh
docker run --rm -it blinkops/blink-cli version
```

### Without package managers (Under Development)

WIP

## Init
To first initialize your credentials after installing run:
```sh-session
blink init
```
You will be prompted to add a hostname (will be https://app.blinkops.com by default unless Blink is run locally within the org)
and an API-key
```sh-session
✔ Hostname: https://app.dev.blinkops.com
✔ Blink API Key (Obtain key by accessing https://app.dev.blinkops.com/api/v1/apikey in your browser)
```
If you wish to use a pre-defined configuration, you can specify a file path by setting the ``--config`` flag when running blink init
Example:
```sh-session
blink init --config ./config.json
```

At a minimum you should set the following values in your config file.

```json
{
  "hostname": "<blink-address>", 
  "blink-api-key": "<apikey>",
  "scheme": "https"
}
```

## Usage

```sh-session
blink [command]

# Run `--help` for detailed information about CLI commands
blink [command] --help
```


### Workspaces

By default, blink uses the workspace specified in your config file.
If a workspace is not set in your config file, or you would like to use a different workspace
than the one set, use the ``--workspace`` flag.

```sh-session
blink [command] --workspace my_workspace
```


## Commands

The Blink CLI supports a broad range of commands. Below is some of the most used ones:
- `automations`
- `runner-groups`
- `connections`

## Feedback

Got feedback for us? Please don't hesitate to shoot us a message.

