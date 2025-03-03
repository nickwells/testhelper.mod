package testhelper

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

var OrigStdout = os.Stdout

// FakeIO holds the details needed to redirect and restore Stdin, Stdout and
// Stderr for the duration of a test. Create this just before the test with a
// NewStdio... func and after the test is complete restore the original
// settings by calling the Done method which will also return the contents of
// anything written to stdout and stderr.
type FakeIO struct {
	sync.Mutex
	finished bool

	origStdin  *os.File
	origStdout *os.File
	origStderr *os.File

	stdoutCh    chan []byte
	stdoutErrCh chan error

	stderrCh    chan []byte
	stderrErrCh chan error

	stdinErrCh chan error

	stdinWriter  *os.File
	stdoutReader *os.File
	stderrReader *os.File
}

// closeIfNotNil closes the passed File pointer if it isn't closeIfNotNil
func closeIfNotNil(f *os.File) {
	if f != nil {
		f.Close()
	}
}

func (fio *FakeIO) resetFakeIO() {
	closeIfNotNil(os.Stdin)
	closeIfNotNil(os.Stdout)
	closeIfNotNil(os.Stderr)

	os.Stdin = fio.origStdin
	os.Stdout = fio.origStdout
	os.Stderr = fio.origStderr

	closeIfNotNil(fio.stdinWriter)
	closeIfNotNil(fio.stdoutReader)
	closeIfNotNil(fio.stderrReader)
}

// reader will copy the contents of the passed file into a buffer passing any
// errors back over the errCh and any bytes read over the byteCh.
func reader(name string, r *os.File, byteCh chan []byte, errCh chan error) {
	var b bytes.Buffer
	if _, err := io.Copy(&b, r); err != nil {
		errCh <- fmt.Errorf("Error copying from %s: %w", name, err)
	}
	byteCh <- b.Bytes()

	r.Close()
	close(byteCh)
	close(errCh)
}

// writer will write the given byte slice to the passed File passing any
// errors back over the errCh.
func writer(name string, w *os.File, b []byte, errCh chan error) {
	if _, err := w.Write(b); err != nil {
		errCh <- fmt.Errorf("Error writing to %s: %w", name, err)
	}

	w.Close()
	close(errCh)
}

// NewStdioFromString will create a Stdio object which will provide access to
// the contents of anything written to stdout or stderr. After this has been
// called any code reading from stdin will get the contents of the passed
// string. Any output to stdout or stderr will be captured
func NewStdioFromString(input string) (fio *FakeIO, err error) {
	fio = &FakeIO{
		origStdin:  os.Stdin,
		origStdout: os.Stdout,
		origStderr: os.Stderr,
	}
	os.Stdin = nil
	os.Stdout = nil
	os.Stderr = nil

	defer func() {
		if err != nil {
			fio.resetFakeIO()
			fio = nil
		}
	}()

	os.Stdin, fio.stdinWriter, err = os.Pipe()
	if err != nil {
		err = fmt.Errorf("Cannot create the stdin pipe: %w", err)
		return
	}

	fio.stdinErrCh = make(chan error)

	go writer("stdin", fio.stdinWriter, []byte(input), fio.stdinErrCh)

	fio.stdoutReader, os.Stdout, err = os.Pipe()
	if err != nil {
		err = fmt.Errorf("Cannot create the stdout pipe: %w", err)
		return
	}

	fio.stdoutCh = make(chan []byte)
	fio.stdoutErrCh = make(chan error)

	go reader("stdout", fio.stdoutReader, fio.stdoutCh, fio.stdoutErrCh)

	fio.stderrReader, os.Stderr, err = os.Pipe()
	if err != nil {
		err = fmt.Errorf("Cannot create the stderr pipe: %w", err)
		return
	}

	fio.stderrCh = make(chan []byte)
	fio.stderrErrCh = make(chan error)

	go reader("stderr", fio.stderrReader, fio.stderrCh, fio.stderrErrCh)

	return
}

// Done tidies up, restores the std IO values to their previous settings and
// returns anything written to stdout and stderr. It is an error to call this
// twice on the same FakeIO.
func (fio *FakeIO) Done() (stdout, stderr []byte, err error) {
	if fio == nil {
		err = errors.New("FakeIO.Done - nil pointer")
		return
	}

	fio.Lock()
	if fio.finished {
		err = errors.New("FakeIO.Done - already called")
	}

	fio.finished = true
	fio.Unlock()

	fio.stdinWriter.Close()

	os.Stdout.Close()

	stdout = <-fio.stdoutCh

	if tmpErr, ok := <-fio.stdoutErrCh; ok {
		err = tmpErr
	}

	os.Stderr.Close()

	stderr = <-fio.stderrCh

	if tmpErr, ok := <-fio.stderrErrCh; ok {
		err = errors.Join(err, tmpErr)
	}

	if tmpErr, ok := <-fio.stdinErrCh; ok && !errors.Is(tmpErr, os.ErrClosed) {
		err = errors.Join(err, tmpErr)
	}

	fio.stdinWriter.Close()
	fio.stdoutReader.Close()
	fio.stderrReader.Close()

	fio.stdinWriter = nil
	fio.stdoutReader = nil
	fio.stderrReader = nil

	fio.resetFakeIO()

	return
}
