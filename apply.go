package dynamicquery

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// ApplyFilters aplica cada FilterSpec como WHERE no db ou retorna erro.
func ApplyFilters(db *gorm.DB, specs []FilterSpec) (*gorm.DB, error) {
	for _, f := range specs {
		col := strings.Join(f.Path, ".")
		switch f.Operator {
		case OpEq:
			db = db.Where(fmt.Sprintf("%s = ?", col), *f.Value)
		case OpNeq:
			db = db.Where(fmt.Sprintf("%s <> ?", col), *f.Value)
		case OpGt:
			db = db.Where(fmt.Sprintf("%s > ?", col), *f.Value)
		case OpGte:
			db = db.Where(fmt.Sprintf("%s >= ?", col), *f.Value)
		case OpLt:
			db = db.Where(fmt.Sprintf("%s < ?", col), *f.Value)
		case OpLte:
			db = db.Where(fmt.Sprintf("%s <= ?", col), *f.Value)
		case OpLike:
			db = db.Where(fmt.Sprintf("%s LIKE ?", col), "%"+*f.Value+"%")
		case OpNotLike:
			db = db.Where(fmt.Sprintf("%s NOT LIKE ?", col), "%"+*f.Value+"%")
		case OpStarts:
			db = db.Where(fmt.Sprintf("%s LIKE ?", col), *f.Value+"%")
		case OpEnds:
			db = db.Where(fmt.Sprintf("%s LIKE ?", col), "%"+*f.Value)
		case OpBetween:
			parts := strings.SplitN(*f.Value, ",", 2)
			if len(parts) != 2 {
				return db, errors.New("dynamicquery: between precisa de dois valores separados por vírgula")
			}
			db = db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", col), parts[0], parts[1])
		case OpIn:
			vals := strings.Split(*f.Value, ",")
			db = db.Where(fmt.Sprintf("%s IN ?", col), vals)
		case OpNotIn:
			vals := strings.Split(*f.Value, ",")
			db = db.Where(fmt.Sprintf("%s NOT IN ?", col), vals)
		case OpIsNull:
			db = db.Where(fmt.Sprintf("%s IS NULL", col))
		case OpIsNotNull:
			db = db.Where(fmt.Sprintf("%s IS NOT NULL", col))
		default:
			return db, fmt.Errorf("dynamicquery: operador desconhecido %q", f.Operator)
		}
	}
	return db, nil
}

// ApplySorts aplica cada SortSpec como ORDER BY.
func ApplySorts(db *gorm.DB, specs []SortSpec) *gorm.DB {
	for _, s := range specs {
		col := strings.Join(s.Path, ".")
		db = db.Order(fmt.Sprintf("%s %s", col, s.Direction))
	}
	return db
}

// ApplySelects faz SELECT col1,col2…; se specs vazio, não altera o db.
func ApplySelects(db *gorm.DB, specs []SelectSpec) *gorm.DB {
	if len(specs) == 0 {
		return db
	}
	cols := make([]string, len(specs))
	for i, sel := range specs {
		cols[i] = strings.Join(sel.Path, ".")
	}
	return db.Select(strings.Join(cols, ", "))
}
