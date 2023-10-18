package codegen

import (
	"github.com/chronos-tachyon/pyooq/schema"
)

type CodeGenerator interface {
	GenerateCode(s schema.Schema, params map[string]any) ([]GeneratedFile, error)
}

type GeneratedFile struct {
	Path     string
	Contents []byte
}
