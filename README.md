# slog

A simple log lib, support log level.

Features:

* Package categorized logger
* Support log level
* Configured via json

## Requires

Only go standard libs

## Usage

### Get a logger for your package

Assume there are there packages: github.com/kuun/slog/demo,
github.com/kuun/slog/demo/pkga, github.com/kuun/slog/demo/pkgb, structured as bellow:

```
src/github.com/kuun/slog/demo
|-- main.go
|-- pkga
|   `-- a.go
`-- pkgb
    `-- b.go
```

src/github.com/kuun/slog/demo/main.go

```go
package main

import (
	"test/pkga"
	"test/pkgb"

	"github.com/kuun/slog"
)

var log = slog.GetLogger()

func main() {
	log.Debug("hello slog")
	pkga.Hello()
	pkgb.Hello()
	slog.Close()
}
```

src/github.com/kuun/slog/demo/pkga/a.go

```go
package pkga

import "github.com/kuun/slog"

var log = slog.GetLogger()

func Hello() {
	log.Debug("hello pkga")
}
```

src/github.com/kuun/slog/demo/pkgb/b.go

```go
package pkgb

import "github.com/kuun/slog"

var log = slog.GetLogger()

func Hello() {
	log.Debug("hello pkgb")
}
```

then build and run the app, outputs on the terminal:

```
D 0531 11:13:15.956194 test/main.go:13] hello slog
D 0531 11:13:15.956257 t/pkga/a.go:8] hello pkga
D 0531 11:13:15.956259 t/pkgb/b.go:8] hello pkgb
```

Get a new logger by calling slog.GetLogger, slog will create a slog.Logger categorized
by current package, if called more than once in the same package, slog will return
the same slog.Logger.

### Configure slog

Slog configure file a json object. If there is not configure file or configure
file is empty, slog uses the default configure:

```javascript
{
    "loggers": [
        {
            "pattern": "*",
            "level": "DEBUG",
            "writers": ["STDOUT"]
        }
    ]
}
```

Now we want to write log to a file, we can modify the configuration like bellow:

```javascript
{
    "writers": [
        {
            "name": "FILE",
            "file": "/tmp/app.log"
        }
    ],
    "loggers": [
        {
            "pattern": "*",
            "level": "DEBUG",
            "writers": ["FILE", "STDOUT"]
        }
    ]
}
```
Then save configuration to slog.json.

Then modify environment SLOG_CONF_FILE and run the application:

```shell
export SLOG_CONF_FILE=slog.json
./test
```

Logs will be outputed to stdout and /tmp/app.log

## Configure slog

### How to specific log configuration file

Slog parse environment value 'SLOG_CONF_FILE' at application startups.

### Configuration spec

* writers

  writers is an array, collects all available writers. There are two predefined
  writer: "STDOUT",  "STDERR".

  * writers.name

    Writer's name, must be unique.

  * writers.file

    The file where writer writes to

* loggers

  loggers is an array, collects all logger configuration

  * loggers.pattern

    The current configuration's match pattern, the first and last character can
    be '\*' to wildcard multi packages.

  * loggers.level

    log level of loggers matched by the pattern.

  * loggers.writers

    writers of loggers matched by the pattern.
