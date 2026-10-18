package vrc_osc

import (
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

func TestInfo_Complete(t *testing.T) {
	tests := []struct {
		name string
		info Info
		want bool
	}{
		{
			name: "complete",
			info: Info{"127.0.0.1", 9000},
			want: true,
		},
		{
			name: "incomplete",
			want: false,
		},
		{
			name: "no host",
			info: Info{Port: 9000},
			want: false,
		},
		{
			name: "no port",
			info: Info{Host: "127.0.0.1"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.EqualT(t, tt.want, tt.info.Complete())
		})
	}
}
