package db

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"gorm.io/gorm"
)

const (
	// defaultTableName  if tableName == "", the Adapter will use this default table name.
	defaultTableName = "casbin_rules"

	// maxParamLength  .
	maxParamLength = 7

	sqlPlaceHolder = "?"
	sqlCreateTable = `
CREATE TABLE %[1]s(
    p_type VARCHAR(32),
    v0     VARCHAR(255),
    v1     VARCHAR(255),
    v2     VARCHAR(255),
    v3     VARCHAR(255),
    v4     VARCHAR(255),
    v5     VARCHAR(255)
);
CREATE INDEX idx_%[1]s ON %[1]s (p_type,v0,v1);`
	sqlTruncateTable = "TRUNCATE TABLE %s"
	sqlIsTableExist  = `SHOW TABLES LIKE '%s'`
	sqlInsertRow     = "INSERT INTO %s (p_type,v0,v1,v2,v3,v4,v5) VALUES (?,?,?,?,?,?,?)"
	sqlUpdateRow     = "UPDATE %s SET p_type=?,v0=?,v1=?,v2=?,v3=?,v4=?,v5=? WHERE p_type=? AND v0=? AND v1=? AND v2=? AND v3=? AND v4=? AND v5=?"
	sqlDeleteAll     = "DELETE FROM %s"
	sqlDeleteRow     = "DELETE FROM %s WHERE p_type=? AND v0=? AND v1=? AND v2=? AND v3=? AND v4=? AND v5=?"
	sqlDeleteByArgs  = "DELETE FROM %s WHERE p_type=?"
	sqlSelectAll     = "SELECT p_type,v0,v1,v2,v3,v4,v5 FROM %s"
	sqlSelectWhere   = "SELECT p_type,v0,v1,v2,v3,v4,v5 FROM %s WHERE "
	sqlCreateTableMysql = `
CREATE TABLE IF NOT EXISTS %[1]s(
    p_type VARCHAR(32)  DEFAULT '' NOT NULL,
    v0     VARCHAR(255) DEFAULT '' NOT NULL,
    v1     VARCHAR(255) DEFAULT '' NOT NULL,
    v2     VARCHAR(255) DEFAULT '' NOT NULL,
    v3     VARCHAR(255) DEFAULT '' NOT NULL,
    v4     VARCHAR(255) DEFAULT '' NOT NULL,
    v5     VARCHAR(255) DEFAULT '' NOT NULL,
    INDEX idx_%[1]s (p_type,v0,v1)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;`
)


// CasbinRule  defines the casbin rule model.
// It used for save or load policy lines from connected database.
type CasbinRule struct {
	PType string
	V0    string
	V1    string
	V2    string
	V3    string
	V4    string
	V5    string
}

// Adapter  defines the database adapter for Casbin.
// It can load policy lines from connected database or save policy lines.
type Adapter struct {
	db         *gorm.DB
	ctx        context.Context
	driverName string
	tableName  string

	isFiltered bool

	sqlPlaceHolder   string
	sqlCreateTable   string
	sqlTruncateTable string
	sqlIsTableExist  string
	sqlInsertRow     string
	sqlUpdateRow     string
	sqlDeleteAll     string
	sqlDeleteRow     string
	sqlDeleteByArgs  string
	sqlSelectAll     string
	sqlSelectWhere   string
}

// Filter  defines the filtering rules for a FilteredAdapter's policy.
// Empty values are ignored, but all others must match the filter.
type Filter struct {
	PType []string
	V0    []string
	V1    []string
	V2    []string
	V3    []string
	V4    []string
	V5    []string
}

// NewAdapter  the constructor for Adapter.
// db should connected to database and controlled by user.
// The Adapter will automatically create a table named "casbin_rule".
func NewAdapter(db *gorm.DB) (*Adapter, error) {
	return NewAdapterContext(context.Background(), db)
}

// NewAdapterContext  the constructor for Adapter.
// db should connected to database and controlled by user.
// The Adapter will automatically create a table named "casbin_rule".
func NewAdapterContext(ctx context.Context, db *gorm.DB) (*Adapter, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	adapter := Adapter{
		db:         db,
		ctx:        ctx,
		driverName: "mysql",
		tableName:  defaultTableName,
	}

	// generate different databases sql
	adapter.genSQL()

	if !adapter.isTableExist() {
		if err := adapter.createTable(); err != nil {
			return nil, err
		}
	}

	return &adapter, nil
}

// genSQL  generate sql based on db driver name.
func (p *Adapter) genSQL() {
	p.sqlPlaceHolder = sqlPlaceHolder
	p.sqlCreateTable = fmt.Sprintf(sqlCreateTable, p.tableName)
	p.sqlTruncateTable = fmt.Sprintf(sqlTruncateTable, p.tableName)
	p.sqlIsTableExist = fmt.Sprintf(sqlIsTableExist, p.tableName)
	p.sqlInsertRow = fmt.Sprintf(sqlInsertRow, p.tableName)
	p.sqlUpdateRow = fmt.Sprintf(sqlUpdateRow, p.tableName)
	p.sqlDeleteAll = fmt.Sprintf(sqlDeleteAll, p.tableName)
	p.sqlDeleteRow = fmt.Sprintf(sqlDeleteRow, p.tableName)
	p.sqlDeleteByArgs = fmt.Sprintf(sqlDeleteByArgs, p.tableName)
	p.sqlSelectAll = fmt.Sprintf(sqlSelectAll, p.tableName)
	p.sqlSelectWhere = fmt.Sprintf(sqlSelectWhere, p.tableName)
	p.sqlCreateTable = fmt.Sprintf(sqlCreateTableMysql, p.tableName)
}

func (p *Adapter) sqlRebind(query string) string {
	if p.sqlPlaceHolder == sqlPlaceHolder {
		return query
	}

	var idx, num int

	result := make([]byte, 0, len(query)+10)

	for {
		idx = strings.Index(query, sqlPlaceHolder)
		if idx == -1 {
			break
		}

		num++

		result = append(result, query[:idx]...)
		result = append(result, p.sqlPlaceHolder...)
		result = strconv.AppendInt(result, int64(num), 10)

		query = query[idx+1:]
	}

	return string(append(result, query...))
}

// createTable  create a not exists table.
func (p *Adapter) createTable() error {
	p.db.WithContext(p.ctx).Exec(p.sqlCreateTable)

	return nil
}

// truncateTable  clear the table.
func (p *Adapter) truncateTable() error {
	err := p.db.WithContext(p.ctx).Raw(p.sqlTruncateTable, defaultTableName).Error

	return err
}

// deleteAll  clear the table.
func (p *Adapter) deleteAll() error {
	err := p.db.WithContext(p.ctx).Raw(p.sqlDeleteAll, defaultTableName).Error

	return err
}

// isTableExist  check the table exists.
func (p *Adapter) isTableExist() bool {
	var val string
	p.db.Raw(p.sqlIsTableExist).Scan(&val)

	return val != ""
}

// deleteRows  delete eligible data.
func (p *Adapter) deleteRows(query string, args ...interface{}) error {
	query = p.sqlRebind(query)

	err := p.db.WithContext(p.ctx).Exec(query, args...).Error

	return err
}

// truncateAndInsertRows  clear table and insert new rows.
func (p *Adapter) truncateAndInsertRows(rules [][]interface{}) error {
	if err := p.truncateTable(); err != nil {
		return err
	}
	return p.execTxSqlRows(p.sqlInsertRow, rules)
}

// deleteAllAndInsertRows  clear table and insert new rows.
func (p *Adapter) deleteAllAndInsertRows(rules [][]interface{}) error {
	if err := p.deleteAll(); err != nil {
		return err
	}
	return p.execTxSqlRows(p.sqlInsertRow, rules)
}

// execTxSqlRows  exec sql rows.
func (p *Adapter) execTxSqlRows(query string, rules [][]interface{}) (err error) {
	tx := p.db.WithContext(p.ctx).Begin()
	if tx.Error != nil {
		return
	}

	var action string

	stmt, err := tx.ConnPool.PrepareContext(p.ctx, query)
	if err != nil {
		action = "prepare context"
		goto ROLLBACK
	}

	for _, rule := range rules {
		if _, err = stmt.ExecContext(p.ctx, rule...); err != nil {
			action = "stmt exec context"
			goto ROLLBACK
		}
	}

	if err = stmt.Close(); err != nil {
		action = "stmt close"
		goto ROLLBACK
	}

	if db := tx.Commit(); db.Error != nil {
		action = "commit"
		goto ROLLBACK
	}

	return

ROLLBACK:

	if err1 := tx.Rollback(); err1 != nil {
		err = fmt.Errorf("%s err: %v, rollback err: %v", action, err, err1)
	}

	return
}

// queryFunc  define func for query
var queryFunc = func(ctx context.Context, db *sql.DB, query string, args ...interface{}) ([]*CasbinRule, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	lines := make([]*CasbinRule, 0, 64)
	for rows.Next() {
		var rule CasbinRule

		err = rows.Scan(&rule.PType, &rule.V0, &rule.V1, &rule.V2, &rule.V3, &rule.V4, &rule.V5)
		if err != nil {
			return nil, err
		}

		lines = append(lines, &rule)
	}

	return lines, nil
}

// selectRows  select eligible data by args from the table.
func (p *Adapter) selectRows(query string, args ...interface{}) ([]*CasbinRule, error) {
	db, _ := p.db.DB()
	if len(args) == 0 {
		return queryFunc(p.ctx, db, query)
	}

	query = p.sqlRebind(query)

	return queryFunc(p.ctx, db, query, args...)
}

// selectWhereIn  select eligible data by filter from the table.
func (p *Adapter) selectWhereIn(filter *Filter) (lines []*CasbinRule, err error) {
	var (
		sqlBuf bytes.Buffer
		buf    bytes.Buffer
	)

	sqlBuf.Grow(64)
	sqlBuf.WriteString(p.sqlSelectWhere)

	args := make([]string, 0, 4)

	for _, col := range [maxParamLength]struct {
		name string
		arg  []string
	}{
		{"p_type", filter.PType},
		{"v0", filter.V0},
		{"v1", filter.V1},
		{"v2", filter.V2},
		{"v3", filter.V3},
		{"v4", filter.V4},
		{"v5", filter.V5},
	} {
		l := len(col.arg)
		if l == 0 {
			continue
		}

		switch sqlBuf.Bytes()[sqlBuf.Len()-1] {
		case '?', ')':
			sqlBuf.WriteString(" AND ")
		}

		sqlBuf.WriteString(col.name)

		if l == 1 {
			sqlBuf.WriteString("=?")
			args = append(args, col.arg[0])
		} else {
			buf.Grow(l * 2)
			for i := 0; i < l; i++ {
				buf.WriteString("?,")
			}
			buf.Truncate(buf.Len() - 1)

			sqlBuf.WriteString(" IN (")
			sqlBuf.Write(buf.Bytes())
			sqlBuf.WriteByte(')')

			args = append(args, col.arg...)

			buf.Reset()
		}
	}

	params := make([]interface{}, len(args))
	for idx := range args {
		params[idx] = args[idx]
	}

	return p.selectRows(sqlBuf.String(), params...)
}

// LoadPolicy  load all policy rules from the storage.
func (p *Adapter) LoadPolicy(model model.Model) error {
	lines, err := p.selectRows(p.sqlSelectAll)
	if err != nil {
		return err
	}

	for _, line := range lines {
		p.loadPolicyLine(line, model)
	}

	return nil
}

// SavePolicy  save policy rules to the storage.
func (p *Adapter) SavePolicy(model model.Model) error {
	args := make([][]interface{}, 0, 32)

	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			arg := p.genArgs(ptype, rule)
			args = append(args, arg)
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			arg := p.genArgs(ptype, rule)
			args = append(args, arg)
		}
	}

	// return p.truncateAndInsertRows(args)
	return p.deleteAllAndInsertRows(args)
}

// AddPolicy  add one policy rule to the storage.
func (p *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	args := p.genArgs(ptype, rule)

	err := p.db.WithContext(p.ctx).Exec(p.sqlInsertRow, args...).Error

	return err
}

// AddPolicies  add multiple policy rules to the storage.
func (p *Adapter) AddPolicies(sec string, ptype string, rules [][]string) error {
	args := make([][]interface{}, 0, 8)

	for _, rule := range rules {
		arg := p.genArgs(ptype, rule)
		args = append(args, arg)
	}

	return p.execTxSqlRows(p.sqlInsertRow, args)
}

// RemovePolicy  remove policy rules from the storage.
func (p *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	var sqlBuf bytes.Buffer

	sqlBuf.Grow(64)
	sqlBuf.WriteString(p.sqlDeleteByArgs)

	args := make([]interface{}, 0, 4)
	args = append(args, ptype)

	for idx, arg := range rule {
		if arg != "" {
			sqlBuf.WriteString(" AND v")
			sqlBuf.WriteString(strconv.Itoa(idx))
			sqlBuf.WriteString("=?")

			args = append(args, arg)
		}
	}

	return p.deleteRows(sqlBuf.String(), args...)
}

// RemoveFilteredPolicy  remove policy rules that match the filter from the storage.
func (p *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	var sqlBuf bytes.Buffer

	sqlBuf.Grow(64)
	sqlBuf.WriteString(p.sqlDeleteByArgs)

	args := make([]interface{}, 0, 4)
	args = append(args, ptype)

	var value string

	l := fieldIndex + len(fieldValues)

	for idx := 0; idx < 6; idx++ {
		if fieldIndex <= idx && idx < l {
			value = fieldValues[idx-fieldIndex]

			if value != "" {
				sqlBuf.WriteString(" AND v")
				sqlBuf.WriteString(strconv.Itoa(idx))
				sqlBuf.WriteString("=?")

				args = append(args, value)
			}
		}
	}

	return p.deleteRows(sqlBuf.String(), args...)
}

func (p *Adapter) RemovePolicies(sec string, ptype string, rules [][]string) (err error) {
	args := make([][]interface{}, 0, 8)

	for _, rule := range rules {
		arg := p.genArgs(ptype, rule)
		args = append(args, arg)
	}

	return p.execTxSqlRows(p.sqlDeleteRow, args)
}

// LoadFilteredPolicy  load policy rules that match the filter.
// filterPtr must be a pointer.
func (p *Adapter) LoadFilteredPolicy(model model.Model, filterPtr interface{}) error {
	if filterPtr == nil {
		return p.LoadPolicy(model)
	}

	filter, ok := filterPtr.(*Filter)
	if !ok {
		return errors.New("invalid filter type")
	}

	lines, err := p.selectWhereIn(filter)
	if err != nil {
		return err
	}

	for _, line := range lines {
		p.loadPolicyLine(line, model)
	}

	p.isFiltered = true

	return nil
}

// IsFiltered  returns true if the loaded policy rules has been filtered.
func (p *Adapter) IsFiltered() bool {
	return p.isFiltered
}

// UpdatePolicy update a policy rule from storage.
// This is part of the Auto-Save feature.
func (p *Adapter) UpdatePolicy(sec, ptype string, oldRule, newPolicy []string) error {
	oldArg := p.genArgs(ptype, oldRule)
	newArg := p.genArgs(ptype, newPolicy)

	err := p.db.WithContext(p.ctx).Exec(p.sqlUpdateRow, append(newArg, oldArg...)...).Error

	return err
}

// UpdatePolicies updates policy rules to storage.
func (p *Adapter) UpdatePolicies(sec, ptype string, oldRules, newRules [][]string) (err error) {
	if len(oldRules) != len(newRules) {
		return errors.New("old rules size not equal to new rules size")
	}

	args := make([][]interface{}, 0, 16)

	for idx := range oldRules {
		oldArg := p.genArgs(ptype, oldRules[idx])
		newArg := p.genArgs(ptype, newRules[idx])
		args = append(args, append(newArg, oldArg...))
	}

	return p.execTxSqlRows(p.sqlUpdateRow, args)
}

// UpdateFilteredPolicies deletes old rules and adds new rules.
func (p *Adapter) UpdateFilteredPolicies(sec, ptype string, newPolicies [][]string, fieldIndex int, fieldValues ...string) (oldPolicies [][]string, err error) {
	var value string

	var whereBuf bytes.Buffer
	whereBuf.Grow(32)

	l := fieldIndex + len(fieldValues)

	whereArgs := make([]interface{}, 0, 4)
	whereArgs = append(whereArgs, ptype)

	for idx := 0; idx < 6; idx++ {
		if fieldIndex <= idx && idx < l {
			value = fieldValues[idx-fieldIndex]

			if value != "" {
				whereBuf.WriteString(" AND v")
				whereBuf.WriteString(strconv.Itoa(idx))
				whereBuf.WriteString("=?")

				whereArgs = append(whereArgs, value)
			}
		}
	}

	var selectBuf bytes.Buffer
	selectBuf.Grow(64)
	selectBuf.WriteString(p.sqlSelectWhere)
	selectBuf.WriteString("p_type=?")
	selectBuf.Write(whereBuf.Bytes())

	var oldRows []*CasbinRule
	value = p.sqlRebind(selectBuf.String())
	oldRows, err = p.selectRows(value, whereArgs...)
	if err != nil {
		return
	}

	var deleteBuf bytes.Buffer
	deleteBuf.Grow(64)
	deleteBuf.WriteString(p.sqlDeleteByArgs)
	deleteBuf.Write(whereBuf.Bytes())

	var tx *sql.Tx
	tx1 := p.db.WithContext(p.ctx).Begin()
	if tx1.Error != nil {
		return
	}

	var (
		stmt   *sql.Stmt
		action string
	)
	value = p.sqlRebind(deleteBuf.String())
	if _, err = tx.ExecContext(p.ctx, value, whereArgs...); err != nil {
		action = "delete old policies"
		goto ROLLBACK
	}

	stmt, err = tx.PrepareContext(p.ctx, p.sqlInsertRow)
	if err != nil {
		action = "prepare context"
		goto ROLLBACK
	}

	for _, policy := range newPolicies {
		arg := p.genArgs(ptype, policy)
		if _, err = stmt.ExecContext(p.ctx, arg...); err != nil {
			action = "stmt exec context"
			goto ROLLBACK
		}
	}

	if err = stmt.Close(); err != nil {
		action = "stmt close"
		goto ROLLBACK
	}

	if err = tx.Commit(); err != nil {
		action = "commit"
		goto ROLLBACK
	}

	oldPolicies = make([][]string, 0, len(oldRows))
	for _, rule := range oldRows {
		oldPolicies = append(oldPolicies, []string{rule.PType, rule.V0, rule.V1, rule.V2, rule.V3, rule.V4, rule.V5})
	}

	return

ROLLBACK:

	if err1 := tx.Rollback(); err1 != nil {
		err = fmt.Errorf("%s err: %v, rollback err: %v", action, err, err1)
	}

	return
}

// loadPolicyLine  load a policy line to model.
func (Adapter) loadPolicyLine(line *CasbinRule, model model.Model) {
	if line == nil {
		return
	}

	var lineBuf bytes.Buffer

	lineBuf.Grow(64)
	lineBuf.WriteString(line.PType)

	args := [6]string{line.V0, line.V1, line.V2, line.V3, line.V4, line.V5}
	for _, arg := range args {
		if arg != "" {
			lineBuf.WriteByte(',')
			lineBuf.WriteString(arg)
		}
	}

	persist.LoadPolicyLine(lineBuf.String(), model)
}

// genArgs  generate args from ptype and rule.
func (Adapter) genArgs(ptype string, rule []string) []interface{} {
	args := make([]interface{}, maxParamLength)
	args[0] = ptype

	for idx := range rule {
		args[idx+1] = rule[idx]
	}

	for idx := len(rule) + 1; idx < maxParamLength; idx++ {
		args[idx] = ""
	}

	return args
}