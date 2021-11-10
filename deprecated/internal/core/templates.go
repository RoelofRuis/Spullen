package core

import (
	"fmt"
	"html/template"
)

func NewTemplates() *Templates {
	return &Templates{
		templates: make(map[string]*template.Template),
	}
}

type Templates struct {
	templates map[string]*template.Template
}

func (v *Templates) Register(name string, files ...string) {
	t := template.Must(template.ParseFiles(files...))
	v.templates[name] = t
}

func (v *Templates) Get(name string) (*template.Template, error) {
	t, has := v.templates[name]
	if !has {
		return nil, fmt.Errorf("no template registered with name [%s]", name)
	}
	return t, nil
}
