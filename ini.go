package ini

import (
	"bytes"
	"io"
)

type INIMarshaler interface {
	MarshalINI() (b []byte, err error)
}

type INIUnmarshaler interface {
	UnmarshalINI(b []byte) error
}

type Section struct {
	Header *Line // FIXME no header for default section ?
	Body   []*Line
}

func (s *Section) Append(l *Line) {
	s.Body = append(s.Body, l)
}

func(s *Section) IsDefault() bool {
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
}

func (i *INI) GetSections(name string) []*Section {
	sections := make([]*Section, 0, 5)
	for _, sec := range i.Sections {
		if ! sec.IsDefault() && sec.Header.Value == name {
			sections = append(sections, sec)
		}
	}
	return sections
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

func (i *INI) Unmarshal(b []byte) error {
	// create scanner
	s := NewScanner(bytes.NewReader(b))

	for {
		line, err := s.Scan()
		if err != nil {
			if err == io.EOF {
				break // FIXME continue processing this line ?
			} else {
				return err
			}
		}
		switch line.ValType {
		case COMMENT: // discard
		case SECTION:
			i.AddSection(line)
		case KEYVAL:
			i.CurrentSection().Append(line)
		}
	}
	return nil
}

func Parse(b []byte) (*INI, error) {
	i := New()
	return i, i.Unmarshal(b)
}
