// Copyright 2013 The ql Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSES/QL-LICENSE file.

// Copyright 2015 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package tables

import (
	"github.com/chrislusf/gleamold/sql/model"
	"github.com/chrislusf/gleamold/sql/table"
)

var ()

type Table struct {
	Name    model.CIStr
	Columns []*table.Column

	meta *model.TableInfo
}

// Meta implements table.Table Meta interface.
func (t *Table) Meta() *model.TableInfo {
	return t.meta
}

func MockTableFromMeta(tableInfo *model.TableInfo) table.Table {
	return &Table{meta: tableInfo}
}

func init() {
	table.MockTableFromMeta = MockTableFromMeta
}
