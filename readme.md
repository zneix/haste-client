# Haste-Client

This is a rewritten Haste-Client in Go which is meant to be a little utility that uploads code via command line from pipe.


# Usage

Examples:

`echo Sample Text | haste`  
> https://haste.zneix.eu/nibazahidu

`cat veryLongScript.js | haste | xsel`
> *copies https://haste.zneix.eu/ibadomuvaq to clipboard*

<br>

## Arguments

`-h`

Shows Help and exits

`-v`

Shows Program version and exits

`-d string`

Changes upload destination, can be another haste server, requires schema prefix (default: https://haste.zneix.eu)  


# Installation

Requirements: [go 1.15](https://golang.org/dl/)

Build with Makefile:

```bash
make build
sudo cp ./haste /usr/local/bin/
```