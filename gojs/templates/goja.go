package templates

const Goja = `
package {{.PackageName}}

import (
	"{{.PackagePath}}"
)

import (
	"github.com/dop251/goja"
	"github.com/xujinzheng/gojs-tool/gojs"
)

var (
	module = gojs.NewGojaModule("{{.PackageName}}")
)

func init() {
	module.Set(
		gojs.Objects{
			{{- $pkgName:=.PackageName -}}
			{{range $objName, $objDefine := .PackageObjects}}
			"{{$objName}}": {{$pkgName}}.{{$objDefine}},
			{{- end}}
			},
		).Register()
}

func Enable(runtime *goja.Runtime) {
	module.Enable(runtime)
}
`
