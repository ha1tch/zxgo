# zxgotools

A collection of cross-platform ZX Spectrum tools written in Go.

## Overview

Many amazing ZX Spectrum tools are available, but sometimes you lack the right compiler, the right OS, or access to source code. These tools are designed to be easily compiled for most architectures and operating systems supported by Go.

## Tools

### LoadTAP

Reads and analyzes ZX Spectrum TAP files. Available for Windows (x64/i386), Linux (x64/i386), and macOS (ARM64).

```bash
loadtap [-d] [-r] <tap-file>
```

Options:
- `-d`: Dump block data as hex
- `-r`: Output raw block data (skips headers)

### MakeTAP

Creates TAP files from binary data. Available for Windows (x64/i386), Linux (x64/i386), and macOS (ARM64).

```bash
maketap [--name NAME] [--address ADDR] input.bin output.tap
```

Options:
- `--name`: Name for code block (max 10 chars)
- `--address`: Start address (default: 32768/0x8000)

### TAP2TZX

Converts TAP files to TZX format with additional metadata and features. Available for Windows (x64/i386), Linux (x64/i386), and macOS (ARM64).

```bash
tap2tzx [-o output.tzx] [-c config.yaml] [options] input1.tap [input2.tap ...]
```

Options:
- `-o`: Output TZX file (required)
- `-c`: YAML configuration file
- `-p`: Pause duration between blocks in ms (default: 1000)
- `-m`: Add metadata block
- `-title`: Program title (requires -m)
- `-author`: Program author (requires -m)
- `-year`: Publication year (requires -m)
- `-128`: Program requires 128K
- `-ay`: Program uses AY sound chip
- `-paging`: Program uses memory paging
- `-model`: Required model: +2, +2A, or +3
- `-multiload`: Program is multiload (adds 48K stop blocks)
- `-group`: Group name for following files

## Building

The project uses Go modules and supports cross-compilation for multiple platforms.

### Available Binaries

After building, you'll find these binaries in the `bin` directory:

Windows:
- `tool.exe` (64-bit/amd64)
- `tool.win32.exe` (32-bit/i386)

Linux:
- `tool.linux` (64-bit/amd64)
- `tool.linux32` (32-bit/i386)

macOS:
- `tool.mac` (ARM64)

Where `tool` is one of: `loadtap`, `maketap`, or `tap2tzx`

### Quick Build

Build all tools for all supported platforms:

```bash
./mk.sh
```

### Manual Building

You can also build individual tools for specific platforms:

```bash
# For Windows 64-bit
GOOS=windows GOARCH=amd64 go build -o bin/tool.exe ./cmd/tool

# For Windows 32-bit
GOOS=windows GOARCH=386 go build -o bin/tool.win32.exe ./cmd/tool

# For Linux 64-bit
GOOS=linux GOARCH=amd64 go build -o bin/tool.linux ./cmd/tool

# For Linux 32-bit
GOOS=linux GOARCH=386 go build -o bin/tool.linux32 ./cmd/tool

# For macOS ARM64
GOOS=darwin GOARCH=arm64 go build -o bin/tool.mac ./cmd/tool
```

Replace `tool` with `loadtap`, `maketap`, or `tap2tzx` as needed.

## Project Structure

```
zxgotools/
├── cmd/
│   ├── loadtap/
│   ├── maketap/
│   └── tap2tzx/
├── pkg/
│   └── tap/
├── bin/
├── LICENSE
├── README.md
└── mk.sh
```

## Contact

- Email: haitch@duck.com
- Mastodon: [@haitchfive@oldbytes.space](https://oldbytes.space/haitchfive)

## License

Licensed under the Apache License, Version 2.0. See LICENSE file for details.