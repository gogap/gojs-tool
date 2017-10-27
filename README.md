# gojs-tool

A tool for export package funcs, types, vars for VM (like goja) to use

### Install

```
go get -u -v github.com/gogap/gojs-tool
```


### Command

```bash
> gojs-tool gen --help
 
NAME:
   gojs-tool gen - generate modules

USAGE:
   gojs-tool gen [command options] packages

OPTIONS:
   --template value, -t value  template filename in templates folder without extension (default: "goja")
   --recusive, -r              recusive generate code
   --namespace, -n             generate submodules loader, access object by dot, only work while recusive enabled
   --gopath value              the package in which GOPATH [$GOPATH]
```


### Example

```bash
example> pwd
github.com/xujinzheng/playgo/example

example> gojs-tool gen -r -n github.com/sirupsen/logrus

example> tree modules
modules
└── github.com
    └── sirupsen
        └── logrus
            ├── hooks
            │   ├── hooks.go
            │   ├── syslog
            │   │   └── syslog.go
            │   └── test
            │       └── test.go
            └── logrus.go

```

#### Auto generated code

`modules/github.com/sirupsen/logrus/logrus.go`

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
			// Functions
			"AddHook":             logrus.AddHook,
			"Debug":               logrus.Debug,
			"Debugf":              logrus.Debugf,
			"Debugln":             logrus.Debugln,
			"Error":               logrus.Error,
			"Errorf":              logrus.Errorf,
			"Errorln":             logrus.Errorln,
			"Exit":                logrus.Exit,
			"Fatal":               logrus.Fatal,
			"Fatalf":              logrus.Fatalf,
			//..........more funcs..........		


			// Var and consts
			"AllLevels":     logrus.AllLevels,
			"DebugLevel":    logrus.DebugLevel,
			"ErrorKey":      logrus.ErrorKey,
			"ErrorLevel":    logrus.ErrorLevel,
			"FatalLevel":    logrus.FatalLevel,
			"FieldKeyLevel": logrus.FieldKeyLevel,
			"FieldKeyMsg":   logrus.FieldKeyMsg,
			"FieldKeyTime":  logrus.FieldKeyTime,
			"InfoLevel":     logrus.InfoLevel,
			"PanicLevel":    logrus.PanicLevel,
			"WarnLevel":     logrus.WarnLevel,

			// Types (value type)
			"Entry":         func() logrus.Entry { return logrus.Entry{} },
			"JSONFormatter": func() logrus.JSONFormatter { return logrus.JSONFormatter{} },
			"Logger":        func() logrus.Logger { return logrus.Logger{} },
			"MutexWrap":     func() logrus.MutexWrap { return logrus.MutexWrap{} },
			"TextFormatter": func() logrus.TextFormatter { return logrus.TextFormatter{} },

			// Types (pointer type)
			"NewJSONFormatter": func() *logrus.JSONFormatter { return &logrus.JSONFormatter{} },
			"NewLogger":        func() *logrus.Logger { return &logrus.Logger{} },
			"NewMutexWrap":     func() *logrus.MutexWrap { return &logrus.MutexWrap{} },
			"NewTextFormatter": func() *logrus.TextFormatter { return &logrus.TextFormatter{} },
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
	"fmt"
)

import (
	"github.com/xujinzheng/playgo/example/modules/github.com/denverdino/aliyungo"
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
	aliyungo.Enable(runtime)

	_, err := runtime.RunString(`
		var entryA = logrus.NewEntry()
		var entryB = logrus.Entry()
		
		logrus.Println("entryA:",entryA)
		logrus.Println("entryB:",entryB)
    	
    	logrus.WithField("Hello", "World").Println("I am gojs")

    	// gojs-tool gen --recusive --namespace github.com/denverdino/aliyungo
    	// wrapper packages in namespace, access by dot
    	var client = aliyungo.cs.NewClient() 

    	logrus.Println("client", client)
    `)

	if err != nil {
		fmt.Println(err)
		return
	}
}


```


```bash
example> go run main.go
INFO[0000] entryA: &{<nil> map[] 0001-01-01 00:00:00 +0000 UTC panic  <nil>}
INFO[0000] entryB: {<nil> map[] 0001-01-01 00:00:00 +0000 UTC panic  <nil>}
INFO[0000] I am gojs                                     Hello=World
INFO[0000] client &{  https://cs.aliyuncs.com 2015-12-15 false  0xc420399e60}
```