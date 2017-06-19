package manager

import (
	"fmt"
	"testing"
)

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		input string
		want  error
	}{
		{
			input: "A",
			want:  nil,
		},
		{
			input: "CNAME",
			want:  nil,
		},
		{
			input: "AAAA",
			want:  fmt.Errorf("Invalid record type 'AAAA', must be CNAME or A"),
		},
	}
	for _, test := range tests {
		config := Config{RecordType: test.input}
		got := config.Validate()
		if got != nil && test.want != nil && got.Error() != test.want.Error() {
			t.Errorf(`Invalid result from Config{}.Validate(), Got: "%v", expected "%v"`, got, test.want)
		}
	}
}
