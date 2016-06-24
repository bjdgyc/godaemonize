#godaemonize [![Build Status](https://travis-ci.org/bjdgyc/godaemonize.svg?branch=master)](https://travis-ci.org/bjdgyc/godaemonize)

## Introduction

godaemonize is a command-line utility that runs a command as a Unix daemon.

The godaemonize is most like [daemonize](https://github.com/bmc/daemonize)

## Installation

`go get github.com/bjdgyc/godaemonize`

The bin file will install in `GOPATH` src/bin dir.

## Usage

```
godaemonize, version 0.1.1
Usage: godaemonize [OPTIONS] -x file [ARGV] ...

OPTIONS

  -E string
    	Pass environment setting to daemon. like [a=b,c=d]
  -d string
    	Set daemon's working directory to <dir>
  -e string
    	Send daemon's stderr to file, default is <stderr>
  -o string
    	Send daemon's stdout to file, default is <stdout>
  -p string
    	Save PID to <pidfile>
  -u string
    	Run daemon as user <user>. Requires invocation as root
```


## Quickstart


```
godaemonize -p /tmp/some.pid -u nobody -E "key1=value1,key2=value2" -e /tmp/some_err.log -o /tmp/some_out.log -x /tmp/some.sh some parameter
```




