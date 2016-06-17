package ini

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
)

var UNAVAIL = fmt.Errorf("Value unavailable.")

type INIMarshaler interface {
	MarshalINI() (b []byte, err error)
}

type INIUnmarshaler interface {
	UnmarshalINI(b []byte) error
}

type Section struct {
	Header *Line
	Body   []*Line
}

func (s *Section) Append(l *Line) {
	s.Body = append(s.Body, l)
}

func (s *Section) IsDefault() bool {
	return s.Header == nil
}

func (s *Section) Values(key string) []string {
	vals := make([]string, 0, 5)
	for _, l := range s.Body {
		if l.ValType == KEYVAL && l.Key == key {
			vals = append(vals, l.Value)
		}
	}
	return vals
}

func (i *INI) AddSection(header *Line) *Section {
	sec := &Section{
		Header: header,
		Body:   make([]*Line, 0, CAP_KEYVALS),
	}
	i.Sections = append(i.Sections, sec)
	return sec
}

type INI struct {
	Sections []*Section
	// The last declaration of a section has precedence
	SectionOverwrite bool
	// The last declaration of a key has precedence.
	ValueOverwrite bool
}

func (i *INI) GetSections(name string) []*Section {
	sections := make([]*Section, 0, 5)
	for _, sec := range i.Sections {
		if !sec.IsDefault() && sec.Header.Value == name {
			sections = append(sections, sec)
		}
	}
	return sections
}

func (s *Section) FirstValue(key string) (string, bool) {
	for _, k := range s.Body {
		if k.Key == key {
			return k.Value, true
		}
	}
	return "", false
}

func (s *Section) LastValue(key string) (string, bool) {
	for i := len(s.Body) - 1; i > 0; i-- {
		k := s.Body[i]
		if k.Key == key {
			return k.Value, true
		}
	}
	return "", false
}

func (s *Section) Value(key string, overwrite bool) (string, bool) {
	if overwrite {
		return s.LastValue(key)
	} else {
		return s.FirstValue(key)
	}
}

func (i *INI) getVal(sec *Section, sectionName, key string) (string, bool) {
	if !sec.IsDefault() && sec.Header.Value == sectionName {
		return sec.Value(key, i.ValueOverwrite)
	}
	return "", false
}

func (i *INI) Value(sectionName, key string) (string, error) {
	// TODO add support for default section !!!!
	if i.SectionOverwrite {
		for x := len(i.Sections) - 1; x > 0; x-- {
			if val, ok := i.getVal(i.Sections[x], sectionName, key); ok {
				return val, nil
			}
		}
	} else {
		for _, sec := range i.Sections {
			if val, ok := i.getVal(sec, sectionName, key); ok {
				return val, nil
			}
		}
	}
	return "", UNAVAIL
}

func (i *INI) IntValue(sectionName, key string) (int, error) {
	val, err := i.Value(sectionName, key)
	if err == nil {
		return strconv.Atoi(val)
	} else {
		return -1, err
	}
}

func New() *INI {
	ini := &INI{
		Sections: make([]*Section, 0, CAP_SECTIONS),
	}
	ini.AddSection(nil) // add default section
	return ini
}

func (i *INI) CurrentSection() *Section {
	return i.Sections[len(i.Sections)-1]
}

func (i *INI) Unmarshal(s *Scanner) error {
	for {
		line, err := s.Scan()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		switch line.ValType {
		case SECTION:
			i.AddSection(line)
		case KEYVAL:
			i.CurrentSection().Append(line)
		}
	}
	return nil
}

// Unmarshal using a Scanner with default settings.
func Unmarshal(b []byte) (*INI, error) {
	i := New()
	s := NewScanner(bytes.NewReader(b))
	return i, i.Unmarshal(s)
}

// Unmarshal from file.
func Load(filePath string) (*INI, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return Unmarshal(b)
}
