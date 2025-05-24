package dynamicquery

import "strings"

// ApplyFilterAliases percorre cada FilterSpec e, se o Path tiver
// só um segmento, substitui por aliasMap[seg]
// ou, se não houver alias, converte CamelCase→snake_case.
func ApplyFilterAliases(specs []FilterSpec, aliasMap map[string]string) []FilterSpec {
	out := make([]FilterSpec, 0, len(specs))
	for _, spec := range specs {
		path := spec.Path
		if len(path) == 1 {
			if alias, ok := aliasMap[path[0]]; ok {
				path = strings.Split(alias, ".")
			} else {
				path[0] = toSnakeCase(path[0])
			}
		}
		spec.Path = path
		out = append(out, spec)
	}
	return out
}

// ApplySortAliases faz o mesmo para SortSpec.
func ApplySortAliases(specs []SortSpec, aliasMap map[string]string) []SortSpec {
	out := make([]SortSpec, 0, len(specs))
	for _, spec := range specs {
		path := spec.Path
		if len(path) == 1 {
			if alias, ok := aliasMap[path[0]]; ok {
				path = strings.Split(alias, ".")
			} else {
				path[0] = toSnakeCase(path[0])
			}
		}
		spec.Path = path
		out = append(out, spec)
	}
	return out
}

// ApplySelectAliases faz o mesmo para SelectSpec.
func ApplySelectAliases(specs []SelectSpec, aliasMap map[string]string) []SelectSpec {
	out := make([]SelectSpec, 0, len(specs))
	for _, spec := range specs {
		path := spec.Path
		if len(path) == 1 {
			if alias, ok := aliasMap[path[0]]; ok {
				path = strings.Split(alias, ".")
			} else {
				path[0] = toSnakeCase(path[0])
			}
		}
		spec.Path = path
		out = append(out, spec)
	}
	return out
}
