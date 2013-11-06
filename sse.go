package goSSEClient

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"
)

type SSEvent struct {
	Id   string
	Data []byte
}

func OpenSSEUrl(url string) (events chan SSEvent, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error: resp.StatusCode == %d\n", resp.StatusCode)
	}

	events = make(chan SSEvent)
	var buf bytes.Buffer

	go func() {
		ev := SSEvent{}

		buffed := bufio.NewReader(resp.Body)
		for {
			line, err := buffed.ReadBytes('\n')
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error during resp.Body read:%s\n", err)
				close(events)
			}

			switch {

			// start of event
			case bytes.HasPrefix(line, []byte("id:")):
				ev.Id = string(line[3:])

			// event data
			case bytes.HasPrefix(line, []byte("data:")):
				buf.Write(line[6:])

			// end of event
			case len(line) == 1:
				ev.Data = buf.Bytes()
				buf.Reset()
				events <- ev
				ev = SSEvent{}

			default:
				fmt.Fprintf(os.Stderr, "Error during EventReadLoop - Default triggerd! len:%d\n%s", len(line), line)
				close(events)

			}
		}
	}()

	return events, nil
}
