package main

import (
	"database/sql"
	"fmt"
)

type QueryType string

const (
	QueryTypeMany QueryType = "many"
	QueryTypeOne  QueryType = "one"
)

type Where struct {
	Field string    `json:"field"`
	Value any       `json:"value"`
	Op    Operation `json:"op"`
}

type Operation string

const (
	OperationEqual Operation = "eq"
	OperationGTE   Operation = "gte"
	OperationLTE   Operation = "lte"
	OperationGT    Operation = "gt"
	OperationLT    Operation = "lt"
)

type Query struct {
	Type         QueryType `json:"type"`
	Where        []Where   `json:"where"` // joined by AND
	ReturnFields []string  `json:"return_fields"`
}

func (o Operation) DBString() string {
	switch o {
	case OperationEqual:
		return "="
	case OperationGTE:
		return ">="
	case OperationLTE:
		return "<="
	case OperationGT:
		return ">"
	case OperationLT:
		return "<"
	default:
		panic("unimplemented")
	}
}

func QueryOne(db *sql.DB, q Query) (map[string]any, error) {
	query, args := q.Build()

	var values []any = make([]any, len(q.ReturnFields))
	valuePointers := make([]any, len(q.ReturnFields))
	for i := range values {
		valuePointers[i] = &values[i]
	}
	err := db.QueryRow(query, args...).Scan(valuePointers...)
	if err != nil {
		return nil, fmt.Errorf("could not query database: %w", err)
	}

	result := make(map[string]any)

	for i, field := range q.ReturnFields {
		result[field] = values[i]
	}

	return result, nil
}

func QueryMany(db *sql.DB, q Query) ([]map[string]any, error) {
	query, args := q.Build()

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("could not query database: %w", err)
	}
	defer rows.Close()

	var results []map[string]any

	for rows.Next() {
		var values []any = make([]any, len(q.ReturnFields))
		valuePointers := make([]any, len(q.ReturnFields))
		for i := range values {
			valuePointers[i] = &values[i]
		}
		if err := rows.Scan(valuePointers...); err != nil {
			return nil, fmt.Errorf("could not scan row: %w", err)
		}

		result := make(map[string]any)

		for i, field := range q.ReturnFields {
			result[field] = values[i]
		}

		results = append(results, result)
	}

	return results, nil
}

func (q Query) Build() (string, []interface{}) {
	var query string
	var args []interface{}

	query = "SELECT "
	for i, field := range q.ReturnFields {
		if i > 0 {
			query += ", "
		}
		query += field
	}
	query += " FROM product_data WHERE "

	for i, where := range q.Where {
		if i > 0 {
			query += " AND "
		}
		query += WhereString(where.Field, where.Op, i+1)
		args = append(args, where.Value)
	}

	return query, args
}

// handle the special case of jsonb arrays containing strings
func WhereString(field string, op Operation, argi int) string {
	if op == OperationEqual {
		switch field {
		case "erzeugniss_nummern", "scip_nummern":
			return fmt.Sprintf("%s ? $%d", field, argi)
		}
	}

	return fmt.Sprintf("%s %s $%d", field, op.DBString(), argi)
}
