package modelc

import (
	"testing"

	"github.com/pacit/dde/common"
)

func TestDockerResourceConfig_CompileTemplatesOrDie(t *testing.T) {
	type fields struct {
		Name  string
		Value string
		File  string
	}
	type args struct {
		ctx   common.DCtx
		props map[string]string
	}
	tests := []struct {
		name     string
		fields   fields
		expected fields
		args     args
	}{
		{
			name:     "all fields",
			fields:   fields{"name-{{.p1}}", "value{{.p2}}", "file-{{.p1}}/{{.p2}}"},
			expected: fields{"name-val1", "valueval2", "file-val1/val2"},
			args: args{
				ctx: common.DCtx{},
				props: map[string]string{
					"p1": "val1",
					"p2": "val2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dr := &DockerResourceConfig{
				Name:  tt.fields.Name,
				Value: tt.fields.Value,
				File:  tt.fields.File,
			}
			dr.CompileTemplatesOrDie(tt.args.ctx, tt.args.props)
			if dr.Name != tt.expected.Name {
				t.Fatal("Name is not correct compile")
			}
			if dr.Value != tt.expected.Value {
				t.Fatal("Value is not correct compile")
			}
			if dr.File != tt.expected.File {
				t.Fatal("File is not correct compile")
			}
		})
	}
}
