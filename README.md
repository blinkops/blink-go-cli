<p align="center">
  <img alt="GoReleaser Logo" src="https://www.blinkops.com/favicon.ico" height="140" />
  <h3 align="center">Blink CLI</h3>
</p>

---
The Blink CLI helps you build and manage your Blink operations right from the terminal.

**With the CLI, you can:**

- Create, retrieve, update, or delete API objects.
- Run playbooks

## Installation

Blink CLI is available for macOS, Windows, and Linux for distros like Ubuntu, Debian, RedHat and CentOS.

### macOS

Blink CLI is available on macOS via [Homebrew](https://brew.sh/):

```sh
brew tap blinkops/blink-go-cli https://github.com/blinkops/blink-go-cli
brew install blink-cli
```

### Linux (Under Development)
WIP

### Windows (Under Development)

Blink CLI is available on Windows via the [Scoop](https://scoop.sh/) package manager:
WIP 

### Docker(Under Development)

The CLI is also available as a Docker image: [`blinkops/blink-cli`](https://hub.docker.com/r/blinkops/blink-cli).

```sh
docker run --rm -it blinkops/blink-cli version
```

### Without package managers (Under Development)

WIP


## Usage

```sh-session
blink-cli [command]

# Run `--help` for detailed information about CLI commands
blink-cli [command] --help
```

### Configuration

By default, blink-cli looks for a file named ``config.json`` in the ``$HOME/.config/blink-cli/config.json`` directory.
You can specify a different file path by setting the ``--config`` flag

```sh-session
blink-cli [command] --config ./config.json
```

At a minimum you should set the following values in your config file.

```json
{
  "hostname": "<blink-address>", 
  "BLINK-API-KEY": "<apikey>",
  "scheme": "https"
}
```

## Commands

The Blink CLI supports a broad range of commands. Below is some of the most used ones:
- `playbooks`
- `runners`

## Feedback

Got feedback for us? Please don't hesitate to shoot us a message.

