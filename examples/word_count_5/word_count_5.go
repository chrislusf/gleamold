// word_count.go
package main

import (
	"github.com/chrislusf/gleamold/flow"
)

func main() {

	flow.New().TextFile("/etc/passwd").FlatMap(`
		function(line)
			return line:gmatch("%w+")
		end
	`).Pipe("tr 'A-Z' 'a-z'").Pipe("tee x.out").Pipe("sort").Pipe("uniq -c").Printlnf("%s").Run()
}
