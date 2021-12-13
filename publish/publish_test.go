package publish

import (
	"os"
	"testing"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"
)

func TestPublisher_CheckReservedAndPublish(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantCode failure.StringCode
	}{
		{"empty", "", ErrFileEmpty},
		{"missing reserved", `
---
date: 2021-01-01T00:00:00Z
title: Happy New Year!
draft: true
---
Happy New Year!
`, ErrContentIsNotReserved},
		{"missing draft", `
---
date: 2021-01-01T00:00:00Z
title: Happy New Year!
reserved: true
---
Happy New Year!
`, ErrContentIsReservedButNotDraft},
		{"not draft", `
---
date: 2021-01-01T00:00:00Z
title: Happy New Year!
reserved: true
draft: false
---
Happy New Year!
`, ErrContentIsReservedButNotDraft},
		{"is not the time", `
---
date: 2100-01-01T00:00:00Z
title: Happy New Year!
reserved: true
draft: true
---
Happy New Year!
`, ErrContentIsNotTheTimeYet},
		{"success", `
---
date: 2021-01-01T00:00:00Z
title: Happy New Year!
reserved: true
draft: true
---
Happy New Year!
`, ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// fake
			readFile = func(filename string) ([]byte, error) {
				return []byte(tt.content), nil
			}
			writeFile = func(filename string, data []byte, perm os.FileMode) error {
				return nil
			}

			err := New("reserved", "draft", "%s").CheckReservedAndPublish("dummy.md")
			if tt.wantCode == "" {
				if err != nil {
					t.Errorf("want no error, got=%+v", err)
				}
				return
			}
			if !failure.Is(err, tt.wantCode) {
				t.Errorf("want error=%+v, got=%+v", tt.wantCode, err)
			}
		})
	}
}

func Test_detectImageURL(t *testing.T) {
	tests := []struct {
		name string
		body string
		want []image
	}{
		{
			name: "no url",
			body: "hello world!",
			want: []image{},
		},
		{
			name: "1 url but not image",
			body: "hello world! https://example.com/",
			want: []image{},
		},
		{
			name: "1 url",
			body: "hello world! ![](https://example.com/example.png)",
			want: []image{
				{
					url: "https://example.com/example.png",
					raw: "![](https://example.com/example.png)",
				},
			},
		},
		{
			name: "2 url but 1 url not image",
			body: "https://example.com/ hello world! ![](https://example.com/example.png)",
			want: []image{
				{
					url: "https://example.com/example.png",
					raw: "![](https://example.com/example.png)",
				},
			},
		},
		{
			name: "2 url",
			body: "![](https://example.com/example.jpeg) hello world! ![](https://example.com/example.png)",
			want: []image{
				{
					url: "https://example.com/example.jpeg",
					raw: "![](https://example.com/example.jpeg)",
				},
				{
					url: "https://example.com/example.png",
					raw: "![](https://example.com/example.png)",
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := detectImageURL(tt.body)
			assert.Equal(t, tt.want, got)
		})
	}
}
