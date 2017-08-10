package executor

import (
	"github.com/chrislusf/gleamold/flow"
	"github.com/chrislusf/gleamold/sql/model"
)

type TableColumn struct {
	ColumnName string
	ColumnType byte
}

type TableSource struct {
	Dataset   *flow.Dataset
	TableInfo *model.TableInfo
}

var (
	Tables = make(map[string]*TableSource)
)
