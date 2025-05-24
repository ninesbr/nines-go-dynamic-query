package dynamicquery

import (
	"fmt"
	"regexp"
	"strings"
)

var selectRx = regexp.MustCompile(`^[a-zA-Z0-9_]+(?:\.[a-zA-Z0-9_]+)*$`)

// SelectSpec define projeção de colunas dinâmicas.
type SelectSpec struct {
	Path []string
}

// GetPath retorna o Path do SelectSpec.
func (s SelectSpec) GetPath() []string { return s.Path }

// SetPath redefine o Path do SelectSpec.
func (s *SelectSpec) SetPath(p []string) { s.Path = p }

// ParseSelect valida raw "coluna" ou "join.col" e retorna SelectSpec.
func ParseSelect(raw string) (SelectSpec, error) {
	if !selectRx.MatchString(raw) {
		return SelectSpec{}, fmt.Errorf("dynamicquery: select inválido: %q", raw)
	}
	return SelectSpec{Path: strings.Split(raw, ".")}, nil
}
