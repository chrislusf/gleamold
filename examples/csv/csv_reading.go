package main

import (
	"github.com/chrislusf/gleamold/distributed"
	. "github.com/chrislusf/gleamold/flow"
	"github.com/chrislusf/gleamold/gio"
	"github.com/chrislusf/gleamold/plugins/csv"
)

func main() {

	gio.Init()

	f := New()

	a := f.Read(csv.New("a?.csv", 3).SetHasHeader(true)).Select(Field(1, 2, 3)).Hint(TotalSize(17))

	b := f.Read(csv.New("b*.csv", 3)).Select(Field(1, 4, 5)).Hint(PartitionSize(13))

	join := a.RightOuterJoin(b).Printlnf("%s : %s %s, %s %s")

	// join.Run(distributed.Planner())

	join.Run(distributed.Option())

}
