package ini

import (
	"strings"
	"testing"
	"github.com/stretchr/testify/assert"
	"io"
)

func TestRead_EmptyLinesAndComments(t *testing.T) {
	s := `

#
# foobar   

#

` // FIXME adding a space in front of here returns a LineError
	rdr := NewReader(strings.NewReader(s))
	line, err := rdr.Read()
	assert.Nil(t, err)

	assert.Equal(t, 4, line.Pos)
	assert.Equal(t, COMMENT, line.ValType)
	assert.Equal(t, "foobar", line.Value)

	line, err = rdr.Read()
	assert.Nil(t, line)
	assert.Equal(t, io.EOF, err)
}

func TestRead_KEYVAL(t *testing.T) {
	s := `foo=bar
	`
	rdr := NewReader(strings.NewReader(s))
	line, err := rdr.Read()
	assert.Nil(t, err)

	assert.Equal(t, 1, line.Pos)
	assert.Equal(t, KEYVAL, line.ValType)
	assert.Equal(t, "foo", line.Key)
	assert.Equal(t, "bar", line.Value)
}
