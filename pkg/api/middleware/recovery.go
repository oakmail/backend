// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style

package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"runtime"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/mattn/go-sqlite3"
	"github.com/oakmail/goqu"
	"github.com/oakmail/logrus"

	"github.com/oakmail/backend/pkg/api/errors"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// Recovery handles panics in the API
func (i *Impl) Recovery(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			if _, ok := err.(error); ok {
				_, file, line, _ := runtime.Caller(3)

				switch err.(type) {
				case
					*pq.Error,
					sqlite3.Error,
					sqlite3.ErrNo,
					sqlite3.ErrNoExtended,
					goqu.GoquError,
					goqu.EncodeError:

					if strings.Contains(file, "database/utils.go") {
						_, file, line, _ = runtime.Caller(4)
					}

					i.Log.WithFields(logrus.Fields{
						"error":    err.(error).Error(),
						"location": file + ":" + strconv.Itoa(line),
					}).Error("Database error")

					errors.Abort(c, http.StatusInternalServerError, errors.DatabaseError)
					return
				}
			}

			stack := stack(3)
			httpreq, _ := httputil.DumpRequest(c.Request, false)
			i.Log.WithFields(logrus.Fields{
				"request": string(httpreq),
				"error":   err,
				"stack":   string(stack),
			}).Error("Panic recovered")
			c.AbortWithStatus(500)
		}
	}()

	c.Next()
}

// stack returns a nicely formated stack frame, skipping skip frames
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
