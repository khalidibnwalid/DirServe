# DirServe

A simple, lightweight, HTTP server for exposing local directories over a local network. DirServe provides a simple way to serve files from any directory on your system with optional security features.

Note: *currently videos are not supported.*

![Screenshot][screenshot]


## Features

- Serve files from any local directory
- Basic authentication support

## Installation

### Via `go install` (Recommended)

```bash
go install github.com/khalidibnwalid/DirServe@latest
```

### Via Manual Build

```bash
# Clone the repository
git clone https://github.com/khalidibnwalid/DirServe

# Build the application
cd dirserve
go build
```

## Usage

```bash
dirserve [flags]
```

### Command-line Flags

| Flag      | Default | Description                             |
|-----------|--------:|-----------------------------------------|
| `-port`   | 8080    | Port to serve on                        |
| `-dir`    | .       | Directory to serve files from           |
| `-auth`   | false   | Enable basic authentication             |
| `-user`   | ""      | Username for basic authentication       |
| `-pass`   | ""      | Password for basic authentication       |

## Examples

### Serve the current directory on the default port (8080)

```bash
./dirserve
```

### Serve a specific directory on port 3000

```bash
./dirserve -port 3000 -dir /path/to/your/files
```

### Serve with basic authentication

```bash
./dirserve -auth -user admin -pass secret
```

## Endpoints

- `/raw/` - Direct file server access
- `/raw/` - web ui for browsing
- `/ping` - Health check endpoint (returns "pong")

## Security Considerations

When using DirServe, be aware that you are exposing files to the network. Use the authentication option when serving sensitive content, and be careful about which directories you choose to expose.

[screenshot]: https://raw.githubusercontent.com/khalidibnwalid/DirServe/refs/heads/main/assets/cover.webp