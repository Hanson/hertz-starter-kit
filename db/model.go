package db

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"hertz-starter-kit/utils/log"
	"time"
)

type Model struct {
	db *gorm.DB

	withTrash bool

	ctx context.Context
}

// NewModel 当有 where 条件并且入参为 struct 时，会自动调用 db.Model
func NewModel(m interface{}, ctx context.Context) *Model {
	return &Model{
		db:        Db.Model(m).WithContext(ctx).Session(&gorm.Session{NewDB: false}),
		withTrash: false,
		ctx:       ctx,
	}
}

func (m *Model) WithTrash() *Model {
	m.withTrash = true
	return m
}

func (m *Model) Where(query interface{}, args ...interface{}) *Model {
	if len(args) > 0 {
		m.db = m.db.Where(query, args[0])
	} else {
		m.db = m.db.Where(query)
	}

	return m
}

func (m *Model) SoftDelete() *gorm.DB {
	return m.db.Update("deleted_at", time.Now().Unix())
}

func (m *Model) WhereIn(field string, args ...interface{}) *Model {
	if len(args) == 0 {
		m.db = m.db.Where("1 = 0")
	}
	m.db = m.db.Where(fmt.Sprintf("%s IN ?", field), args[0])
	return m
}

func (m *Model) addScopeBeforeExec() *Model {
	if !m.withTrash {
		m.db = m.db.Where("deleted_at = 0")
	}
	return m
}

func (m *Model) First(res interface{}) error {
	m.addScopeBeforeExec()

	return m.db.First(res).Error
}

func (m *Model) Find(res interface{}) error {
	m.addScopeBeforeExec()

	return m.db.Find(res).Error
}

type ExecResult struct {
	RowsAffected int64
	IsCreate     bool
}

func (m *Model) Update(arg ...interface{}) (*ExecResult, error) {
	execResult := &ExecResult{}

	m.addScopeBeforeExec()

	if len(arg) < 1 || len(arg) > 2 {
		return nil, fmt.Errorf("invalid argument count")
	}

	var res *gorm.DB
	if len(arg) == 1 {
		res = m.db.Updates(arg[0])
	} else {
		res = m.db.Update(arg[0].(string), arg[1])
	}

	execResult.RowsAffected = res.RowsAffected

	return execResult, res.Error
}

func (m *Model) Create(res interface{}) error {
	return m.db.Create(res).Error
}

func (m *Model) FirstOrCreate(query interface{}, attrs interface{}, res interface{}) *gorm.DB {
	m.addScopeBeforeExec()

	return m.db.Where(query).Attrs(attrs).FirstOrCreate(res)
}

//	func (m *Model) UpdateOrCreate(query interface{}, assign interface{}, res interface{}) *gorm.DB {
//		m.addScopeBeforeExec()
//
//		return m.db.Where(query).Assign(assign).FirstOrCreate(res)
//	}
func (m *Model) UpdateOrCreate(attributes map[string]interface{}, values map[string]interface{}, res interface{}) (*ExecResult, error) {
	m.addScopeBeforeExec()

	execResult := &ExecResult{}

	queryTx := m.db.Session(&gorm.Session{}).Limit(1).Order(clause.OrderByColumn{
		Column: clause.Column{Table: clause.CurrentTable, Name: clause.PrimaryKey},
	})

	result := queryTx.Where(attributes).Find(res)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		execResult.IsCreate = true
		all := map[string]interface{}{}
		for k, v := range values {
			all[k] = v
		}
		for k, v := range attributes {
			all[k] = v
		}
		all["created_at"] = time.Now().Unix()
		all["updated_at"] = all["created_at"]
		err := m.db.Create(all).Error
		if err != nil {
			log.Errorf(m.ctx, "err: %+v", err)
			return nil, err
		}

		return execResult, nil
	} else {
		result := queryTx.Where(attributes).Updates(values)
		if result.Error != nil {
			return nil, result.Error
		}

		execResult.RowsAffected = result.RowsAffected
		return execResult, nil
	}
}

func (m *Model) Increment(column string, cnt int32) (*ExecResult, error) {
	execResult := &ExecResult{}
	m.addScopeBeforeExec()

	res := m.db.UpdateColumn(column, gorm.Expr(column+" + ?", cnt))
	execResult.RowsAffected = res.RowsAffected
	return execResult, res.Error
}
