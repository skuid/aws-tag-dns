package manager

import (
	"reflect"
	"testing"
)

func TestBuildRoutes(t *testing.T) {
	tests := []struct {
		input []string
		want  map[string]string
	}{
		{
			input: []string{},
			want:  map[string]string{},
		},
		{
			input: []string{"10.100.0.1", "10.0.100.1", "10.101.0.1"},
			want: map[string]string{
				"etcd0.skuid.com": "10.0.100.1",
				"etcd1.skuid.com": "10.100.0.1",
				"etcd2.skuid.com": "10.101.0.1",
			},
		},
		{
			input: []string{"10.100.0.1", "10.100.0.1", "10.100.0.1"},
			want: map[string]string{
				"etcd0.skuid.com": "10.100.0.1",
				"etcd1.skuid.com": "10.100.0.1",
				"etcd2.skuid.com": "10.100.0.1",
			},
		},
	}
	m := Manager{
		Domain: "skuid.com",
		config: Config{
			SubdomainPrefix: "etcd",
		},
	}
	for _, test := range tests {
		got := m.buildRoutes(test.input)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("buildRoutes returned unexpected results, got: %v, want %v", got, test.want)
		}
	}
}
