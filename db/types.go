package db

import "database/sql"

func NullStringFrom(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}

func MapStrings[T any](list []T, field func(T) string) []string {
	result := make([]string, len(list))
	for i, item := range list {
		result[i] = field(item)
	}
	return result
}
