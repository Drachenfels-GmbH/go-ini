package ini

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"bytes"
)

var s = `[foobar]
foo=bar
hello=world

[foobar1]
foo1=bar1
hello1=world1

[foobar]
foo=barA
hello=worldA
foo=barB
blub=bla

`

func TestUnmarshal(t *testing.T) {
	i, err := Unmarshal([]byte(s))
	assert.Nil(t, err)

	assert.Equal(t, 4, len(i.Sections))
	assert.Equal(t, 0, len(i.Sections[0].Body))
	assert.True(t, i.Sections[0].IsDefault())

	assert.Equal(t, "foobar", i.Sections[1].Header.Value)
	assert.Equal(t, "foobar1", i.Sections[2].Header.Value)
	assert.Equal(t, "foobar", i.Sections[3].Header.Value)
}

func TestSections(t *testing.T) {
	i, err := Unmarshal([]byte(s))
	assert.Nil(t, err)

	sections := i.GetSections("foobar")
	assert.Equal(t, 2, len(sections))

	assert.Equal(t, "bar", sections[0].Values("foo")[0])
	assert.Equal(t, "barA", sections[1].Values("foo")[0])
	assert.Equal(t, "barB", sections[1].Values("foo")[1])
}

func TestValue(t *testing.T) {
	i, err := Unmarshal([]byte(s))
	assert.Nil(t, err)
	assert.False(t, i.SectionOverwrite)
	assert.False(t, i.ValueOverwrite)
	val, exist := i.Value("foobar", "foo")
	assert.True(t, exist)
	assert.Equal(t, "bar", val)
}

func TestValue_WithSectionOverwrite(t *testing.T) {
	s := NewScanner(bytes.NewReader([]byte(s)))
	i := New()
	i.SectionOverwrite = true

	err := i.Unmarshal(s)
	assert.Nil(t, err)

	val, exist := i.Value("foobar", "foo")
	assert.True(t, exist)
	assert.Equal(t, "barA", val)
}

func TestValue_WithSectionAndValueOverwrite(t *testing.T) {
	s := NewScanner(bytes.NewReader([]byte(s)))
	i := New()
	i.SectionOverwrite = true
	i.ValueOverwrite = true

	err := i.Unmarshal(s)
	assert.Nil(t, err)

	val, exist := i.Value("foobar", "foo")
	assert.True(t, exist)
	assert.Equal(t, "barB", val)
}
