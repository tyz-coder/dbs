package dba

import (
	"bytes"
	"strings"
	"strconv"
	"database/sql"
	"fmt"
	"errors"
)

type DeleteBuilder struct {
	prefixes     expressions
	options      expressions
	tables       []string
	using        string
	joins        []string
	wheres       whereExpressions
	orderBys     []string
	limit        uint64
	updateLimit  bool
	offset       uint64
	updateOffset bool
	suffixes     expressions
}

func (this *DeleteBuilder) Prefix(sql string, args ...interface{}) *DeleteBuilder {
	this.prefixes = append(this.prefixes, Expression(sql, args...))
	return this
}

func (this *DeleteBuilder) Options(options ...string) *DeleteBuilder {
	for _, c := range options {
		this.options = append(this.options, Expression(c))
	}
	return this
}

func (this *DeleteBuilder) Table(from ...string) *DeleteBuilder {
	this.tables = from
	return this
}

func (this *DeleteBuilder) USING(sql string) *DeleteBuilder {
	this.using = sql
	return this
}

func (this *DeleteBuilder) Join(join, table string) *DeleteBuilder {
	return this.join(join, table)
}

func (this *DeleteBuilder) RightJoin(table string) *DeleteBuilder {
	return this.join("RIGHT JOIN", table)
}

func (this *DeleteBuilder) LeftJoin(table string) *DeleteBuilder {
	return this.join("LEFT JOIN", table)
}

func (this *DeleteBuilder) join(join, table string) *DeleteBuilder {
	this.joins = append(this.joins, join, fmt.Sprintf("`%s`", table))
	return this
}

func (this *DeleteBuilder) Where(sql string, args ...interface{}) *DeleteBuilder {
	this.wheres = append(this.wheres, WhereExpression(sql, args...))
	return this
}

func (this *DeleteBuilder) OrderBy(sql ...string) *DeleteBuilder {
	this.orderBys = append(this.orderBys, sql...)
	return this
}

func (this *DeleteBuilder) Limit(limit uint64) *DeleteBuilder {
	this.limit = limit
	this.updateLimit = true
	return this
}

func (this *DeleteBuilder) Offset(offset uint64) *DeleteBuilder {
	this.offset = offset
	this.updateOffset = true
	return this
}

func (this *DeleteBuilder) Suffix(sql string, args ...interface{}) *DeleteBuilder {
	this.suffixes = append(this.suffixes, Expression(sql, args...))
	return this
}

func (this *DeleteBuilder) ToSQL() (sql string, args []interface{}, err error) {
	if len(this.tables) == 0 {
		return "", nil, errors.New("delete statements must specify a table")
	}

	var sqlBuffer = &bytes.Buffer{}
	if len(this.prefixes) > 0 {
		args, _ = this.prefixes.appendToSQL(sqlBuffer, " ", args)
		sqlBuffer.WriteString(" ")
	}

	sqlBuffer.WriteString("DELETE ")

	if len(this.options) > 0 {
		args, _ = this.options.appendToSQL(sqlBuffer, " ", args)
		sqlBuffer.WriteString(" ")
	}

	sqlBuffer.WriteString("FROM ")

	if len(this.tables) > 0 {
		sqlBuffer.WriteString(strings.Join(this.tables, ", "))
	}

	if len(this.using) > 0 {
		sqlBuffer.WriteString(" USING ")
		sqlBuffer.WriteString(this.using)
	}

	if len(this.joins) > 0 {
		sqlBuffer.WriteString(" ")
		sqlBuffer.WriteString(strings.Join(this.joins, " "))
	}

	if len(this.wheres) == 0 {
		return "", nil, errors.New("delete statements must have WHERE condition")
	}

	if len(this.wheres) > 0 {
		sqlBuffer.WriteString(" WHERE ")
		args, _ = this.wheres.appendToSQL(sqlBuffer, " ", args)
	}

	if len(this.orderBys) > 0 {
		sqlBuffer.WriteString(" ORDER BY ")
		sqlBuffer.WriteString(strings.Join(this.orderBys, ", "))
	}

	if this.updateLimit {
		sqlBuffer.WriteString(" LIMIT ")
		sqlBuffer.WriteString(strconv.FormatUint(this.limit, 10))
	}

	if this.updateOffset {
		sqlBuffer.WriteString(" OFFSET ")
		sqlBuffer.WriteString(strconv.FormatUint(this.offset, 10))
	}

	if len(this.suffixes) > 0 {
		sqlBuffer.WriteString(" ")
		args, _ = this.suffixes.appendToSQL(sqlBuffer, " ", args)
	}

	sql = sqlBuffer.String()

	return sql, args, err
}

func (this *DeleteBuilder) Exec(s StmtPrepare) (sql.Result, error) {
	sql, args, err := this.ToSQL()
	if err != nil {
		return nil, err
	}
	return Exec(s, sql, args...)
}

func NewDeleteBuilder() *DeleteBuilder {
	return &DeleteBuilder{}
}