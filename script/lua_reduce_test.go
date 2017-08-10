package script

import (
	"io"
	"testing"
	"time"

	"github.com/chrislusf/gleam/util"
)

func TestLuaReduce(t *testing.T) {

	ts := time.Now().UnixNano() / int64(time.Millisecond)

	testLuaScript(
		"test ReduceBy",
		func(script Script) {
			script.Reduce(`
				function(x, y)
					return x+y
				end
			`)
		},
		func(inputWriter io.Writer) {
			util.WriteRow(inputWriter, ts, 100)
			util.WriteRow(inputWriter, ts, 101)
			util.WriteRow(inputWriter, ts, 101)
			util.WriteRow(inputWriter, ts, 103)
		},
		func(outputReader io.Reader) {
			ts1, row, _ := util.ReadRow(outputReader)
			t.Logf("row1: %+v", row)
			if ts1 != ts {
				t.Errorf("failed timestamp: %+v", row)
			}
			if !(row[0].(uint64) == 405) {
				t.Errorf("failed ReduceBy results: %+v", row)
			}
		},
	)
}

func TestLuaReduceByMultipleValues(t *testing.T) {

	ts := time.Now().UnixNano() / int64(time.Millisecond)

	testLuaScript(
		"test ReduceBy",
		func(script Script) {
			script.ReduceBy(`
				function(x, y, a, b)
					return x+a, y+b
				end
			`, []int{1})
		},
		func(inputWriter io.Writer) {
			util.WriteRow(inputWriter, ts, "key1", 100, 133)
			util.WriteRow(inputWriter, ts, "key2", 101, 3)
			util.WriteRow(inputWriter, ts, "key2", 101, 4)
			util.WriteRow(inputWriter, ts, "key3", 103, 138)
		},
		func(outputReader io.Reader) {
			ts1, row, _ := util.ReadRow(outputReader)
			if ts1 != ts {
				t.Errorf("failed timestamp: %+v", row)
			}
			t.Logf("row1: %+v", row)
			if !(row[1].(uint64) == 100 && row[2].(uint64) == 133) {
				t.Errorf("failed ReduceBy results: %+v", row)
			}
			ts1, row, _ = util.ReadRow(outputReader)
			t.Logf("row2: %+v", row)
			if !(row[1].(uint64) == 202 && row[2].(uint64) == 7) {
				t.Errorf("failed ReduceBy results: %+v", row)
			}
			ts1, row, _ = util.ReadRow(outputReader)
			t.Logf("row3: %+v", row)
			if !(row[1].(uint64) == 103 && row[2].(uint64) == 138) {
				t.Errorf("failed ReduceBy results: %+v", row)
			}
		},
	)
}

func TestLuaReduceBySingleValues(t *testing.T) {

	ts := time.Now().UnixNano() / int64(time.Millisecond)

	testLuaScript(
		"test ReduceBy",
		func(script Script) {
			script.ReduceBy(`
				function(x, y)
					return x+y
				end
			`, []int{1, 2})
		},
		func(inputWriter io.Writer) {
			util.WriteRow(inputWriter, ts, "key1", 100, 133)
			util.WriteRow(inputWriter, ts, "key2", 101, 3)
			util.WriteRow(inputWriter, ts, "key2", 101, 4)
			util.WriteRow(inputWriter, ts, "key3", 103, 138)
		},
		func(outputReader io.Reader) {
			ts1, row, _ := util.ReadRow(outputReader)
			if ts1 != ts {
				t.Errorf("failed timestamp: %+v", row)
			}
			t.Logf("row1: %+v", row)
			if !(row[2].(uint64) == 133) {
				t.Errorf("failed ReduceBy results: %+v", row)
			}
			ts1, row, _ = util.ReadRow(outputReader)
			t.Logf("row2: %+v", row)
			if !(row[2].(uint64) == 7) {
				t.Errorf("failed ReduceBy results: %+v", row)
			}
			ts1, row, _ = util.ReadRow(outputReader)
			t.Logf("row3: %+v", row)
			if !(row[2].(uint64) == 138) {
				t.Errorf("failed ReduceBy results: %+v", row)
			}
		},
	)
}

func TestLuaReduceByWithNil(t *testing.T) {

	ts := time.Now().UnixNano() / int64(time.Millisecond)

	testLuaScript(
		"test ReduceBy with nil",
		func(script Script) {
			script.ReduceBy(`
				function(x, y, a, b)
					return a, b
				end
			`, []int{1})
		},
		func(inputWriter io.Writer) {
			util.WriteRow(inputWriter, ts, "key1", 100, nil)
			util.WriteRow(inputWriter, ts, "key1", 101, 3)
		},
		func(outputReader io.Reader) {
			ts1, row, _ := util.ReadRow(outputReader)
			t.Logf("row1: %+v", row)
			if ts1 != ts {
				t.Errorf("failed timestamp: %+v", row)
			}
			if !(row[1].(uint64) == 101 && row[2].(uint64) == 3) {
				t.Errorf("failed ReduceBy results: [%s %d %d]", row...)
			}
		},
	)
}

func TestLuaGroupByMultipleValue(t *testing.T) {

	ts := time.Now().UnixNano() / int64(time.Millisecond)

	testLuaScript(
		"test GroupBy",
		func(script Script) {
			script.GroupBy([]int{1, 2})
		},
		func(inputWriter io.Writer) {
			util.WriteRow(inputWriter, ts, "key1", 100, 133, "r1")
			util.WriteRow(inputWriter, ts, "key2", 101, 3, "r2")
			util.WriteRow(inputWriter, ts, "key2", 101, 4, "r3")
			util.WriteRow(inputWriter, ts, "key3", 103, 138, "r4")
		},
		func(outputReader io.Reader) {
			ts1, row, _ := util.ReadRow(outputReader)
			if ts1 != ts {
				t.Errorf("failed timestamp: %+v", row)
			}
			t.Logf("row1: %+v", row)
			if !(row[2].([]interface{})[0].(uint64) == 133) {
				t.Errorf("failed GroupBy results: %+v", row)
			}
			ts1, row, _ = util.ReadRow(outputReader)
			t.Logf("row2: %+v", row)
			if !(row[2].([]interface{})[0].(uint64) == 3) {
				t.Errorf("failed GroupBy results: %+v", row)
			}
			ts1, row, _ = util.ReadRow(outputReader)
			t.Logf("row3: %+v", row)
			if !(row[2].([]interface{})[0].(uint64) == 138) {
				t.Errorf("failed GroupBy results: %+v", row)
			}
		},
	)
}

func TestLuaGroupByZeroValue(t *testing.T) {

	ts := time.Now().UnixNano() / int64(time.Millisecond)

	testLuaScript(
		"test GroupBy",
		func(script Script) {
			script.GroupBy([]int{1, 2})
		},
		func(inputWriter io.Writer) {
			util.WriteRow(inputWriter, ts, "key1", 100)
			util.WriteRow(inputWriter, ts, "key2", 101)
			util.WriteRow(inputWriter, ts, "key2", 101)
			util.WriteRow(inputWriter, ts, "key3", 103)
		},
		func(outputReader io.Reader) {
			ts1, row, _ := util.ReadRow(outputReader)
			t.Logf("row1: %+v", row)
			if ts1 != ts {
				t.Errorf("failed timestamp: %+v", row)
			}
			if !(row[2].(uint64) == 1) {
				t.Errorf("failed GroupBy results: %+v", row)
			}
			ts1, row, _ = util.ReadRow(outputReader)
			t.Logf("row2: %+v", row)
			if !(row[2].(uint64) == 2) {
				t.Errorf("failed GroupBy results: %+v", row)
			}
			ts1, row, _ = util.ReadRow(outputReader)
			t.Logf("row3: %+v", row)
			if !(row[2].(uint64) == 1) {
				t.Errorf("failed GroupBy results: %+v", row)
			}
		},
	)
}

func TestLuaGroupBySingleValues(t *testing.T) {

	ts := time.Now().UnixNano() / int64(time.Millisecond)

	testLuaScript(
		"test GroupBy",
		func(script Script) {
			script.GroupBy([]int{1, 2})
		},
		func(inputWriter io.Writer) {
			util.WriteRow(inputWriter, ts, "key1", 100, 133)
			util.WriteRow(inputWriter, ts, "key2", 101, 3)
			util.WriteRow(inputWriter, ts, "key2", 101, 4)
			util.WriteRow(inputWriter, ts, "key3", 103, 138)
		},
		func(outputReader io.Reader) {
			ts1, row, _ := util.ReadRow(outputReader)
			t.Logf("row1: %+v", row)
			if ts1 != ts {
				t.Errorf("failed timestamp: %+v", row)
			}
			if !(row[2].([]interface{})[0].(uint64) == 133) {
				t.Errorf("failed GroupBy results: %+v", row)
			}
			ts1, row, _ = util.ReadRow(outputReader)
			t.Logf("row2: %+v", row)
			if !(row[2].([]interface{})[0].(uint64) == 3 && row[2].([]interface{})[1].(uint64) == 4) {
				t.Errorf("failed GroupBy results: %+v", row)
			}
			ts1, row, _ = util.ReadRow(outputReader)
			t.Logf("row3: %+v", row)
			if !(row[2].([]interface{})[0].(uint64) == 138) {
				t.Errorf("failed GroupBy results: %+v", row)
			}
		},
	)
}
