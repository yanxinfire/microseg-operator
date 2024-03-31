package policy

import (
	"reflect"
	"testing"

	microsegv1 "github.com/yanxinfire/microseg-operator/api/v1"

	"github.com/projectcalico/api/pkg/lib/numorstring"

	v3 "github.com/projectcalico/api/pkg/apis/projectcalico/v3"
	projectcalicov3 "github.com/projectcalico/api/pkg/client/clientset_generated/clientset/typed/projectcalico/v3"
)

func TestCalicoEngine_buildCalicoRule(t *testing.T) {
	type fields struct {
		client projectcalicov3.ProjectcalicoV3Interface
	}
	type args struct {
		rl []microsegv1.MicrosegNetworkPolicyIngress
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		wantRule *v3.Rule
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := calicoEngine{
				client: tt.fields.client,
			}
			if gotRule, _ := e.buildIngressRule(tt.args.rl); !reflect.DeepEqual(gotRule, tt.wantRule) {
				t.Errorf("buildIngressRule() = %v, want %v", gotRule, tt.wantRule)
			}
		})
	}
}

func TestCalicoEngine_fillEntityRuleField(t *testing.T) {
	type fields struct {
		client projectcalicov3.ProjectcalicoV3Interface
	}
	type args struct {
		namespaceSelector map[string]string
		resourceSelector  map[string]string
		ipBlocks          []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := calicoEngine{
				client: tt.fields.client,
			}
			rule := e.fillEntityRuleField(tt.args.namespaceSelector, tt.args.resourceSelector, tt.args.ipBlocks)
			t.Log(rule)
		})
	}
}

func Test_rulePorts(t *testing.T) {
	singlePort, _ := numorstring.PortFromString("80")
	port90, _ := numorstring.PortFromString("90")
	port100 := numorstring.SinglePort(100)
	rangePort, _ := numorstring.PortFromRange(80, 1800)

	tests := []struct {
		name string
		args string
		want []numorstring.Port
	}{
		{
			name: "single port",
			args: "80",
			want: []numorstring.Port{singlePort},
		},
		{
			name: "multi port",
			args: "80, 90, 100,",
			want: []numorstring.Port{singlePort, port90, port100},
		},
		{
			name: "range",
			args: "80-1800",
			want: []numorstring.Port{rangePort},
		},
		{
			name: "invalid",
			args: "abd",
			want: []numorstring.Port{},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calicoRulePorts(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calicoRulePorts() = %v, want %v", got, tt.want)
			}
		})
	}
}
