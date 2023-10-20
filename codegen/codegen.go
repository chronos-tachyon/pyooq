package codegen

import (
	"github.com/chronos-tachyon/pyooq/repr"
)

type CodeGenerator interface {
	GenerateCode(s repr.Schema, params map[string]any) ([]GeneratedFile, error)
}

type GeneratedFile struct {
	Path     string
	Contents []byte
}
