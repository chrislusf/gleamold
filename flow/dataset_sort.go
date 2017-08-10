package flow

import (
	"github.com/chrislusf/gleam/instruction"
	"github.com/ugorji/go/codec"
)

var (
	msgpackHandler codec.MsgpackHandle
)

type pair struct {
	keys []interface{}
	data []byte
}

// Distinct sort on specific fields and pick the unique ones.
// Required Memory: about same size as each partition.
// example usage: Distinct(Field(1,2)) means
// distinct on field 1 and 2.
// TODO: optimize for low cardinality case.
func (d *Dataset) Distinct(sortOptions ...*SortOption) *Dataset {
	sortOption := concat(sortOptions)

	ret := d.LocalSort(sortOption).LocalDistinct(sortOption)
	if len(d.Shards) > 1 {
		ret = ret.MergeSortedTo(1, sortOption).LocalDistinct(sortOption)
	}
	return ret
}

// Sort sort on specific fields, default to the first field.
// Required Memory: about same size as each partition.
// example usage: Sort(Field(1,2)) means
// sorting on field 1 and 2.
func (d *Dataset) Sort(sortOptions ...*SortOption) *Dataset {
	sortOption := concat(sortOptions)

	ret := d.LocalSort(sortOption)
	ret = ret.TreeMergeSortedTo(1, 10, sortOption)
	return ret
}

// Top streams through total n items, picking reverse ordered k items with O(n*log(k)) complexity.
// Required Memory: about same size as n items in memory
func (d *Dataset) Top(k int, sortOptions ...*SortOption) *Dataset {
	sortOption := concat(sortOptions)

	ret := d.LocalTop(k, sortOption)
	if len(d.Shards) > 1 {
		ret = ret.MergeSortedTo(1, sortOption).LocalLimit(k, 0)
	}
	return ret
}

func (d *Dataset) LocalDistinct(sortOptions ...*SortOption) *Dataset {
	sortOption := concat(sortOptions)

	ret, step := add1ShardTo1Step(d)
	ret.IsLocalSorted = sortOption.orderByList
	ret.IsPartitionedBy = d.IsPartitionedBy
	step.SetInstruction(instruction.NewLocalDistinct(sortOption.orderByList))
	return ret
}

func (d *Dataset) LocalSort(sortOptions ...*SortOption) *Dataset {
	sortOption := concat(sortOptions)

	if isOrderByEquals(d.IsLocalSorted, sortOption.orderByList) {
		return d
	}

	ret, step := add1ShardTo1Step(d)
	ret.IsLocalSorted = sortOption.orderByList
	ret.IsPartitionedBy = d.IsPartitionedBy
	step.SetInstruction(instruction.NewLocalSort(sortOption.orderByList, int(d.GetPartitionSize())*3))
	return ret
}

func (d *Dataset) LocalTop(n int, sortOptions ...*SortOption) *Dataset {
	sortOption := concat(sortOptions)

	if isOrderByExactReverse(d.IsLocalSorted, sortOption.orderByList) {
		return d.LocalLimit(n, 0)
	}

	ret, step := add1ShardTo1Step(d)
	ret.IsLocalSorted = sortOption.orderByList
	ret.IsPartitionedBy = d.IsPartitionedBy
	step.SetInstruction(instruction.NewLocalTop(n, sortOption.orderByList))
	return ret
}

func isOrderByEquals(a []instruction.OrderBy, b []instruction.OrderBy) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v.Index != b[i].Index || v.Order != b[i].Order {
			return false
		}
	}
	return true
}

func isOrderByExactReverse(a []instruction.OrderBy, b []instruction.OrderBy) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v.Index != b[i].Index || v.Order == b[i].Order {
			return false
		}
	}
	return true
}

func getOrderBysFromIndexes(indexes []int) (orderBys []instruction.OrderBy) {
	for _, i := range indexes {
		orderBys = append(orderBys, instruction.OrderBy{i, instruction.Ascending})
	}
	return
}
