package publish

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/cavaliergopher/grab/v3"
	"github.com/morikuni/failure"
	"github.com/sters/hugo-utilities/content"
)

const targetFile = ".md"

var (
	// ErrNotTarget on specified filepath.
	ErrNotTarget = failure.StringCode("not target file")
	// ErrFileCannotLoad on specified filepath.
	ErrFileCannotLoad = failure.StringCode("file cannot load")
	// ErrFileEmpty on specified filepath.
	ErrFileEmpty = failure.StringCode("file is empty")
	// ErrFileContentMismatch on specified filepath.
	ErrFileContentMismatch = failure.StringCode("file content mismatch")
	// ErrContentIsReservedButNotDraft on specified filepath.
	ErrContentIsReservedButNotDraft = failure.StringCode("content is reserved but not draft")
	// ErrContentIsNotReserved on specified filepath.
	ErrContentIsNotReserved = failure.StringCode("content is not reserved")
	// ErrContentIsNotTheTimeYet on specified filepath.
	ErrContentIsNotTheTimeYet = failure.StringCode("content is not the time yet")

	readFile  = ioutil.ReadFile
	writeFile = ioutil.WriteFile
)

// New is constructor of Publisher.
func New(reservedKey string, draftKey string, imageReplaceFormat string) *Publisher {
	return &Publisher{
		reservedKey:        reservedKey,
		draftKey:           draftKey,
		imageReplaceFormat: imageReplaceFormat,
	}
}

// Publisher doing publish reserved content.
type Publisher struct {
	reservedKey        string
	draftKey           string
	imageReplaceFormat string
}

// CheckReservedAndPublish reserved content.
func (p *Publisher) CheckReservedAndPublish(fpath string) error {
	if !strings.Contains(fpath, targetFile) {
		return failure.New(ErrNotTarget)
	}

	rawContent, err := readFile(fpath)
	if err != nil {
		return failure.Wrap(err, failure.WithCode(ErrFileCannotLoad))
	}

	if len(rawContent) == 0 {
		return failure.New(ErrFileEmpty)
	}

	content, err := content.ParseMarkdownWithYaml(bytes.NewBuffer(rawContent))
	if err != nil {
		return failure.New(ErrFileContentMismatch)
	}

	if _, ok := content.FrontMatter[p.reservedKey]; !ok {
		return failure.New(ErrContentIsNotReserved)
	}
	if d, ok := content.FrontMatter[p.draftKey]; !ok || d != true {
		return failure.New(ErrContentIsReservedButNotDraft)
	}

	t, err := time.Parse(time.RFC3339, content.FrontMatter["date"].(string))
	if err != nil {
		return failure.Wrap(err, failure.WithCode(ErrFileContentMismatch))
	}

	now := time.Now()
	if t.UnixNano() > now.UnixNano() {
		return failure.New(ErrContentIsNotTheTimeYet)
	}

	delete(content.FrontMatter, p.reservedKey)
	delete(content.FrontMatter, p.draftKey)

	if images := detectImageURL(content.Body); len(images) > 0 {
		dir := filepath.Dir(fpath)
		for _, image := range images {
			r := regexp.MustCompile(regexp.QuoteMeta(image.raw))
			resp, err := grab.Get(dir, image.url)
			if err != nil || resp == nil {
				content.Body = r.ReplaceAllString(content.Body, "<!-- $1 (failed to get file) -->")
			} else {
				content.Body = r.ReplaceAllStringFunc(content.Body, func(s string) string {
					fp, _ := filepath.Rel(dir, resp.Filename)
					return regexp.MustCompile("^(.+)$").ReplaceAllString(fp, p.imageReplaceFormat)
				})
			}
		}
	}

	result, err := content.Dump()
	if err != nil {
		return failure.Wrap(err, failure.WithCode(ErrFileContentMismatch))
	}

	if err := writeFile(fpath, []byte(result), os.ModePerm); err != nil {
		return failure.Wrap(err)
	}

	return nil
}

const urlRegex = "!\\[.*?\\]\\((https?://[^\\s\\]]+)\\)"

type image struct {
	url string
	raw string
}

func detectImageURL(body string) []image {
	r := regexp.MustCompile(urlRegex)

	raw := r.FindAllStringSubmatch(body, -1)
	files := make([]image, len(raw))
	for i, r := range raw {
		files[i] = image{
			url: r[1],
			raw: r[0],
		}
	}

	return files
}
