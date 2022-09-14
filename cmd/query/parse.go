package query

import (
	"github.com/blastrain/vitess-sqlparser/sqlparser"
)

// UnknownType 0
// Select 1
// Update 2
// Insert 3
// Delete 4

// Unknown => 0
// Eq -> "=" 1
// Ne -> "!=" 2
// Gt -> ">" 3
// Lt -> "<" 4
// Gte -> ">=" 5
// Lte -> "<=" 6

func ParserSql(sql string) {
	query, err := sqlparser.Parse(sql)
	if err != nil {
		panic(err)
	}
	query.WalkSubtree(walk)
}

func walk(node sqlparser.SQLNode) (kontinue bool, err error) {
	if node != nil {
		sqlparser.String(node)
	}
	return true, nil
}
