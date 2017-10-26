package gojs

import (
	"bytes"
	"go/format"
	"text/template"

	"github.com/gogap/gocoder"
	"github.com/gogap/gojs-tool/gojs/templates"
)

func Parser(pkgPath string, gopath string) (vars *TemplateVars, err error) {
	pkg, err := gocoder.NewGoPackage(pkgPath,
		gocoder.OptionGoPath(gopath),
	)

	if err != nil {
		return
	}

	tmplVars := &TemplateVars{
		PackageName:  pkg.Name(),
		PackagePath:  pkg.Path(),
		PackageFuncs: make(map[string]string),
		PackageTypes: make(map[string]string),
		PackageVars:  make(map[string]string),
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

		tmplVars.PackageFuncs[goFunc.Name()] = goFunc.Name()
	}

	numVars := pkg.NumVars()

	for i := 0; i < numVars; i++ {
		goVar := pkg.Var(i)

		if !isExported(goVar.Name()) {
			continue
		}

		tmplVars.PackageVars[goVar.Name()] = goVar.Name()
	}

	numTypes := pkg.NumTypes()

	for i := 0; i < numTypes; i++ {
		goType := pkg.Type(i)

		if !isExported(goType.Name()) {
			continue
		}

		tmplVars.PackageTypes[goType.Name()] = goType.Name()
	}

	vars = tmplVars

	return
}

func GenerateCode(tmplName, pkgPath, pkgAlias, gopath string) (code string, err error) {
	vars, err := Parser(pkgPath, gopath)
	if err != nil {
		return
	}

	tmplBytes, err := templates.Asset(tmplName + ".tmpl")
	if err != nil {
		return
	}

	tmpl, err := template.New("Goja").Parse(string(tmplBytes))
	if err != nil {
		return
	}

	buf := bytes.NewBuffer(nil)
	err = tmpl.Execute(buf, vars)

	if err != nil {
		return
	}

	if len(pkgAlias) > 0 {
		vars.PackageName = pkgAlias
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
