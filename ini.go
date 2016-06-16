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

type Scanner struct {
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
	//
	//PanicOnError bool

	rdr     *bufio.Reader
	pos     int   // line counter (beginning with line 1)
	line    *Line // unmodified line value (e.g for error reporting)
	err     error
	section string // current section name
}

// Process line by line
func (r *Scanner) Scan() (*Line, error) {
	// minimum line length:
	// comment: 2 characters '#f'
	// section: 3 characters '[f]'
	// keyvalue: 3 characters 'k=v'
	for {
		if r.err != nil {
			return r.line, r.err
		}

		r.pos++
		r.line = &Line{Pos: r.pos}

		var lineBufferExceeded bool
		r.line.Val, lineBufferExceeded, r.err = r.rdr.ReadLine()

		if lineBufferExceeded {
			return r.line, fmt.Errorf("Internal error, line buffer exceeded.")
		}

		// Stop line processing on undefined error but continue on EOF.
		if r.err != nil && r.err != io.EOF {
			return r.line, r.err
		}

		s := string(r.line.Val)
		// trim line
		if r.TrimLeadingSpace {
			s = strings.TrimLeftFunc(s, unicode.IsSpace)
		}
		// trailing whitespace is only removed from sections and comments

		// skip empty lines
		if len(s) == 0 {
			continue
		}

		if len(s) == 1 {
			if rune(s[0]) == r.Comment {
				continue
			} else {
				//if len(strings.TrimSpace(s)) == 0 {
				//	continue
				//} else {
				return r.line, fmt.Errorf("Malformed line.")
				//}
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
				return r.line, r.err
			}
		case '[':
			// section
			s := strings.TrimRightFunc(s[1:], unicode.IsSpace)
			if len(s) > 0 && s[len(s)-1] == ']' {
				r.line.ValType = SECTION
				r.line.Value = s[:len(s)-1]
				return r.line, r.err
			} else {
				return r.line, fmt.Errorf("Malformed section.")
			}
		default:
			// keyval
			vals := strings.SplitN(s, r.Separator, 2)
			if len(vals) != 2 {
				return r.line, fmt.Errorf("Malformed keyval.")
			}
			r.line.ValType = KEYVAL
			r.line.Key = vals[0]
			r.line.Value = vals[1]
			return r.line, r.err
		}
	}
}

func New() INI {
	return make(map[int]map[string]string, CAP_DEFAULT)
}

func (i *INI) Unmarshal(b []byte) error {
	return fmt.Errorf("Not implemented")
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		rdr:              bufio.NewReader(r),
		Comment:          '#',
		Separator:        "=",
		TrimLeadingSpace: true,
	}
}
