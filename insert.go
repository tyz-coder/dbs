package dbs

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"strings"
	"context"
)

type InsertBuilder struct {
	prefixes statements
	options  statements
	columns  []string
	table    string
	values   [][]interface{}
	suffixes statements
}

func (this *InsertBuilder) Prefix(sql string, args ...interface{}) *InsertBuilder {
	this.prefixes = append(this.prefixes, NewStatement(sql, args...))
	return this
}

func (this *InsertBuilder) Options(options ...string) *InsertBuilder {
	for _, c := range options {
		this.options = append(this.options, NewStatement(c))
	}
	return this
}

func (this *InsertBuilder) Columns(columns ...string) *InsertBuilder {
	this.columns = append(this.columns, columns...)
	return this
}

func (this *InsertBuilder) Column(column string) *InsertBuilder {
	this.columns = append(this.columns, column)
	return this
}

func (this *InsertBuilder) Table(table string) *InsertBuilder {
	this.table = table
	return this
}

func (this *InsertBuilder) Values(values ...interface{}) *InsertBuilder {
	this.values = append(this.values, values)
	return this
}

func (this *InsertBuilder) Suffix(sql string, args ...interface{}) *InsertBuilder {
	this.suffixes = append(this.suffixes, NewStatement(sql, args...))
	return this
}

func (this *InsertBuilder) SET(column string, value interface{}) *InsertBuilder {
	this.columns = append(this.columns, column)
	if len(this.values) == 0 {
		this.values = append(this.values, make([]interface{}, 0, 0))
	}
	var vList = this.values[0]
	vList = append(vList, value)
	this.values[0] = vList
	return this
}

func (this *InsertBuilder) ToSQL() (string, []interface{}, error) {
	var sqlBuffer = &bytes.Buffer{}
	var args = newArgs()
	if err := this.AppendToSQL(sqlBuffer, args); err != nil {
		return "", nil, err
	}
	sql := sqlBuffer.String()
	log(sql, args.values)
	return sql, args.values, nil
}

func (this *InsertBuilder) AppendToSQL(w io.Writer, args *Args) error {
	if len(this.table) == 0 {
		return errors.New("insert statements must specify a table")
	}
	if len(this.values) == 0 {
		return errors.New("insert statements must have at least one set of values")
	}

	if len(this.prefixes) > 0 {
		if err := this.prefixes.AppendToSQL(w, " ", args); err != nil {
			return err
		}
		if _, err := io.WriteString(w, " "); err != nil {
			return err
		}
	}

	if _, err := io.WriteString(w, "INSERT "); err != nil {
		return err
	}

	if len(this.options) > 0 {
		if err := this.options.AppendToSQL(w, " ", args); err != nil {
			return err
		}
		if _, err := io.WriteString(w, " "); err != nil {
			return err
		}
	}

	if _, err := io.WriteString(w, fmt.Sprintf("INTO `%s` ", this.table)); err != nil {
		return err
	}

	if len(this.columns) > 0 {
		if _, err := io.WriteString(w, "(`"); err != nil {
			return err
		}
		if _, err := io.WriteString(w, strings.Join(this.columns, "`, `")); err != nil {
			return err
		}
		if _, err := io.WriteString(w, "`)"); err != nil {
			return err
		}
	}

	if _, err := io.WriteString(w, " VALUES "); err != nil {
		return err
	}

	var valuesPlaceholder = make([]string, len(this.values))
	for index, value := range this.values {
		var valuePlaceholder = make([]string, len(value))
		for i, v := range value {
			switch vt := v.(type) {
			case Statement:
				vSQL, vArgs, _ := vt.ToSQL()
				valuePlaceholder[i] = vSQL
				args.Append(vArgs...)
			default:
				valuePlaceholder[i] = "?"
				args.Append(v)
			}
		}
		valuesPlaceholder[index] = fmt.Sprintf("(%s)", strings.Join(valuePlaceholder, ", "))
	}
	if _, err := io.WriteString(w, strings.Join(valuesPlaceholder, ", ")); err != nil {
		return err
	}

	if len(this.suffixes) > 0 {
		if _, err := io.WriteString(w, " "); err != nil {
			return err
		}
		if err := this.suffixes.AppendToSQL(w, " ", args); err != nil {
			return err
		}
	}
	return nil
}

func (this *InsertBuilder) Exec(s SQLExecutor) (sql.Result, error) {
	sql, args, err := this.ToSQL()
	if err != nil {
		return nil, err
	}
	return s.Exec(sql, args...)
}

func (this *InsertBuilder) ExecContext(ctx context.Context, s SQLExecutor) (sql.Result, error) {
	sql, args, err := this.ToSQL()
	if err != nil {
		return nil, err
	}
	return s.ExecContext(ctx, sql, args...)
}

func NewInsertBuilder() *InsertBuilder {
	return &InsertBuilder{}
}

func Insert(s SQLExecutor, table string, data map[string]interface{}) (sql.Result, error) {
	var in = NewInsertBuilder()
	in.Table(table)

	var values []interface{}
	for k, v := range data {
		in.Column(k)
		values = append(values, v)
	}
	in.Values(values...)
	return in.Exec(s)
}