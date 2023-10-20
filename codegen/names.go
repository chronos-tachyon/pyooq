package codegen

import (
	"regexp"
	"strings"
)

var reSep = regexp.MustCompile(`[\s_-]+`)

type Name []string

func (name Name) Extend(str string) Name {
	pieces := reSep.Split(str, -1)
	return append(name, pieces...)
}

func (name Name) Join(sep string, mapfn func(string) string) string {
	var out []byte
	needSep := false
	for _, str := range name {
		if mapfn != nil {
			str = mapfn(str)
		}
		if str == "" {
			continue
		}
		if needSep {
			out = append(out, sep...)
		}
		out = append(out, str...)
		needSep = true
	}
	return string(out)
}

func (name Name) CamelCase() string {
	return name.Join("", nil)
}

func (name Name) LowerCamelCase() string {
	isFirst := true
	return name.Join("", func(str string) string {
		if isFirst {
			isFirst = false
			return strings.ToLower(str)
		}
		return str
	})
}

func (name Name) SnakeCase() string {
	return name.Join("_", nil)
}

func (name Name) UpperSnakeCase() string {
	return name.Join("_", strings.ToUpper)
}

func (name Name) LowerSnakeCase() string {
	return name.Join("_", strings.ToLower)
}

func (name Name) KebabCase() string {
	return name.Join("-", nil)
}

func (name Name) UpperKebabCase() string {
	return name.Join("-", strings.ToUpper)
}

func (name Name) LowerKebabCase() string {
	return name.Join("-", strings.ToLower)
}
