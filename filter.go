// Package dynamicquery fornece utilitários para filtros dinâmicos em GORM.
package dynamicquery

import (
	"fmt"
	"regexp"
	"strings"
)

// FilterOp representa um operador de comparação em um filtro.
type FilterOp string

const (
	OpEq        FilterOp = "eq"
	OpNeq       FilterOp = "neq"
	OpGt        FilterOp = "gt"
	OpGte       FilterOp = "gte"
	OpLt        FilterOp = "lt"
	OpLte       FilterOp = "lte"
	OpLike      FilterOp = "like"
	OpNotLike   FilterOp = "nlike"
	OpStarts    FilterOp = "starts"
	OpEnds      FilterOp = "ends"
	OpBetween   FilterOp = "between"
	OpIn        FilterOp = "in"
	OpNotIn     FilterOp = "nin"
	OpIsNull    FilterOp = "isnull"
	OpIsNotNull FilterOp = "isnotnull"
)

var filterRx = regexp.MustCompile(
	`^([^:]+):(eq|neq|gt|gte|lt|lte|like|nlike|starts|ends|between|in|nin|isnull|isnotnull)(?::(.*))?$`,
)

// FilterSpec define um filtro dinâmico: Path (coluna ou join),
// Operator (um dos FilterOp) e Value (nil p/ isnull/isnotnull).
type FilterSpec struct {
	Path     []string
	Operator FilterOp
	Value    *string
}

// GetPath retorna o Path do FilterSpec.
func (f FilterSpec) GetPath() []string { return f.Path }

// SetPath redefine o Path do FilterSpec.
func (f *FilterSpec) SetPath(p []string) { f.Path = p }

// ParseFilter analisa raw "campo:op:valor" e retorna um FilterSpec ou erro.
func ParseFilter(raw string) (FilterSpec, error) {
	m := filterRx.FindStringSubmatch(raw)
	if m == nil {
		return FilterSpec{}, fmt.Errorf("dynamicquery: filtro inválido: %q", raw)
	}
	pathStr, opStr, val := m[1], m[2], m[3]
	var vptr *string
	if opStr != string(OpIsNull) && opStr != string(OpIsNotNull) {
		vptr = &val
	}
	return FilterSpec{
		Path:     strings.Split(pathStr, "."),
		Operator: FilterOp(opStr),
		Value:    vptr,
	}, nil
}
