# gojs-tool

A tool for export package funcs, types, vars for VM (like goja) to use

### Install

```
go get github.com/gogap/gojs-tool
```

### Example

```bash
example> pwd
github.com/xujinzheng/playgo/example

example> gojs-tool gen -r github.com/sirupsen/logrus

example> tree modules
modules
└── github.com
    └── sirupsen
        └── logrus
            ├── hooks
            │   ├── syslog
            │   │   └── syslog.go
            │   └── test
            │       └── test.go
            └── logrus.go

```

`logrus.go`

```go
package logrus

import (
	"github.com/sirupsen/logrus"
)

import (
	"github.com/dop251/goja"
	"github.com/gogap/gojs-tool/gojs"
)

var (
	module = gojs.NewGojaModule("logrus")
)

func init() {
	module.Set(
		gojs.Objects{
			"AddHook":             logrus.AddHook,
			"Debug":               logrus.Debug,
			"Debugf":              logrus.Debugf,
			"Debugln":             logrus.Debugln,
			"Error":               logrus.Error,
			"Errorf":              logrus.Errorf,
			// ...
			// ... more code ...
			// ...
		},
	).Register()
}

func Enable(runtime *goja.Runtime) {
	module.Enable(runtime)
}
```


#### Use logurs in goja

```go
package main

import (
	"github.com/xujinzheng/playgo/example/modules/github.com/sirupsen/logrus"
)

import (
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

func main() {
	registry := new(require.Registry)

	runtime := goja.New()

	registry.Enable(runtime)
	logrus.Enable(runtime)

	runtime.RunString(`
    	logrus.WithField("Hello", "World").Println("I am gojs")
    `)
}
```


```bash
example> go run main.go
INFO[0000] I am gojs      Hello=World
```