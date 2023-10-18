package codegen

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"sort"
	"text/template"

	"github.com/chronos-tachyon/pyooq/schema"
)

//go:embed go.txt
var goTemplateRaw string

var goTemplate = template.Must(
	template.New("root").
		Delims("@@{{", "}}@@").
		Option("missingkey=error").
		Parse(goTemplateRaw))

type Go struct{}

type goTemplateData struct {
	Schema      schema.Schema
	FileName    string
	PackageName string
	SymbolName  string
}

func (Go) GenerateCode(s schema.Schema, params map[string]any) ([]GeneratedFile, error) {
	var data goTemplateData
	data.Schema = s
	data.FileName = stringParam(params, "fileName", nil)
	data.PackageName = stringParam(params, "packageName", nil)
	data.SymbolName = stringParam(params, "symbolName", nil)
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
