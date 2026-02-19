// Package query provides a lightweight utility for parsing search, sort,
// and pagination query parameters from HTTP requests and applying them
// as safe, parameterized SQL clauses.
package query

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Filter represents a single WHERE clause with a parameterized value.
type Filter struct {
	Column string // SQL column, e.g. "is_active"
	Op     string // SQL operator, e.g. "="
	Value  any    // parameterized value
}

// Params holds parsed and validated query parameters.
type Params struct {
	Search  string   // free-text search term
	Sort    string   // validated column name
	Order   string   // "ASC" or "DESC"
	Limit   int
	Offset  int
	Filters []Filter // exact-match filters
}

// Config defines constraints for a specific list endpoint.
type Config struct {
	AllowedSorts  map[string]string // map[apiName]sqlColumn, e.g. {"name": "name", "created_at": "created_at"}
	DefaultSort   string            // SQL column to sort by when none requested
	DefaultOrder  string            // "ASC" or "DESC"
	SearchColumns []string          // SQL columns to ILIKE search across
	MaxLimit      int               // hard cap, default 100
	DefaultLimit  int               // default if not provided, default 50
}

// Parse extracts and validates query parameters from a request.
func Parse(r *http.Request, cfg Config) Params {
	if cfg.MaxLimit <= 0 {
		cfg.MaxLimit = 100
	}
	if cfg.DefaultLimit <= 0 {
		cfg.DefaultLimit = 50
	}
	if cfg.DefaultOrder == "" {
		cfg.DefaultOrder = "DESC"
	}

	q := r.URL.Query()

	// Pagination
	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit <= 0 || limit > cfg.MaxLimit {
		limit = cfg.DefaultLimit
	}
	offset, _ := strconv.Atoi(q.Get("offset"))
	if offset < 0 {
		offset = 0
	}

	// Sort
	sortCol := cfg.DefaultSort
	if s := q.Get("sort"); s != "" {
		if col, ok := cfg.AllowedSorts[s]; ok {
			sortCol = col
		}
	}

	order := cfg.DefaultOrder
	if o := strings.ToUpper(q.Get("order")); o == "ASC" || o == "DESC" {
		order = o
	}

	// Search
	search := strings.TrimSpace(q.Get("search"))

	return Params{
		Search: search,
		Sort:   sortCol,
		Order:  order,
		Limit:  limit,
		Offset: offset,
	}
}

// Apply builds the WHERE + ORDER BY + LIMIT/OFFSET clauses and returns
// the complete SQL string along with positional arguments.
//
// baseSQL must be a complete SELECT without WHERE/ORDER/LIMIT, e.g.
//
//	"SELECT id, email, full_name, is_active, created_at FROM users"
//
// searchCols are the SQL columns to match the search term against via ILIKE.
func Apply(baseSQL string, p Params, searchCols []string) (string, []any) {
	var (
		clauses []string
		args    []any
		idx     = 1 // positional parameter counter
	)

	// Exact-match filters: is_active = $1, etc.
	for _, f := range p.Filters {
		clauses = append(clauses, fmt.Sprintf("%s %s $%d", f.Column, f.Op, idx))
		args = append(args, f.Value)
		idx++
	}

	// Search clause: (col1 ILIKE $N OR col2 ILIKE $N ...)
	if p.Search != "" && len(searchCols) > 0 {
		like := "%" + sanitizeLike(p.Search) + "%"
		args = append(args, like)

		var ors []string
		for _, col := range searchCols {
			ors = append(ors, fmt.Sprintf("%s ILIKE $%d", col, idx))
		}
		clauses = append(clauses, "("+strings.Join(ors, " OR ")+")")
		idx++
	}

	// Build SQL
	sql := baseSQL
	if len(clauses) > 0 {
		sql += " WHERE " + strings.Join(clauses, " AND ")
	}

	if p.Sort != "" {
		sql += fmt.Sprintf(" ORDER BY %s %s", p.Sort, p.Order)
	}

	sql += fmt.Sprintf(" LIMIT $%d OFFSET $%d", idx, idx+1)
	args = append(args, p.Limit, p.Offset)

	return sql, args
}

// ApplyCount builds a COUNT query with the same WHERE clause as Apply.
func ApplyCount(baseCountSQL string, p Params, searchCols []string) (string, []any) {
	var (
		clauses []string
		args    []any
		idx     = 1
	)

	// Exact-match filters
	for _, f := range p.Filters {
		clauses = append(clauses, fmt.Sprintf("%s %s $%d", f.Column, f.Op, idx))
		args = append(args, f.Value)
		idx++
	}

	if p.Search != "" && len(searchCols) > 0 {
		like := "%" + sanitizeLike(p.Search) + "%"
		args = append(args, like)

		var ors []string
		for _, col := range searchCols {
			ors = append(ors, fmt.Sprintf("%s ILIKE $%d", col, idx))
		}
		clauses = append(clauses, "("+strings.Join(ors, " OR ")+")")
	}

	sql := baseCountSQL
	if len(clauses) > 0 {
		sql += " WHERE " + strings.Join(clauses, " AND ")
	}

	return sql, args
}

// sanitizeLike escapes special LIKE characters (%, _, \).
func sanitizeLike(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}
