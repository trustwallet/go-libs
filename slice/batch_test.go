package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSlice(t *testing.T) {
	type args struct {
		slice []uint
		size  uint
	}

	tests := []struct {
		name    string
		args    args
		want    [][]uint
		wantErr bool
	}{
		{
			"Test uint batch size of 2",
			args{
				[]uint{1, 2, 3, 4}, 2,
			},
			[][]uint{{1, 2}, {3, 4}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSlice(tt.args.slice, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
