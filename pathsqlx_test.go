package pathsqlx

import (
	"testing"

	_ "github.com/lib/pq"
)

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "// TODO: Add test cases."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
