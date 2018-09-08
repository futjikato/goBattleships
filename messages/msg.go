package messages

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

// Msg is a message send to and from players
//
// This simple message protocol contains of a message type and a payload
type Msg struct {
	MsgType string
	Payload string
}

// WriteTo writes the message to the given writer
//
// This naive implementation uses a semicolon and a pipe as seperators.
// Thus those characters should not be in a type or payload
func (m *Msg) WriteTo(w io.Writer) (int64, error) {
	i, err := w.Write([]byte(m.MsgType + ";" + m.Payload + "|"))
	return int64(i), err
}

// ReadFrom reads messages from the reader
//
// To allow a message to be larger then the read buffer we keep a buffer
// around with remaining bytes that did not contain a complete message.
func ReadFrom(r io.Reader, remains []byte) ([]*Msg, []byte) {
	buf := make([]byte, 1024)
	messages := make([]*Msg, 0)

	i, err := r.Read(buf)
	if i > 0 {
		merged := append(remains, buf[:i]...)
		written := string(merged)
		nextEnd := strings.Index(written, "|")
		for nextEnd > 0 {
			msgStr := written[:nextEnd+1]
			written = written[nextEnd+1:]
			nextEnd = strings.Index(written, "|")

			msg, pErr := parse(msgStr)
			if err != nil {
				fmt.Printf("Failed to parse complete message %s: %s", msgStr, pErr)
			} else {
				messages = append(messages, msg)
			}
		}

		return messages, []byte(written)
	}
	if err != nil && err != io.EOF {
		panic(err)
	}

	return messages, remains
}

// parse a msg from a string
//
// given string must be exactly one message and start with the message type and
// end with a pipe.
// Message types can contain pipes.
// Payloads can contains semicolons
func parse(s string) (*Msg, error) {
	// regex matches in group 1 everything but a semicolon
	// followed by a semicolon
	// followed by group 2 with everything but a pipe
	// ends with a pipe
	matches := regexp.MustCompile("^([^;]+);([^|]*)|$").FindStringSubmatch(s)
	if len(matches) != 3 {
		return nil, fmt.Errorf("not enough matches %d", len(matches))
	}

	return &Msg{
		MsgType: matches[1],
		Payload: matches[2],
	}, nil
}
