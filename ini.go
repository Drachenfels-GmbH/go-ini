/*
	IMPORTANT considerations	

	- Last line must be terminated with a newline

	TODO associate preceding (multi-line) comments with lines ? 
  TODO associate trailing comments to lines ?
  TODO For KeyVal split up the Value part into Value and CommentValue
  IDEA map comments to line numbers and use line number of line to collect 
  comments (same line, continuous lines before)

*/
package ini

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
)

const (
	CAP_DEFAULT = 10
)

type INI map[int]map[string]string

type INIMarshaler interface {
	MarshalINI() (ini []byte, err error)
}

type INIUnmarshaler interface {
	UnmarshalINI(ini []byte) error
}

type ValType int

const (
	COMMENT ValType = iota
	SECTION
	KEYVAL
)

// use a line reader
// FIXME protect line values against modification ?
type Line struct {
	ValType
	Pos   int
	Val   []byte
	Key   string // Non empty for ValType KEYVAL
	Value string
}

func (l *Line) String() string {
	return string(l.Val)
}

type Reader struct {
	// The last declaration of a section
	// overwrites any previous declaration.
	SectionOverwrite bool
	// The last declaration of a key overwrites
	// any previous declaration.
	KeyOverwrite bool
	// Comment character for the start of a line.
	Comment rune
	// Key value separator token
	Separator string
	// Trim leading space in every line.
	TrimLeadingSpace bool

	rdr     *bufio.Reader
	pos     int    // line counter (beginning with line 1)
	line    *Line  // unmodified line value (e.g for error reporting)
	section string // current section name
}

type LineError struct {
	line *Line
	msg  string
}

func (e *LineError) Error() string {
	return fmt.Sprintf("Error at line[%d] - %s: %q", e.line.Pos, e.msg, e.line)
}

func (r *Reader) NewLineError(s string, args ...interface{}) *LineError {
	return &LineError{
		line: r.line,
		msg:  fmt.Sprintf(s, args...),
	}
}


// Process line by line
func (r *Reader) Read() (*Line, error) {
	// minimum line length:
	// comment: 2 characters '#f'
	// section: 3 characters '[f]'
	// keyvalue: 3 characters 'k=v'

	for {
		r.pos++
		r.line = nil

		l, lineBufferExceeded, err := r.rdr.ReadLine()
		if err != nil {
			return nil, err
		}
		if lineBufferExceeded {
			return nil, r.NewLineError("Internal error, line buffer exceeded.")
		}

		s := string(l)
		// trim line
		if r.TrimLeadingSpace {
			s = strings.TrimLeftFunc(s, unicode.IsSpace)
		}
		// trailing whitespace is only removed from sections and comments

		// skip empty lines
		if len(s) == 0 {
			continue
		}

		r.line = &Line{Pos: r.pos, Val: l}

		if len(s) == 1 {
			if rune(s[0]) == r.Comment {
				continue
			} else {
				return nil, r.NewLineError("Malformed")
			}
		}

		switch rune(s[0]) {
		case r.Comment:
			s := strings.TrimSpace(s[1:])
			// skip empty comments
			if len(s) == 0 {
				continue
			} else {
				r.line.ValType = COMMENT
				r.line.Value = s
				return r.line, nil
			}
		case '[':
			// section
			s := strings.TrimRightFunc(s[1:], unicode.IsSpace)
			if len(s) > 0 && s[len(s)-1] == ']' {
				r.line.ValType = SECTION
				r.line.Value = s[:len(s)-1]
				return r.line, nil
			} else {
				return nil, r.NewLineError("Malformed section")
			}
		default:
			// keyval
			vals := strings.SplitN(s, r.Separator, 2)
			if len(vals) != 2 {
				return nil, r.NewLineError("Malformed keyval")
			}
			r.line.ValType = KEYVAL
			r.line.Key = vals[0]
			r.line.Value = vals[1]
		}
	}
}

func New() INI {
	return make(map[int]map[string]string, CAP_DEFAULT)
}

func (i *INI) Unmarshal(b []byte) error {
	return fmt.Errorf("Not implemented")
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		rdr:       bufio.NewReader(r),
		Comment:   '#',
		Separator: "=",
	}
}
