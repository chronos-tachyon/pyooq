package codegen

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"sort"
	"text/template"

	"github.com/chronos-tachyon/pyooq/repr"
)

//go:embed go.txt
var goTemplateRaw string

var goTemplate = template.Must(
	template.New("root").
		Option("missingkey=error").
		Funcs(goFuncMap).
		Parse(goTemplateRaw))

var goFuncMap = template.FuncMap{
	"SchemaTypeName": func(str string) string {
		var scratch [8]string
		name := Name(scratch[:0]).Extend(str)
		name = append(name, "Schema")
		return name.CamelCase()
	},
	"SchemaVarName": func(str string) string {
		var scratch [8]string
		return Name(scratch[:0]).Extend(str).CamelCase()
	},
	"TableTypeName": func(str string, t *repr.Table) string {
		var scratch [8]string
		name := Name(scratch[:0]).Extend(str)
		name = name.Extend(t.Name)
		name = append(name, "Table")
		return name.CamelCase()
	},
	"TableFieldName": func(t *repr.Table) string {
		var scratch [8]string
		return Name(scratch[:0]).Extend(t.Name).CamelCase()
	},
	"TableSQLName": func(t *repr.Table) string {
		var scratch [8]string
		return Name(scratch[:0]).Extend(t.Name).LowerSnakeCase()
	},
	"ColumnFieldName": func(c *repr.Column) string {
		var scratch [8]string
		return Name(scratch[:0]).Extend(c.Name).CamelCase()
	},
	"ColumnSQLName": func(c *repr.Column) string {
		var scratch [8]string
		return Name(scratch[:0]).Extend(c.Name).LowerSnakeCase()
	},
}

type Go struct{}

type goTemplateData struct {
	Schema      repr.Schema
	FileName    string
	PackageName string
	SchemaName  string
}

func (Go) GenerateCode(s repr.Schema, params map[string]any) ([]GeneratedFile, error) {
	var data goTemplateData
	data.Schema = s
	data.FileName = stringParam(params, "fileName", nil)
	data.PackageName = stringParam(params, "packageName", nil)
	data.SchemaName = stringParam(params, "schemaName", nil)
	noMoreParams(params)

	var contents bytes.Buffer
	err := goTemplate.Execute(&contents, data)
	if err != nil {
		return nil, err
	}

	formatted, err := format.Source(contents.Bytes())
	if err != nil {
		return nil, err
	}

	files := []GeneratedFile{
		{
			Path:     data.FileName,
			Contents: formatted,
		},
	}
	return files, nil
}

var _ CodeGenerator = Go{}

func stringParam(params map[string]any, key string, fn func() string) string {
	if str, found := params[key]; found {
		delete(params, key)
		return str.(string)
	}
	if fn == nil {
		panic(fmt.Errorf("missing required parameter field: %q", key))
	}
	return fn()
}

func noMoreParams(params map[string]any) {
	if len(params) <= 0 {
		return
	}

	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	panic(fmt.Errorf("found unknown parameter fields: %v", keys))
}
