package database

import "database/sql"

func toNullString(s string) sql.Null[string] {
	return sql.Null[string]{
		V:     s,
		Valid: s != "",
	}
}
