package content

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseMarkdownWithYaml(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    *MarkdownContent
		wantErr bool
	}{
		{
			name:    "empty",
			raw:     "",
			want:    nil,
			wantErr: true,
		},
		{
			name: "simple",
			raw: `---
foo: bar
baz: 1
---
foo
`,
			want: &MarkdownContent{
				FrontMatter: map[string]interface{}{
					"foo": "bar",
					"baz": uint64(1),
				},
				Body: `foo
`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMarkdownWithYaml(bytes.NewBuffer([]byte(tt.raw)))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMarkdownWithYaml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("ParseMarkdownWithYaml() got = %+v, want %+v", got, tt.want)
				return
			}
		})
	}
}

func TestMarkdownContent_Dump(t *testing.T) {
	tests := []struct {
		name    string
		m       *MarkdownContent
		want    string
		wantErr bool
	}{
		{
			name: "simple",
			m: &MarkdownContent{
				FrontMatter: map[string]interface{}{
					"foo": "bar",
					"baz": 1,
				},
				Body: `foo
bar
`,
			},
			want: `---
baz: 1
foo: bar
---
foo
bar
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Dump()
			if (err != nil) != tt.wantErr {
				t.Errorf("got error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarkdownContent_GetInt(t *testing.T) {
	tests := []struct {
		name  string
		m     *MarkdownContent
		field string
		want  int
	}{
		{
			name: "simple",
			m: &MarkdownContent{
				FrontMatter: map[string]interface{}{
					"foo": "bar",
					"baz": 1,
				},
				Body: `
foo
bar
`,
			},
			field: "foo",
			want:  0,
		},
		{
			name: "simple",
			m: &MarkdownContent{
				FrontMatter: map[string]interface{}{
					"foo": "bar",
					"baz": 1,
				},
				Body: `
foo
bar
`,
			},
			field: "baz",
			want:  1,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.FrontMatter.GetInt(tt.field)
			if got != tt.want {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarkdownContent_GetString(t *testing.T) {
	tests := []struct {
		name  string
		m     *MarkdownContent
		field string
		want  string
	}{
		{
			name: "simple",
			m: &MarkdownContent{
				FrontMatter: map[string]interface{}{
					"foo": "bar",
					"baz": 1,
				},
				Body: `
foo
bar
`,
			},
			field: "foo",
			want:  "bar",
		},
		{
			name: "simple",
			m: &MarkdownContent{
				FrontMatter: map[string]interface{}{
					"foo": "bar",
					"baz": 1,
				},
				Body: `
foo
bar
`,
			},
			field: "baz",
			want:  "",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.FrontMatter.GetString(tt.field)
			if got != tt.want {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarkdownContent_GetStrings(t *testing.T) {
	tests := []struct {
		name  string
		m     *MarkdownContent
		field string
		want  []string
	}{
		{
			name: "simple",
			m: &MarkdownContent{
				FrontMatter: map[string]interface{}{
					"foo": "bar",
					"baz": 1,
					"new": []string{
						"foo",
						"bar",
						"baz",
					},
				},
				Body: `
foo
bar
`,
			},
			field: "foo",
			want:  nil,
		},
		{
			name: "simple",
			m: &MarkdownContent{
				FrontMatter: map[string]interface{}{
					"foo": "bar",
					"baz": 1,
					"new": []string{
						"foo",
						"bar",
						"baz",
					},
				},
				Body: `
foo
bar
`,
			},
			field: "baz",
			want:  nil,
		},
		{
			name: "simple",
			m: &MarkdownContent{
				FrontMatter: map[string]interface{}{
					"foo": "bar",
					"baz": 1,
					"new": []string{
						"foo",
						"bar",
						"baz",
					},
				},
				Body: `
foo
bar
`,
			},
			field: "new",
			want: []string{
				"foo",
				"bar",
				"baz",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.FrontMatter.GetStrings(tt.field)
			if !cmp.Equal(tt.want, got) {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}
