# htproxctl

`htproxctl` sets the `HTTP_PROXY` and `HTTPS_PROXY` environment variables based on your configuration before executing the specified command. This allows the command to use the configured proxy without modifying its own behavior.

## Features

- Configure multiple proxy environments
- Automatically set HTTP_PROXY and HTTPS_PROXY environment variables

## Installation

Build for tool:

```
$ go build
```

## Configuration

Create a configuration file at `~/.config/htproxctl.yaml` with the following structure:

```yaml
defaults:
  proxy: socks5://localhost
  port: 11090

environments:
  dev:
    proxy: http://dev-proxy.example.com
    port: 11095
  stg:
    port: 11096
  prod:
    proxy: https://prod-proxy.example.com
    port: 11097
```

## Usage

The basic syntax for using `htproxctl` is:

```
$ htproxctl [-env <environment>] [--] <command> [args...]
```

### Examples

Run a command using the default proxy settings:

```
$ htproxctl -- kubectl get pods
```

Run a command using a specific environment's proxy settings:

```
$ htproxctl -env dev -- kubectl get pods
```

## License

This project is licensed under the [MIT License](./LICENSE).
