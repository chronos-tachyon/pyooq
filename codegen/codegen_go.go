package codegen

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/chronos-tachyon/pyooq"
	"github.com/chronos-tachyon/pyooq/repr"
)

type Go struct{}

//go:embed go.txt
var goTemplateRaw string

type Renderable interface {
	Render([]byte, string) []byte
}

func (Go) GenerateCode(s *repr.Schema, params map[string]any) ([]GeneratedFile, error) {
	t, err := template.New("root").
		Option("missingkey=error").
		Funcs(template.FuncMap{
			"Render": func(r Renderable, sep string) string {
				return string(r.Render(nil, sep))
			},
			"Quote": func(str string) string {
				return strconv.Quote(str)
			},
		}).
		Parse(goTemplateRaw)
	if err != nil {
		return nil, fmt.Errorf("BUG: failed to parse internal template: %w", err)
	}

	fileName := stringParam(params, "fileName", nil)
	packageName := stringParam(params, "packageName", nil)
	schemaName := stringParam(params, "schemaName", nil)
	noMoreParams(params)

	var gs goSchema
	err = gs.Set(s, schemaName, packageName, fileName)
	if err != nil {
		return nil, err
	}

	var contents bytes.Buffer
	err = t.Execute(&contents, &gs)
	if err != nil {
		return nil, fmt.Errorf("BUG: failed to execute internal template: %w", err)
	}

	raw := contents.Bytes()
	formatted, err := format.Source(raw)
	if err != nil {
		tmp := string(raw)
		tmp = strings.TrimRight(tmp, "\n")
		tmp = strings.ReplaceAll(tmp, "\n", "\n\t|")
		return nil, fmt.Errorf("BUG: gofmt failed: %w\n\t|%s", err, tmp)
	}

	files := []GeneratedFile{
		{
			Path:     fileName,
			Contents: formatted,
		},
	}
	return files, nil
}

var _ CodeGenerator = Go{}

type goSchema struct {
	FileName    string
	PackageName string
	Var         string
	Tables      []*goTable
}

func (gs *goSchema) Set(s *repr.Schema, schemaName string, packageName string, fileName string) error {
	name := splitName(schemaName)
	*gs = goSchema{
		FileName:    fileName,
		PackageName: packageName,
		Var:         camelCaseName(name),
		Tables:      make([]*goTable, 0, len(s.Tables)),
	}
	for _, t := range s.Tables {
		gt := &goTable{}
		gs.Tables = append(gs.Tables, gt)
		if err := gt.Set(gs, t); err != nil {
			return fmt.Errorf("%s: %w", gs.Var, err)
		}
	}
	return nil
}

func (gs *goSchema) Render(out []byte, sep string) []byte {
	out = fmt.Appendf(out, "pyooq.BuildSchema()")
	for _, gt := range gs.Tables {
		out = gt.Render(out, sep)
	}
	out = fmt.Appendf(out, ".%sBuild()", sep)
	return out
}

type goTable struct {
	SQL     string
	Var     string
	Columns []*goColumn
	Indices []*goIndex
}

func (gt *goTable) Set(gs *goSchema, t *repr.Table) error {
	name := splitName(t.Name)
	*gt = goTable{
		SQL:     lowerSnakeCaseName(name),
		Var:     gs.Var + "_" + camelCaseName(name),
		Columns: make([]*goColumn, 0, len(t.Columns)),
		Indices: make([]*goIndex, 0, len(t.Indices)),
	}
	for _, c := range t.Columns {
		gc := &goColumn{}
		gt.Columns = append(gt.Columns, gc)
		if err := gc.Set(gs, gt, c); err != nil {
			return fmt.Errorf("%s: %w", gt.Var, err)
		}
	}
	for _, i := range t.Indices {
		gi := &goIndex{}
		gt.Indices = append(gt.Indices, gi)
		if err := gi.Set(gs, gt, i); err != nil {
			return fmt.Errorf("%s: %w", gt.Var, err)
		}
	}
	return nil
}

func (gt *goTable) Ref() string {
	return fmt.Sprintf("ByName(%q)", gt.SQL)
}

func (gt *goTable) Render(out []byte, sep string) []byte {
	out = fmt.Appendf(out, ".%sTable(%q)", sep, gt.SQL)
	for _, gc := range gt.Columns {
		out = gc.Render(out, sep)
	}
	for _, gi := range gt.Indices {
		out = gi.Render(out, sep)
	}
	out = fmt.Appendf(out, ".%sBuildTable()", sep)
	return out
}

type goColumn struct {
	SQL  string
	Var  string
	Type pyooq.Type
}

func (gc *goColumn) Set(gs *goSchema, gt *goTable, c *repr.Column) error {
	name := splitName(c.Name)
	*gc = goColumn{
		SQL:  lowerSnakeCaseName(name),
		Var:  gt.Var + "_" + camelCaseName(name),
		Type: c.Type,
	}
	return nil
}

func (gc *goColumn) Ref() string {
	return fmt.Sprintf("ByName(%q)", gc.SQL)
}

func (gc *goColumn) Render(out []byte, sep string) []byte {
	return fmt.Appendf(out, ".%sColumn(%q, %#v)", sep, gc.SQL, gc.Type)
}

type goIndex struct {
	SQL     string
	Var     string
	Type    pyooq.IndexType
	Columns []*goIndexedColumn
}

func (gi *goIndex) Set(gs *goSchema, gt *goTable, i *repr.Index) error {
	switch itype := i.Type; itype {
	case pyooq.Index_PrimaryKey:
		*gi = goIndex{
			SQL:     "*PK*",
			Var:     gt.Var + "_PK",
			Type:    itype,
			Columns: make([]*goIndexedColumn, 0, len(i.Columns)),
		}

	default:
		name := splitName(i.Name)
		*gi = goIndex{
			SQL:     lowerSnakeCaseName(name),
			Var:     gt.Var + "_" + camelCaseName(name) + "Index",
			Type:    itype,
			Columns: make([]*goIndexedColumn, 0, len(i.Columns)),
		}
	}
	for _, ic := range i.Columns {
		gic := &goIndexedColumn{}
		gi.Columns = append(gi.Columns, gic)
		if err := gic.Set(gs, gt, gi, ic); err != nil {
			return fmt.Errorf("%s: %w", gi.Var, err)
		}
	}
	return nil
}

func (gi *goIndex) Ref() string {
	switch gi.Type {
	case pyooq.Index_PrimaryKey:
		return "PrimaryKey()"
	default:
		return fmt.Sprintf("IndexByName(%q)", gi.SQL)
	}
}

func (gi *goIndex) Render(out []byte, sep string) []byte {
	switch gi.Type {
	case pyooq.Index_PrimaryKey:
		out = fmt.Appendf(out, ".%sPrimaryKey()", sep)
	case pyooq.Index_Unique:
		out = fmt.Appendf(out, ".%sUnique(%q)", sep, gi.SQL)
	default:
		out = fmt.Appendf(out, ".%sIndex(%q)", sep, gi.SQL)
	}
	for _, gic := range gi.Columns {
		out = gic.Render(out, sep)
	}
	out = fmt.Appendf(out, ".BuildIndex()")
	return out
}

type goIndexedColumn struct {
	Column string
	IsDesc bool
}

func (gic *goIndexedColumn) Set(gs *goSchema, gt *goTable, gi *goIndex, ic *repr.IndexedColumn) error {
	name := splitName(ic.Column)
	colSQL := lowerSnakeCaseName(name)
	found := false
	for _, gc := range gt.Columns {
		if colSQL == gc.SQL {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("column %q does not exist", lowerSnakeCaseName(name))
	}
	*gic = goIndexedColumn{
		Column: colSQL,
		IsDesc: ic.IsDesc,
	}
	return nil
}

func (gic *goIndexedColumn) Render(out []byte, sep string) []byte {
	out = fmt.Appendf(out, ".Add(%q", gic.Column)
	if gic.IsDesc {
		out = fmt.Appendf(out, ", pyooq.Desc(true)")
	}
	out = fmt.Appendf(out, ")")
	return out
}

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
