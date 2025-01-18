# zxgio
#### Portable ZX Spectrum tools written in Go

Many amazing tools are available out there, and have been available for decades now, but sometimes you lack the right compiler, or the right OS. Sometimes you don't have access to the full source code, sometimes you've just turned off your last Windows machine 10 years ago and don't want to use a Windows emulator.

Say no more, these are trivial to compile for most architectures and operating system pairs, provided the pair is supported by Go.

```bash
GOOS=windows GOARCH=amd64 go build -o makketap.exe  maketap.go
GOOS=linux   GOARCH=amd64 go build -o maketap.linux maketap.go
GOOS=darwin  GOARCH=arm64 go build -o maketop.mac   maketap.go
```

```bash
go build loadtap.go
```

 
