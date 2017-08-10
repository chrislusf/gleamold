package executor

import (
	"github.com/chrislusf/gleamold/flow"
	"github.com/chrislusf/gleamold/sql/expression"
)

type Executor interface {
	Exec() *flow.Dataset
	Schema() expression.Schema
}
