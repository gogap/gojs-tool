package gojs

import (
	"bytes"
	"fmt"
	"go/format"
	"text/template"

	"github.com/gogap/gocoder"
	"github.com/gogap/gojs-tool/gojs/templates"
)

type GenerateOptions struct {
	TemplateName string
	PackagePath  string
	PackageAlias string
	GoPath       string
	ProjectPath  string

	TemplateVars *TemplateVars

	Args interface{}
}

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

		if goType.IsStruct() {
			tmplVars.PackageTypes[goType.Name()] = goType.Name()
		}

	}

	vars = tmplVars

	return
}

func GenerateCode(options GenerateOptions) (code string, err error) {

	tmplBytes, err := templates.Asset(options.TemplateName + ".tmpl")
	if err != nil {
		return
	}

	tmpl, err := template.New(options.TemplateName).Funcs(templateFuncs()).Parse(string(tmplBytes))
	if err != nil {
		return
	}

	buf := bytes.NewBuffer(nil)

	err = tmpl.Execute(buf, options.TemplateVars)

	if err != nil {
		return
	}

	if len(options.PackageAlias) > 0 {
		options.TemplateVars.PackageName = options.PackageAlias
	}

	codeBytes, err := format.Source(buf.Bytes())

	if err != nil {
		fmt.Println(buf.String())
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

func templateFuncs() map[string]interface{} {
	return map[string]interface{}{
		"exist": func(v map[string]string, key string) bool {
			_, exist := v[key]
			return exist
		},
	}
}
