package event

import (
	"bufio"
	"io"
	"sync"

	"github.com/biensupernice/krane/logger"
)

// Stream messages from an ioreader to a channel
func Stream(in *io.Reader, out chan string, done chan bool) {
	reader := bufio.NewReader(*in)

	// Reader mutual exclusive lock
	var mu sync.RWMutex
	go func() {
		for {
			mu.Lock()

			// Read lines from the reader
			str, _, err := reader.ReadLine()
			if err != nil {
				logger.Debugf("Stream Error: %s", err.Error())
				done <- true
				return
			}

			// send the lines to channel
			out <- string(str)

			mu.Unlock()
		}
	}()
}
