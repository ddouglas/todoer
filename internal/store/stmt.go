package store

import (
	sq "github.com/Masterminds/squirrel"
)

func psql() sq.StatementBuilderType {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return psql
}
