package adapter

import (
	"github.com/Placons/oneapp-logger/logger"
	"io"
)

// closes http body stream while preserving error in defer
func closeBody(body io.ReadCloser, err *error, l *logger.StandardLogger) {
	l.Debug("Closing body stream.")
	e := body.Close()
	if err == nil && e != nil {
		l.Debug("Got an error while closing body stream. Returning error.")
		err = &e
		l.ErrorWithErr("Got an error while closing body stream.", *err)
	}
	l.Debug("Closed body stream.")
}
