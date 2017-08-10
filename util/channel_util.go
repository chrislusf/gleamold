package util

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"sync"
	"sync/atomic"

	"github.com/chrislusf/gleam/pb"
	"gopkg.in/vmihailenco/msgpack.v2"
)

/*
On the wire Message format, pipe:
  32 bits byte length
  []byte encoded in msgpack format

Channel Message format:
  []byte
    consecutive sections of []byte, each section is an object encoded in msgpack format

	This is not actually an array object,
	but just a consecutive list of encoded bytes for each object,
	because msgpack can sequentially decode the objects

When used by Shell scripts:
  from input channel:
    decode the msgpack-encoded []byte into strings that's tab and '\n' separated
	and feed into the shell script
  to output channel:
    encode the tab and '\n' separated lines into msgpack-format []byte
	and feed into the output channel

When used by Lua scripts:
  from input channel:
    decode the msgpack-encoded []byte into array of objects
	and pass these objects as function parameters
  to output channel:
    encode returned objects as an array of objects, into msgpack encoded []byte
	and feed into the output channel

Output Message format:
  decoded objects

Lua scripts need to decode the input and encode the output in msgpack format.
Go code also need to decode the input to "see" the data, e.g. Sort(),
and encode the output, e.g. Source().

Shell scripts via Pipe should see clear data, so the
*/

const (
	BUFFER_SIZE = 1024 * 512
)

// setup asynchronously to merge multiple channels into one channel
func CopyMultipleReaders(readers []io.Reader, writer io.Writer) (inCounter int64, outCounter int64, e error) {

	writerChan := make(chan []byte, 16*len(readers))
	errChan := make(chan error, len(readers))
	for _, reader := range readers {
		go func(reader io.Reader) {
			err := ProcessMessage(reader, func(data []byte) error {
				writerChan <- data
				atomic.AddInt64(&inCounter, 1)
				return nil
			})
			errChan <- err
		}(reader)
	}
	go func() {
		for data := range writerChan {
			if err := WriteMessage(writer, data); err != nil {
				errChan <- fmt.Errorf("WriteMessage Error: %v", err)
				break
			}
			atomic.AddInt64(&outCounter, 1)
		}
	}()
	for range readers {
		err := <-errChan
		if err != nil {
			return inCounter, outCounter, err
		}
	}
	close(writerChan)

	return inCounter, outCounter, nil
}

func LinkChannel(wg *sync.WaitGroup, inChan, outChan chan []byte) {
	wg.Add(1)
	defer wg.Done()
	for bytes := range inChan {
		outChan <- bytes
	}
	close(outChan)
}

func ReaderToChannel(wg *sync.WaitGroup, name string, reader io.ReadCloser, writer io.WriteCloser, closeOutput bool, errorOutput io.Writer) error {
	defer wg.Done()
	defer reader.Close()
	if closeOutput {
		defer writer.Close()
	}

	buf := make([]byte, BUFFER_SIZE)
	var counter int64
	err := copyBuffer(writer, reader, buf, &counter)
	if err != nil {
		fmt.Fprintf(errorOutput, "%s>Read %d bytes from input to channel: %v\n", name, counter, err)
		return err
	}
	// println("reader", name, "copied", n, "bytes.")
	return nil
}

func ChannelToWriter(wg *sync.WaitGroup, name string, reader io.Reader, writer io.WriteCloser, errorOutput io.Writer) error {
	defer wg.Done()
	defer writer.Close()

	buf := make([]byte, BUFFER_SIZE)
	var counter int64
	err := copyBuffer(writer, reader, buf, &counter)
	if err != nil {
		fmt.Fprintf(errorOutput, "%s>Moved %d bytes: %v\n", name, counter, err)
	}
	// println("writer", name, "moved", n, "bytes.")
	return err
}

func LineReaderToChannel(wg *sync.WaitGroup, stat *pb.InstructionStat, name string, reader io.Reader, ch io.WriteCloser, closeOutput bool, errorOutput io.Writer) {
	defer wg.Done()
	if closeOutput {
		defer ch.Close()
	}

	r := bufio.NewReaderSize(reader, BUFFER_SIZE)
	w := bufio.NewWriterSize(ch, BUFFER_SIZE)
	defer w.Flush()

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		stat.InputCounter++
		// fmt.Printf("%s>line input: %s\n", name, scanner.Text())
		parts := bytes.Split(scanner.Bytes(), []byte{'\t'})
		var buf bytes.Buffer
		encoder := msgpack.NewEncoder(&buf)
		encoder.Encode(Now())
		for _, p := range parts {
			if err := encoder.Encode(p); err != nil {
				if err != nil {
					fmt.Fprintf(errorOutput, "%s>Failed to encode bytes from channel to writer: %v\n", name, err)
					return
				}
			}
		}
		// fmt.Printf("%s>encoded input: %s\n", name, string(buf.Bytes()))
		WriteMessage(w, buf.Bytes())
		stat.OutputCounter++
	}
	if err := scanner.Err(); err != nil {
		// seems the program could have ended when reading the output.
		fmt.Fprintf(errorOutput, "Failed to read from input to channel: %v\n", err)
	}
}

func ChannelToLineWriter(wg *sync.WaitGroup, stat *pb.InstructionStat, name string, reader io.Reader, writer io.WriteCloser, errorOutput io.Writer) {
	defer wg.Done()
	defer writer.Close()
	w := bufio.NewWriterSize(writer, BUFFER_SIZE)
	defer w.Flush()

	r := bufio.NewReaderSize(reader, BUFFER_SIZE)

	if err := PrintDelimited(stat, r, w, "\t", "\n"); err != nil {
		fmt.Fprintf(errorOutput, "%s>Failed to decode bytes from channel to writer: %v\n", name, err)
		return
	}

}

func copyBuffer(dst io.Writer, src io.Reader, buf []byte, written *int64) (err error) {
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				atomic.AddInt64(written, int64(nw))
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	return err
}
