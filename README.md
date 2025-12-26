# File Crow

A peer-to-peer file server for sharing files over a local network.

![screenshot](./screenshot.webp)

## Motivation

I often transfer files between my PC and mobile device using a USB cable. But every time I want to share something I have to take the cable out of the charger and connect both devices. So I built File Crow to share files over the local network.

## Quick Start

\* [Go](https://go.dev/doc/install) v1.25.5 or higher and [Node](https://nodejs.org/en/download) v24.11.1 or higher are required.

1. Install filecrow

```bash
# clone the repo
git clone https://github.com/syeero7/filecrow
cd filecrow

# install dependencies
go mod tidy
npm install

# build the compiled binary
npm run build
```

2. Start the file server

```bash
filecrow
```

3. Access the web interface at `http://localhost:<PORT>` on the first device. The default port is 8080.

4. On other devices, access via `http://<DEVICE_IP>:<PORT>` replacing `<DEVICE_IP>` with the local IP address.

5. Upload files to share from any device and download them on other devices.

## Usage

Available flags:

- `--port` - The port number to run the server on (default 8080)
- `-h`, `--help` - Show help

### Examples

```bash
# Start with a custom port
filecrow --port 8000
```

## Contributing

### Clone the repo

```bash
git clone https://github.com/syeero7/filecrow
cd filecrow
```

### Build the compiled binary

```bash
npm run build
```

### Submit a pull request

If you'd like to contribute, please fork the repository and open a pull request to the `main` branch.
