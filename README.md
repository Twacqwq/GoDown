# A Golang-based command line concurrent download tool

## Installation
```bash
$ go install github.com/Twacqwq/godown@latest
```

## Usage
```bash
godown -u https://example.com/example.zip
```
```bash
godown -u https://example.com/example.zip -o /tmp/
```

## help
```bash
$ godown --help
```

## TODO
- 支持断点续传
- 支持磁力解析