package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogap/gojs-tool/gojs"
	"github.com/urfave/cli"
)

type Module struct {
	PackageName string
	PackagePath string
}

func main() {

	app := cli.NewApp()

	app.Usage = "A tool for export package funcs, types, vars for VM (like goja) to use"

	app.Version = "0.0.1"

	app.Commands = cli.Commands{
		cli.Command{
			Name:      "gen",
			Action:    gen,
			Usage:     "generate modules",
			ArgsUsage: "packages",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "template, t",
					Value: "goja",
					Usage: "template filename in templates folder without extension",
				},
				cli.BoolFlag{
					Name:  "recusive, r",
					Usage: "recusive generate code",
				},
				cli.BoolFlag{
					Name:  "namespace, n",
					Usage: "generate submodules loader, access object by dot, only work while recusive enabled",
				},
				cli.StringFlag{
					Name:   "gopath",
					EnvVar: "GOPATH",
					Usage:  "the package in which GOPATH",
				},
			},
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(128)
	}
}

func gen(ctx *cli.Context) (err error) {

	if len(ctx.Args()) == 0 {
		err = errors.New("no package input")
		return
	}

	recusive := ctx.Bool("recusive")
	gopath := ctx.String("gopath")
	template := ctx.String("template")
	namespace := ctx.Bool("namespace")

	wd, err := os.Getwd()
	if err != nil {
		return
	}

	modulesDir := filepath.Join(wd, "modules")

	if !recusive {
		for i := 0; i < len(ctx.Args()); i++ {
			pkgPath := ctx.Args()[i]

			err = generateModulePackage(template, modulesDir, pkgPath, gopath)
			if err != nil {
				return
			}
		}

		return
	}

	for i := 0; i < len(ctx.Args()); i++ {
		pkgPath := ctx.Args()[i]

		walkFn := func(path string, info os.FileInfo, err error) error {

			if info == nil {
				return errors.New("read path error:" + path)
			}

			if !info.IsDir() {
				return nil
			}

			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}

			matched, e := filepath.Glob(filepath.Join(path, "*.go"))
			if e != nil {
				return e
			}

			if len(matched) == 0 {
				return nil
			}

			e = generateModulePackage(template, modulesDir, strings.TrimPrefix(path, filepath.Join(gopath, "src")+"/"), gopath)

			if e != nil {
				return e
			}

			return nil
		}

		err = filepath.Walk(filepath.Join(gopath, "src", pkgPath), walkFn)
		if err != nil {
			return
		}
	}

	if namespace {

		for i := 0; i < len(ctx.Args()); i++ {

			pkgPath := ctx.Args()[i]

			pkgModules := make(map[string][]string)

			walkFn := func(path string, info os.FileInfo, err error) error {

				if info == nil {
					return errors.New("read path error:" + path)
				}

				if !info.IsDir() {
					return nil
				}

				if strings.HasPrefix(info.Name(), ".") {
					return filepath.SkipDir
				}

				files, e := ioutil.ReadDir(path)
				if e != nil {
					return e
				}

				for j := 0; j < len(files); j++ {

					if files[j] == nil {
						return errors.New("read path error:" + filepath.Join(path, files[j].Name()))
					}

					if !files[j].IsDir() {
						continue
					}

					if strings.HasPrefix(files[j].Name(), ".") {
						continue
					}

					hasGoFiles, e := filepath.Glob(filepath.Join(path, "*.go"))
					if e != nil {
						return e
					}

					if len(hasGoFiles) > 0 {
						return nil
					}

					pkgModules[path] = append(pkgModules[path], files[j].Name())
				}

				return nil
			}

			err = filepath.Walk(filepath.Join(modulesDir, pkgPath), walkFn)
			if err != nil {
				return
			}

			for path, moduleNames := range pkgModules {
				pkgPath := strings.TrimPrefix(path, modulesDir+"/")
				rootPath := strings.TrimPrefix(path, filepath.Join(gopath, "src")+"/")

				var modules []Module
				for _, moduleName := range moduleNames {
					m := Module{
						PackagePath: filepath.Join(rootPath, moduleName),
						PackageName: moduleName,
					}

					modules = append(modules, m)
				}

				err = generateNamespace(template+"_namespace", modulesDir, pkgPath, gopath, modules)
				if err != nil {
					return
				}
			}
		}
	}

	return
}

func generateModulePackage(tmplName, modulesDir, pkgPath, gopath string) (err error) {

	vars, err := gojs.Parser(pkgPath, gopath)
	if err != nil {
		return
	}

	options := gojs.GenerateOptions{
		TemplateName: tmplName,
		PackagePath:  pkgPath,
		GoPath:       gopath,
		TemplateVars: vars,
	}

	codeDir := filepath.Join(modulesDir, pkgPath)

	code, err := gojs.GenerateCode(options)
	if err != nil {
		return
	}

	err = os.MkdirAll(codeDir, 0755)
	if err != nil {
		return
	}

	codeFile := filepath.Join(codeDir, filepath.Base(pkgPath)+".go")

	err = ioutil.WriteFile(codeFile, []byte(code), 0644)
	if err != nil {
		return
	}

	return
}

func generateNamespace(tmplName, modulesDir, pkgPath, gopath string, modules []Module) (err error) {

	options := gojs.GenerateOptions{
		TemplateName: tmplName,
		PackagePath:  pkgPath,
		GoPath:       gopath,
		TemplateVars: &gojs.TemplateVars{
			PackageName: filepath.Base(pkgPath),
			PackagePath: strings.TrimPrefix(pkgPath, filepath.Join(gopath, "src")+"/"),
			Args: map[string]interface{}{
				"SubModules": modules,
			},
		},
	}

	code, err := gojs.GenerateCode(options)
	if err != nil {
		return
	}

	codeDir := filepath.Join(modulesDir, pkgPath)

	err = os.MkdirAll(codeDir, 0755)
	if err != nil {
		return
	}

	codeFile := filepath.Join(codeDir, filepath.Base(pkgPath)+".go")

	err = ioutil.WriteFile(codeFile, []byte(code), 0644)
	if err != nil {
		return
	}

	return
}
