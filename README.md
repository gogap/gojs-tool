# gojs-tool

A tool for export package funcs, types, vars for VM (like goja and otto) to use

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

example> gojs-tool gen -t otto -r -n github.com/sirupsen/logrus ## for otto
example> gojs-tool gen -t goja -r -n github.com/sirupsen/logrus ## for goja

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

#### Auto generated code (goja)

`modules/github.com/sirupsen/logrus/logrus.go`

```go
package logrus

import (
	original_logrus "github.com/sirupsen/logrus"
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
			"AddHook":             original_logrus.AddHook,
			"Debug":               original_logrus.Debug,
			"Debugf":              original_logrus.Debugf,
			"Debugln":             original_logrus.Debugln,
			"Error":               original_logrus.Error,
			"Errorf":              original_logrus.Errorf,
			"Errorln":             original_logrus.Errorln,
			"Exit":                original_logrus.Exit,
			"Fatal":               original_logrus.Fatal,
			//..........more funcs..........

			// Var and consts
			"AllLevels":     original_logrus.AllLevels,
			"DebugLevel":    original_logrus.DebugLevel,
			"ErrorKey":      original_logrus.ErrorKey,
			"ErrorLevel":    original_logrus.ErrorLevel,
			"FatalLevel":    original_logrus.FatalLevel,
			"FieldKeyLevel": original_logrus.FieldKeyLevel,
			"FieldKeyMsg":   original_logrus.FieldKeyMsg,
			"FieldKeyTime":  original_logrus.FieldKeyTime,
			"InfoLevel":     original_logrus.InfoLevel,
			"PanicLevel":    original_logrus.PanicLevel,
			"WarnLevel":     original_logrus.WarnLevel,

			// Types (value type)
			"Entry":         func() original_logrus.Entry { return original_logrus.Entry{} },
			"JSONFormatter": func() original_logrus.JSONFormatter { return original_logrus.JSONFormatter{} },
			"Logger":        func() original_logrus.Logger { return original_logrus.Logger{} },
			"MutexWrap":     func() original_logrus.MutexWrap { return original_logrus.MutexWrap{} },
			"TextFormatter": func() original_logrus.TextFormatter { return original_logrus.TextFormatter{} },

			// Types (pointer type)
			"NewJSONFormatter": func() *original_logrus.JSONFormatter { return &original_logrus.JSONFormatter{} },
			"NewLogger":        func() *original_logrus.Logger { return &original_logrus.Logger{} },
			"NewMutexWrap":     func() *original_logrus.MutexWrap { return &original_logrus.MutexWrap{} },
			"NewTextFormatter": func() *original_logrus.TextFormatter { return &original_logrus.TextFormatter{} },
		},
	).Register()
}

func Enable(runtime *goja.Runtime) {
	module.Enable(runtime)
}


```

#### Auto generated code (otto)

```go
package logrus

import (
	original_logrus "github.com/sirupsen/logrus"
)

import (
	"github.com/robertkrimen/otto"
)

type Logrus struct {
	// Functions
	AddHook             interface{}
	Debug               interface{}
	Debugf              interface{}
	Debugln             interface{}
	Error               interface{}
   //..........more funcs..........

	// Var and consts
	AllLevels     interface{}
	DebugLevel    interface{}
	ErrorKey      interface{}
	ErrorLevel    interface{}
	//..........more vars..........

	// Types (value type)
	Entry         interface{}
	JSONFormatter interface{}
	Logger        interface{}
	MutexWrap     interface{}
	TextFormatter interface{}

	// Types (pointer type)
	NewJSONFormatter interface{}
	NewLogger        interface{}
	NewMutexWrap     interface{}
	NewTextFormatter interface{}
}

var (
	module = &Logrus{
		// Functions
		AddHook:             original_logrus.AddHook,
		Debug:               original_logrus.Debug,
		Debugf:              original_logrus.Debugf,
		Debugln:             original_logrus.Debugln,
		//..........more funcs..........

		// Var and consts
		AllLevels:     original_logrus.AllLevels,
		DebugLevel:    original_logrus.DebugLevel,
		ErrorKey:      original_logrus.ErrorKey,
		ErrorLevel:    original_logrus.ErrorLevel,
		//..........more vars..........

		// Types (value type)
		Entry:         func() original_logrus.Entry { return original_logrus.Entry{} },
		JSONFormatter: func() original_logrus.JSONFormatter { return original_logrus.JSONFormatter{} },
		Logger:        func() original_logrus.Logger { return original_logrus.Logger{} },
		MutexWrap:     func() original_logrus.MutexWrap { return original_logrus.MutexWrap{} },
		TextFormatter: func() original_logrus.TextFormatter { return original_logrus.TextFormatter{} },

		// Types (pointer type)
		NewJSONFormatter: func() *original_logrus.JSONFormatter { return &original_logrus.JSONFormatter{} },
		NewLogger:        func() *original_logrus.Logger { return &original_logrus.Logger{} },
		NewMutexWrap:     func() *original_logrus.MutexWrap { return &original_logrus.MutexWrap{} },
		NewTextFormatter: func() *original_logrus.TextFormatter { return &original_logrus.TextFormatter{} },
	}
)

func Enable(o *otto.Otto) {
	o.Set("logrus", module)
}

```


#### Use native module in goja

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

    	// gojs-tool gen --template goja --recusive --namespace github.com/denverdino/aliyungo
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


#### Use native module in otto

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
	"github.com/robertkrimen/otto"
)

func main() {

	runtime := otto.New()

	logrus.Enable(runtime)
	aliyungo.Enable(runtime)

	_, err := runtime.Run(`
		var entryA = logrus.NewEntry(undefined)
		var entryB = logrus.Entry()
		
		logrus.Println("entryA:",entryA)
		logrus.Println("entryB:",entryB)
    	
    	logrus.WithField("Hello", "World").Println("I am gojs")

    	// gojs-tool gen --template otto --recusive --namespace github.com/denverdino/aliyungo
    	// wrapper packages in namespace, access by dot
    	var client = aliyungo.cs.NewClient("","") 

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
INFO[0000] client &{  https://cs.aliyuncs.com 2015-12-15 false  0xc4201824b0}
```