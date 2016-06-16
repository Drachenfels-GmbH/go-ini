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

`
	rdr := NewScanner(strings.NewReader(s))
	//rdr.PanicOnError = true
	line, err := rdr.Scan()
	assert.Nil(t, err)

	assert.Equal(t, 4, line.Pos)
	assert.Equal(t, COMMENT, line.ValType)
	assert.Equal(t, "foobar", line.Value)
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

func TestScan_SECTION(t *testing.T) {
	s := `[foo]

[blubber
`
	rdr := NewScanner(strings.NewReader(s))
	line, err := rdr.Scan()
	assert.Nil(t, err)
	assert.Equal(t, 1, line.Pos)
	assert.Equal(t, SECTION, line.ValType)
	assert.Equal(t, "foo", line.Value)

	line, err = rdr.Scan()
	assert.Equal(t, 3, line.Pos)
	assert.Equal(t, SECTION, line.ValType)
	assert.Equal(t, "", line.Value)
	assert.NotNil(t, err)
}
