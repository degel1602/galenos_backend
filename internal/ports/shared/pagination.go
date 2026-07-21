// Package shared contiene tipos reutilizables entre puertos de entrada y
// salida (paginación) que no pertenecen a ningún agregado de dominio en
// particular.
package shared

// PageRequest normaliza los parámetros de paginación recibidos desde un
// adaptador de entrada.
type PageRequest struct {
	Page     int
	PageSize int
}

const maxPageSize = 100

// NewPageRequest aplica valores por defecto y límites razonables.
func NewPageRequest(page, pageSize int) PageRequest {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > maxPageSize {
		pageSize = 20
	}
	return PageRequest{Page: page, PageSize: pageSize}
}

// Offset calcula el desplazamiento SQL (OFFSET) para esta página.
func (p PageRequest) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// PageResponse es el sobre estándar de una respuesta paginada hacia Angular.
type PageResponse[T any] struct {
	Items      []T `json:"items"`
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

// NewPageResponse construye la respuesta paginada a partir de los items de
// la página actual y el total de registros que reportó el repositorio.
func NewPageResponse[T any](items []T, req PageRequest, totalItems int) PageResponse[T] {
	totalPages := 0
	if req.PageSize > 0 {
		totalPages = (totalItems + req.PageSize - 1) / req.PageSize
	}
	return PageResponse[T]{
		Items:      items,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}
