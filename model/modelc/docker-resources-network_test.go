package modelc

import (
	"testing"

	"github.com/pacit/dde/common"
)

func TestDockerResourceNetwork_CompileTemplatesOrDie(t *testing.T) {
	type fields struct {
		Name       string
		Driver     string
		Subnet     string
		Gateway    string
		IpRange    string
		Attachable bool
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
			fields:   fields{"name-{{.p1}}", "driver-{{.p2}}", "subnet-{{.p1}}", "gateway-{{.p2}}", "ip-{{.p1}}", true},
			expected: fields{"name-val1", "driver-val2", "subnet-val1", "gateway-val2", "ip-val1", true},
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
			dr := &DockerResourceNetwork{
				Name:       tt.fields.Name,
				Driver:     tt.fields.Driver,
				Subnet:     tt.fields.Subnet,
				Gateway:    tt.fields.Gateway,
				IpRange:    tt.fields.IpRange,
				Attachable: tt.fields.Attachable,
			}
			dr.CompileTemplatesOrDie(tt.args.ctx, tt.args.props)
			if dr.Name != tt.expected.Name {
				t.Fatal("Name is not correct compile")
			}
			if dr.Driver != tt.expected.Driver {
				t.Fatal("Driver is not correct compile")
			}
			if dr.Subnet != tt.expected.Subnet {
				t.Fatal("Subnet is not correct compile")
			}
			if dr.Gateway != tt.expected.Gateway {
				t.Fatal("Gateway is not correct compile")
			}
			if dr.IpRange != tt.expected.IpRange {
				t.Fatal("IpRange is not correct compile")
			}
		})
	}
}
