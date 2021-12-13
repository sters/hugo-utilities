package content

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	yaml "github.com/goccy/go-yaml"
	"github.com/morikuni/failure"
)

// FrontMatter is metadata for MarkdownContent.
type FrontMatter map[string]interface{}

func (f FrontMatter) getField(fieldname string) interface{} {
	field, ok := f[fieldname]
	if !ok {
		return nil
	}
	return field
}

func (f FrontMatter) GetInt(fieldname string) int {
	r, ok := f.getField(fieldname).(int)
	if !ok {
		return 0
	}

	return r
}

func (f FrontMatter) GetString(fieldname string) string {
	r, ok := f.getField(fieldname).(string)
	if !ok {
		return ""
	}

	return r
}

func (f FrontMatter) GetStrings(fieldname string) []string {
	ss, ok := f.getField(fieldname).([]string)
	if ok {
		return ss
	}

	si, ok := f.getField(fieldname).([]interface{})
	if !ok {
		return nil
	}

	result := make([]string, 0, len(si))
	for _, s := range si {
		r, ok := s.(string)
		if !ok {
			continue
		}

		result = append(result, r)
	}

	return result
}

// MarkdownContent for hugo.
// https://gohugo.io/content-management/formats/
type MarkdownContent struct {
	// FrontMatter is metadata for this content
	FrontMatter FrontMatter
	// Body for this content
	Body string
}

// Dump to string from this content.
func (m *MarkdownContent) Dump() (string, error) {
	meta, err := yaml.Marshal(m.FrontMatter)
	if err != nil {
		return "", failure.Wrap(err, failure.WithCode(ErrFileContentMismatch))
	}

	return fmt.Sprintf(`%s%s
%s%s`,
		hugoSeparator, strings.TrimSpace(string(meta)), hugoSeparator, m.Body), nil
}

// see https://gohugo.io/content-management/front-matter/#front-matter-formats
const (
	hugoSeparator     = "---\n"
	hugoContentLength = 3
)

// ErrFileContentMismatch on specified filepath.
var ErrFileContentMismatch = failure.StringCode("file content mismatch")

// ParseMarkdownWithYaml from any reader to make MarkdownContent struct.
func ParseMarkdownWithYaml(r io.Reader) (*MarkdownContent, error) {
	raw, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	content := strings.Split(string(raw), hugoSeparator)
	if len(content) < hugoContentLength {
		return nil, failure.New(ErrFileContentMismatch)
	}

	c := &MarkdownContent{
		Body: strings.Join(content[2:], hugoSeparator),
	}

	if err := yaml.Unmarshal([]byte(content[1]), &c.FrontMatter); err != nil {
		return nil, failure.Wrap(err, failure.WithCode(ErrFileContentMismatch))
	}

	return c, nil
}
