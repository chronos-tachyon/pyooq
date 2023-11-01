package codegen

import (
	"fmt"
	"regexp"
	"strings"
)

var reSep = regexp.MustCompile(`[\s_-]+`)

func concat(list ...any) []string {
	var n uint
	for _, item := range list {
		switch x := item.(type) {
		case string:
			n++
		case []string:
			n += uint(len(x))
		default:
			panic(fmt.Errorf("unexpected type %T", item))
		}
	}
	if n <= 0 {
		return nil
	}

	out := make([]string, 0, n)
	for _, item := range list {
		switch x := item.(type) {
		case string:
			out = append(out, x)
		case []string:
			out = append(out, x...)
		}
	}
	return out
}

func equalLists(a []string, b []string) bool {
	an := uint(len(a))
	bn := uint(len(b))
	if an != bn {
		return false
	}
	for i := uint(0); i < an; i++ {
		ai, bi := a[i], b[i]
		if ai != bi {
			return false
		}
	}
	return true
}

func splitName(str string) []string {
	return reSep.Split(str, -1)
}

func joinName(name []string, sep string, mapfn func(string) string) string {
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

func camelCaseName(name []string) string {
	return joinName(name, "", nil)
}

func lowerCamelCaseName(name []string) string {
	isFirst := true
	return joinName(name, "", func(str string) string {
		if isFirst {
			isFirst = false
			return strings.ToLower(str)
		}
		return str
	})
}

func snakeCaseName(name []string) string {
	return joinName(name, "_", nil)
}

func upperSnakeCaseName(name []string) string {
	return joinName(name, "_", strings.ToUpper)
}

func lowerSnakeCaseName(name []string) string {
	return joinName(name, "_", strings.ToLower)
}

func kebabCaseName(name []string) string {
	return joinName(name, "-", nil)
}

func upperKebabCaseName(name []string) string {
	return joinName(name, "-", strings.ToUpper)
}

func lowerKebabCaseName(name []string) string {
	return joinName(name, "-", strings.ToLower)
}
