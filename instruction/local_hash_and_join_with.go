package instruction

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/chrislusf/gleam/pb"
	"github.com/chrislusf/gleam/util"
)

func init() {
	InstructionRunner.Register(func(m *pb.Instruction) Instruction {
		if m.GetLocalHashAndJoinWith() != nil {
			return NewLocalHashAndJoinWith(
				toInts(m.GetLocalHashAndJoinWith().GetIndexes()),
			)
		}
		return nil
	})
}

type LocalHashAndJoinWith struct {
	indexes []int
}

func NewLocalHashAndJoinWith(indexes []int) *LocalHashAndJoinWith {
	return &LocalHashAndJoinWith{indexes}
}

func (b *LocalHashAndJoinWith) Name() string {
	return "LocalHashAndJoinWith"
}

func (b *LocalHashAndJoinWith) Function() func(readers []io.Reader, writers []io.Writer, stats *pb.InstructionStat) error {
	return func(readers []io.Reader, writers []io.Writer, stats *pb.InstructionStat) error {
		return DoLocalHashAndJoinWith(readers[0], readers[1], writers[0], b.indexes, stats)
	}
}

func (b *LocalHashAndJoinWith) SerializeToCommand() *pb.Instruction {
	return &pb.Instruction{
		Name: b.Name(),
		LocalHashAndJoinWith: &pb.Instruction_LocalHashAndJoinWith{
			Indexes: getIndexes(b.indexes),
		},
	}
}

func (b *LocalHashAndJoinWith) GetMemoryCostInMB(partitionSize int64) int64 {
	return int64(float32(partitionSize) * 1.1)
}

func DoLocalHashAndJoinWith(leftReader, rightReader io.Reader, writer io.Writer, indexes []int, stats *pb.InstructionStat) error {
	hashmap := make(map[string][]interface{})
	err := util.ProcessMessage(leftReader, func(input []byte) error {
		if keys, vals, err := genKeyBytesAndValues(input, indexes); err != nil {
			return fmt.Errorf("%v: %+v", err, input)
		} else {
			stats.InputCounter++
			hashmap[string(keys)] = vals
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Sort>Failed to read input data:%v\n", err)
		return err
	}
	if len(hashmap) == 0 {
		io.Copy(ioutil.Discard, rightReader)
		return nil
	}

	err = util.ProcessMessage(rightReader, func(input []byte) error {
		if ts, keys, vals, err := util.DecodeRowKeysValues(input, indexes); err != nil {
			return fmt.Errorf("%v: %+v", err, input)
		} else {
			stats.InputCounter++
			keyBytes, err := util.EncodeKeys(keys...)
			if err != nil {
				return fmt.Errorf("Failed to encoded row %+v: %v", keys, err)
			}
			if mappedValues, ok := hashmap[string(keyBytes)]; ok {
				row := append(keys, vals...)
				row = append(row, mappedValues...)
				util.WriteRow(writer, ts, row...)
				stats.OutputCounter++
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("LocalHashAndJoinWith>Failed to process the bigger input data:%v\n", err)
	}
	return err
}
