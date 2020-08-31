# Haste-Client

This is a rewritten Haste-Client in Go which is meant to be a little utility that uploads code via command line from pipe.
Available for both Windows and Linux.


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

Changes upload destination, can be another haste server, requires schema prefix (default: https://haste.zneix.eu).  

# Easy Installation

Download binary the latest release from [releases page](https://github.com/zneix/haste-client/releases/tag/1.0).  
Version without extension if for Linux, and `.exe` is for Windows.  

To make sure you can use this tool globally in command line, make sure to put it into your PATH Environment Variable:  
[Guide for windows](https://helpdeskgeek.com/windows-10/add-windows-path-environment-variable/)  
[Guide for Linux](https://helpdeskgeek.com/windows-10/add-windows-path-environment-variable/)


# Installing on Windows

Download [Go 1.15](https://golang.org/doc/install?download=go1.15.windows-amd64.msi)

After installing go, open Command Line with `Win + R` and typing `cmd`.  
Execute these commands:
```bash
go get github.com/zneix/haste-client
cd %userprofile%\go\bin
ren haste-client.exe haste.exe
```


# Installing on Linux

Download [Go 1.15](https://golang.org/doc/install?download=go1.15.linux-amd64.tar.gz)

```bash
go get github.com/zneix/github
make build
sudo cp ~/go/bin/haste-client /usr/local/bin/haste
```

Or just clone GitHub repository and build with `make`.