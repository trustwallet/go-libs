package client

import "fmt"

type Path struct {
	template string
	values   []any
}

func NewStaticPath(path string) Path {
	return Path{template: path}
}

func NewEmptyPath() Path {
	return Path{}
}

func NewPath(template string, values []any) Path {
	return Path{template: template, values: values}
}

func (p Path) String() string {
	return fmt.Sprintf(p.template, p.values...)
}
