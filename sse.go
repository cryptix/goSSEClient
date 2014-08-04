package goSSEClient

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/visionmedia/go-debug"
)

var dbg = debug.Debug("goSSEClient")

type SSEvent struct {
	Id   string
	Data []byte
}

func OpenSSEUrl(url string) (events chan SSEvent, err error) {
	tr := &http.Transport{
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error: resp.StatusCode == %d\n", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "text/event-stream" {
		return nil, fmt.Errorf("Error: invalid Content-Type == %s\n", resp.Header.Get("Content-Type"))
	}

	dbg("Response: %+v", resp)

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
			// dbg("NewLine: %s", string(line))

			switch {

			// start of event
			case bytes.HasPrefix(line, []byte("id:")):
				ev.Id = string(line[3:])
				dbg("eventID:", ev.Id)

			// event data
			case bytes.HasPrefix(line, []byte("data:")):
				buf.Write(line[6:])
				dbg("data: %s", string(line[5:]))

			// end of event
			case len(line) == 1:
				ev.Data = buf.Bytes()
				buf.Reset()
				events <- ev
				dbg("Event send.")
				ev = SSEvent{}

			default:
				fmt.Fprintf(os.Stderr, "Error during EventReadLoop - Default triggerd! len:%d\n%s", len(line), line)
				close(events)

			}
		}
	}()

	return events, nil
}
