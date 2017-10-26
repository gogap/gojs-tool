package gojs

type Objects map[string]interface{}

type Runtime interface {
	Set(string, interface{})
}

type Object interface {
	Set(string, interface{})
	Get(string) interface{}
}

type Module interface {
	Name() string
	Set(objects Objects) Module
	Enable(Runtime)
	Register() Module
}

type TemplateVars struct {
	PackageName    string
	PackagePath    string
	PackageObjects map[string]string
}
