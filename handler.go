// handler.go
package dynamicquery

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// PageMeta contém metadados de paginação.
type PageMeta struct {
	Page            int  `json:"page"`
	Take            int  `json:"take"`
	ItemCount       int  `json:"itemCount"`
	PageCount       int  `json:"pageCount"`
	HasPreviousPage bool `json:"hasPreviousPage"`
	HasNextPage     bool `json:"hasNextPage"`
}

// QueryHandlerConfig configura o handler genérico.
type QueryHandlerConfig[T any] struct {
	DB           *gorm.DB
	Model        T
	AliasMap     map[string]string
	DefaultLimit int
	MaxLimit     int
}

// NewQueryHandler retorna um http.HandlerFunc para expor consultas dinâmicas.
func NewQueryHandler[T any](cfg QueryHandlerConfig[T]) http.HandlerFunc {
	if cfg.DefaultLimit <= 0 {
		cfg.DefaultLimit = 100
	}
	if cfg.MaxLimit < cfg.DefaultLimit {
		cfg.MaxLimit = cfg.DefaultLimit
	}

	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		// paginação
		page, _ := strconv.Atoi(q.Get("page"))
		if page < 0 {
			page = 0
		}
		take, _ := strconv.Atoi(q.Get("take"))
		if take <= 0 || take > cfg.MaxLimit {
			take = cfg.DefaultLimit
		}
		offset := page * take

		// parse de specs
		var (
			filters []FilterSpec
			sorts   []SortSpec
			selects []SelectSpec
			err     error
		)
		for _, raw := range q["filter"] {
			var spec FilterSpec
			if spec, err = ParseFilter(raw); err != nil {
				http.Error(w, "filtro inválido: "+err.Error(), http.StatusBadRequest)
				return
			}
			filters = append(filters, spec)
		}
		for _, raw := range q["sort"] {
			var spec SortSpec
			if spec, err = ParseSort(raw); err != nil {
				http.Error(w, "sort inválido: "+err.Error(), http.StatusBadRequest)
				return
			}
			sorts = append(sorts, spec)
		}
		for _, raw := range q["select"] {
			var spec SelectSpec
			if spec, err = ParseSelect(raw); err != nil {
				http.Error(w, "select inválido: "+err.Error(), http.StatusBadRequest)
				return
			}
			selects = append(selects, spec)
		}

		// aplica aliases a partir do cfg.AliasMap
		filters = ApplyFilterAliases(filters, cfg.AliasMap)
		sorts = ApplySortAliases(sorts, cfg.AliasMap)
		selects = ApplySelectAliases(selects, cfg.AliasMap)

		// executa query genérica
		data, total, err := QueryDynamic[T](
			r.Context(), cfg.DB, cfg.Model,
			filters, sorts, selects,
			take, offset,
		)
		if err != nil {
			http.Error(w, "erro interno: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// monta resposta
		pageCount := int((total + int64(take) - 1) / int64(take))
		resp := struct {
			Data []map[string]interface{} `json:"data"`
			Meta PageMeta                 `json:"meta"`
		}{
			Data: data,
			Meta: PageMeta{
				Page:            page,
				Take:            take,
				ItemCount:       len(data),
				PageCount:       pageCount,
				HasPreviousPage: page > 0,
				HasNextPage:     page+1 < pageCount,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

// queryDynamic (sem alterações) permanece igual...
// queryDynamic executa a query e retorna rowsData + total.
func QueryDynamic[T any](
	ctx context.Context,
	db *gorm.DB,
	model T,
	filters []FilterSpec,
	sorts []SortSpec,
	selects []SelectSpec,
	limit, offset int,
) (rowsData []map[string]interface{}, total int64, err error) {
	tx := db.WithContext(ctx).Model(model)
	if tx, err = ApplyFilters(tx, filters); err != nil {
		return nil, 0, err
	}
	if err = tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	tx = ApplySorts(tx, sorts)

	if len(selects) > 0 {
		cols := make([]string, len(selects))
		for i, sel := range selects {
			cols[i] = strings.Join(sel.Path, ".")
		}
		rows, err2 := tx.Select(cols).Limit(limit).Offset(offset).Rows()
		if err2 != nil {
			return nil, total, err2
		}
		defer rows.Close()
		colsNames, _ := rows.Columns()
		for rows.Next() {
			vals := make([]interface{}, len(colsNames))
			ptrs := make([]interface{}, len(vals))
			for i := range vals {
				ptrs[i] = &vals[i]
			}
			if err2 = rows.Scan(ptrs...); err2 != nil {
				return nil, total, err2
			}
			m := make(map[string]interface{}, len(colsNames))
			for i, c := range colsNames {
				if b, ok := vals[i].([]byte); ok {
					m[c] = string(b)
				} else {
					m[c] = vals[i]
				}
			}
			rowsData = append(rowsData, m)
		}
		return rowsData, total, nil
	}

	var list []T
	if err = tx.Limit(limit).Offset(offset).Find(&list).Error; err != nil {
		return nil, total, err
	}
	for _, item := range list {
		raw, err2 := json.Marshal(item)
		if err2 != nil {
			return nil, total, err2
		}
		var m map[string]interface{}
		if err2 = json.Unmarshal(raw, &m); err2 != nil {
			return nil, total, err2
		}
		rowsData = append(rowsData, m)
	}
	return rowsData, total, nil
}
