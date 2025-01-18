# TAP2TZX

A tool to convert ZX Spectrum TAP files to TZX format with support for advanced features like metadata, hardware requirements, grouping, and flow control.

## Usage

```bash
tap2tzx [-o output.tzx] [-c config.yaml] [options] input1.tap [input2.tap ...]
```

### Basic Example

Convert multiple TAP files with groups:

```bash
tap2tzx -o simple.tzx \
  --group "The Three Lights of Glaurung" loader1.tap data1.tap \
  --group "Level 2" loader2.tap data2.tap
```

### Advanced Example with Configuration File

Create a complex TZX with metadata, hardware requirements, and advanced features using a YAML configuration:

```yaml
metadata:
  title: "This is a very complex project"
  author: "haitch"
  year: "2025"

hardware:
  128k_only: true
  use_ay: true
  model: "+2"

blocks:
  # Menu section
  - group: "Menu"
    file: menu.tap
    desc: "A Fantastic Menu"

  # Level sections with various features
  - group: "Breathtaking Level 1"
    file: loader1.tap
    desc: "Loader of the amazing Level 1"
  - file: screen1.tap
  - file: data1.tap

  # Level with jump option
  - id: level2_start
    group: "The Astonishing Level 2"
    file: loader2.tap
  - file: screen2.tap
  - file: data2.tap
  - jump_to: end  # Skip level 3 if needed

  # Sound data with loop
  - group: "Sound Data"
    desc: "Background Music"
  - loop_start: 2
  - file: sound1.tap
  - file: sound2.tap
  - id: loop_end
```

Run with configuration:

```bash
tap2tzx -128 -author haitch -ay -m -title "Fabulous Game" -year 2025 \
  -c complex_tzx.yaml -o complex.tzx
```

## Command Line Options

### Output Options
- `-o`: Output TZX file (required)
- `-c`: YAML configuration file

### Metadata Options
- `-m`: Add metadata block
- `-title`: Program title
- `-author`: Program author
- `-year`: Publication year

### Hardware Options
- `-128`: Program requires 128K
- `-ay`: Program uses AY sound chip
- `-paging`: Program uses memory paging
- `-model`: Required model (+2, +2A, or +3)

### Loading Options
- `-p`: Pause duration between blocks in ms (default: 1000)
- `-multiload`: Program is multiload (adds 48K stop blocks)
- `-group`: Group name for following files

## YAML Configuration Features

### Metadata Section
```yaml
metadata:
  title: "Game Title"
  author: "Author Name"
  year: "2025"
```

### Hardware Requirements
```yaml
hardware:
  128k_only: true
  use_ay: true
  model: "+2"  # +2, +2A, or +3
```

### Block Features

Groups:
```yaml
- group: "Level 1"
  file: level1.tap
```

Descriptions:
```yaml
- desc: "Loading screen"
  file: screen.tap
```

Jump Points:
```yaml
- id: skip_point
  file: data.tap
- jump_to: skip_point
```

Loops:
```yaml
- loop_start: 2
- file: music.tap
- loop_end: true
```

## Building

```bash
# Build for all supported platforms
./mk.sh

# Or build manually for specific platform
GOOS=linux GOARCH=amd64 go build -o ../../bin/tap2tzx.linux
```

## License

Licensed under the Apache License, Version 2.0. See LICENSE file for details.

## Contact

- Email: haitch@duck.com
- Mastodon: [@haitchfive@oldbytes.space](https://oldbytes.space/@haitchfive)