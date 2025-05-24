package dynamicquery

import "errors"

// Defina aqui erros compartilhados, se desejar:
var ErrInvalidBetween = errors.New("dynamicquery: between precisa de dois valores separados por vírgula")
