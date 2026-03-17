# `csvv`

A CLI tool to inspect CSV structured data.

## Installation

### GitHub Releases

Download binary from [GitHub Releases](https://github.com/kostya-zero/csvv/releases) for you system.

### With Go toolchain

If you have Go installed you can run this command to install it:

```shell
go install github.com/kostya-zero/csvv
```

## Usage

```shell
# You can just give path to the CSV file and it will print it's content.
csvv data.csv

# You can also do additional operations.
csvv data.csv --first 7
csvv data.csv --last 3
csvv data.csv --select 103
```
