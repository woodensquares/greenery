package testhelper

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
)

type outGrabber struct {
	saveOut      **os.File
	saveOutValue *os.File

	outChannel chan string
	outWriter  *os.File
}

func chanReader(ch chan string, r io.ReadCloser) {
	var buf bytes.Buffer

	// If there are errors this should show up in the text comparison later on
	// in the tests, so send the error string through the channel also.
	if _, err := io.Copy(&buf, r); err != nil {
		ch <- err.Error()
		return
	}

	if err := r.Close(); err != nil {
		ch <- err.Error()
		return
	}

	ch <- buf.String()
}

// NewGrabber is FIXME:DOC
func NewGrabber() *outGrabber {
	return &outGrabber{}
}

// Start will start saving output for an *os.File like stdout/stderr backing
// up the descriptor so it can be restored later.
func (sc *outGrabber) Start(out **os.File) error {
	sc.saveOut = out
	sc.saveOutValue = *out

	w, err := sc.StartFile()
	if err != nil {
		return err
	}

	// Substitute our values and set up the grabber goroutine
	*out = w
	return nil
}

// StartFile will return an os.File that saves the output for later analysis.
func (sc *outGrabber) StartFile() (*os.File, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	sc.outWriter = w
	sc.outChannel = make(chan string)
	go chanReader(sc.outChannel, r)
	return w, nil
}

// Stop will stop, restore the original descriptor if set, and return the
// captured output.
func (sc *outGrabber) Stop() (string, error) {
	if sc.saveOut != nil {
		*sc.saveOut = sc.saveOutValue
		sc.saveOut = nil
	}

	if sc.outWriter == nil {
		return "", fmt.Errorf("Start was not called before calling Stop")
	}

	if err := sc.outWriter.Close(); err != nil {
		return "", errors.WithMessage(err, "error closing the writer")
	}

	sc.outWriter = nil
	return <-sc.outChannel, nil
}
