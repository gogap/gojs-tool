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

func main() {

	app := cli.NewApp()

	app.Usage = "A tool for export package funcs, types, vars for VM (like goja) to use"

	app.Version = "1.0.0"

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

	wd, err := os.Getwd()
	if err != nil {
		return
	}

	modulesDir := filepath.Join(wd, "modules")

	for i := 0; i < len(ctx.Args()); i++ {
		pkgPath := ctx.Args()[i]

		if !recusive {
			err = generateModulePackage(template, modulesDir, pkgPath, gopath)
			if err != nil {
				return
			}
		} else {

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

			filepath.Walk(filepath.Join(gopath, "src", pkgPath), walkFn)
		}
	}

	return
}

func generateModulePackage(tmplName, modulesDir, pkgPath, gopath string) (err error) {
	code, err := gojs.GenerateCode(tmplName, pkgPath, "", gopath)
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
