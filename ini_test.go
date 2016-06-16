package ini

import (
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestScan_EmptyLinesAndComments(t *testing.T) {
	s := `

#
# foobar   

#

` // FIXME adding a space in front of here returns a LineError
	rdr := NewScanner(strings.NewReader(s))
	//rdr.PanicOnError = true
	line, err := rdr.Scan()
	assert.Nil(t, err)

	assert.Equal(t, 4, line.Pos)
	assert.Equal(t, COMMENT, line.ValType)
	assert.Equal(t, "foobar", line.Value)

	// FIXME return EOF from Scaner
	//line, err = rdr.Scan()
	//assert.Nil(t, line)
	//assert.Equal(t, io.EOF, err)
}

func TestScan_KEYVAL(t *testing.T) {
	s := `foo=bar
`
	rdr := NewScanner(strings.NewReader(s))
	//rdr.PanicOnError = true
	line, err := rdr.Scan()
	assert.Nil(t, err)

	assert.Equal(t, 1, line.Pos)
	assert.Equal(t, KEYVAL, line.ValType)
	assert.Equal(t, "foo", line.Key)
	assert.Equal(t, "bar", line.Value)

	line, err = rdr.Scan()
	assert.Equal(t, 2, line.Pos)
	assert.Equal(t, io.EOF, err)
}

func TestScan_EOF(t *testing.T) {
	s := `foo=bar

`
	rdr := NewScanner(strings.NewReader(s))
	line, err := rdr.Scan()
	assert.Nil(t, err)

	line, err = rdr.Scan()
	assert.Equal(t, 3, line.Pos)
	assert.Equal(t, io.EOF, err)

	line, err = rdr.Scan()
	assert.Equal(t, 3, line.Pos)
	assert.Equal(t, io.EOF, err)
}
