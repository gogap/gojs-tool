package gojs

import (
	"bytes"
	"go/format"

	"github.com/gogap/gocoder"
	"github.com/xujinzheng/gojs-tool/gojs/templates"
	"text/template"
)

func Parser(pkgPath string, gopath string) (vars *TemplateVars, err error) {
	pkg, err := gocoder.NewGoPackage(pkgPath,
		gocoder.OptionGoPath(gopath),
	)

	if err != nil {
		return
	}

	tmplVars := &TemplateVars{
		PackageName:    pkg.Name(),
		PackagePath:    pkg.Path(),
		PackageObjects: make(map[string]string),
	}

	numFuncs := pkg.NumFuncs()

	for i := 0; i < numFuncs; i++ {
		goFunc := pkg.Func(i)

		if !isExported(goFunc.Name()) {
			continue
		}

		if len(goFunc.Receiver()) > 0 {
			continue
		}

		tmplVars.PackageObjects[goFunc.Name()] = goFunc.Name()
	}

	numVars := pkg.NumVars()

	for i := 0; i < numVars; i++ {
		goVar := pkg.Var(i)

		if !isExported(goVar.Name()) {
			continue
		}

		tmplVars.PackageObjects[goVar.Name()] = goVar.Name()
	}

	vars = tmplVars

	return
}

func GenerateGojaModule(pkgPath string, gopath string) (code string, err error) {
	vars, err := Parser(pkgPath, gopath)
	if err != nil {
		return
	}

	tmpl, err := template.New("Goja").Parse(templates.Goja)
	if err != nil {
		return
	}

	buf := bytes.NewBuffer(nil)
	err = tmpl.Execute(buf, vars)

	if err != nil {
		return
	}

	codeBytes, err := format.Source(buf.Bytes())

	if err != nil {
		return
	}

	code = string(codeBytes)

	return
}

func isExported(v string) bool {
	if len(v) == 0 {
		return false
	}

	if v[0] >= 'A' && v[0] <= 'Z' {
		return true
	}

	return false
}
