package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/chronos-tachyon/pyooq/codegen"
	"github.com/chronos-tachyon/pyooq/repr"
)

var reGeneratedCode = regexp.MustCompile(`(?i)generated\s+code`)

func main() {
	var flagSchema string
	var flagLang string
	var flagParams string
	var flagOutputDir string
	flag.StringVar(&flagSchema, "schema", "", "path to JSON describing the schema")
	flag.StringVar(&flagLang, "lang", "go", "language to generate; must be 'go'")
	flag.StringVar(&flagParams, "params", "{}", "additional lang-specific parameters as JSON object")
	flag.StringVar(&flagOutputDir, "output_dir", ".", "path to directory containing generated files")
	flag.Parse()

	if flagLang != "go" {
		panic(fmt.Errorf("--lang=%q is not implemented", flagLang))
	}

	raw, err := os.ReadFile(flagSchema)
	if err != nil {
		panic(err)
	}

	var s repr.Schema
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	d.UseNumber()
	err = d.Decode(&s)
	if err != nil {
		panic(err)
	}

	var params map[string]any
	d = json.NewDecoder(strings.NewReader(flagParams))
	d.DisallowUnknownFields()
	d.UseNumber()
	err = d.Decode(&params)
	if err != nil {
		panic(err)
	}

	var g codegen.CodeGenerator = codegen.Go{}
	files, err := g.GenerateCode(s, params)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		filePath := filepath.Join(flagOutputDir, filepath.FromSlash(file.Path))

		contents, err := os.ReadFile(filePath)
		ok := false
		switch {
		case err == nil:
			match := reGeneratedCode.Find(contents)
			ok = (match != nil)
		case errors.Is(err, fs.ErrNotExist):
			ok = true
		default:
			panic(err)
		}

		if !ok {
			panic(fmt.Errorf("refusing to overwrite existing file %q", filePath))
		}

		os.WriteFile(filePath, file.Contents, 0o666)
	}
}
