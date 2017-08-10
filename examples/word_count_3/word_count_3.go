// word_count.go
package main

import (
	"log"
	"path/filepath"

	"github.com/chrislusf/gleam/flow"
)

func main() {

	fileNames, err := filepath.Glob("/Users/chris/Downloads/txt/en/ep-08-03-*.txt")
	if err != nil {
		log.Fatal(err)
	}

	flow.New().Strings(fileNames).Partition(3).PipeAsArgs("cat $1").FlatMap(`
      function(line)
	    -- log("input:"..line)
        return line:gmatch("%w+")
      end
    `).Map(`
      function(word)
        return word, 1
      end
    `).ReduceBy(`
      function(x, y)
        return x + y
      end
    `).Map(`
      function(k, v)
        return k .. " " .. v
      end
    `).Pipe("sort -n -k 2").Printlnf("%s").Run()

}
