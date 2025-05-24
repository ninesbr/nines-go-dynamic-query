package dynamicquery

import (
	"fmt"
	"regexp"
	"strings"
)

var sortRx = regexp.MustCompile(`^([^.]+(?:\.[^.]+)*)\.(asc|desc)$`)

// SortSpec representa ordenação dinâmica: Path e Direction ("ASC"/"DESC").
type SortSpec struct {
	Path      []string
	Direction string
}

// GetPath retorna o Path do SortSpec.
func (s SortSpec) GetPath() []string { return s.Path }

// SetPath redefine o Path do SortSpec.
func (s *SortSpec) SetPath(p []string) { s.Path = p }

// ParseSort analisa raw "campo.asc" ou "campo.desc" e retorna SortSpec.
func ParseSort(raw string) (SortSpec, error) {
	m := sortRx.FindStringSubmatch(raw)
	if m == nil {
		return SortSpec{}, fmt.Errorf("dynamicquery: sort inválido: %q", raw)
	}
	return SortSpec{
		Path:      strings.Split(m[1], "."),
		Direction: strings.ToUpper(m[2]),
	}, nil
}
