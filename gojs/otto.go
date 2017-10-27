package gojs

import (
	"github.com/robertkrimen/otto"
)

type OttoModule struct {
	name string
	sets map[string]interface{}

	runtime *otto.Otto
}

func NewOttoModule(name string) Module {
	return &OttoModule{
		name: name,
		sets: make(map[string]interface{}),
	}
}

func (p *OttoModule) String() string {
	return p.name
}

func (p *OttoModule) Name() string {
	return p.name
}

func (p *OttoModule) Register() Module {
	return p
}

func (p *OttoModule) Set(objects Objects) Module {

	for k, v := range objects {
		p.sets[k] = v
	}

	return p
}

func (p *OttoModule) Enable(runtime Runtime) {
	for k, v := range p.sets {
		runtime.Set(k, v)
	}
}
