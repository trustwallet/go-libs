package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPath_String(t *testing.T) {
	type fields struct {
		template string
		values   []any
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "empty template, empty values",
			fields: fields{
				template: "",
				values:   nil,
			},
			want: "",
		},
		{
			name: "empty template only",
			fields: fields{
				template: "",
				values:   []any{1, 2, 3},
			},
			want: "%!(EXTRA int=1, int=2, int=3)",
		},
		{
			name: "empty values only",
			fields: fields{
				template: "/api/v1/blocks",
				values:   nil,
			},
			want: "/api/v1/blocks",
		},
		{
			name: "both exist",
			fields: fields{
				template: "/nft/collections/%s/tokens",
				values:   []any{"123"},
			},
			want: "/nft/collections/123/tokens",
		},
		{
			name: "missing values",
			fields: fields{
				template: "/nft/collections/%s/tokens/%d",
				values:   []any{"123"},
			},
			want: "/nft/collections/123/tokens/%!d(MISSING)",
		},
		{
			name: "multiple values",
			fields: fields{
				template: "/nft/collections/%s/tokens/%s",
				values:   []any{"123", "bnb"},
			},
			want: "/nft/collections/123/tokens/bnb",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Path{
				template: tt.fields.template,
				values:   tt.fields.values,
			}
			assert.Equalf(t, tt.want, p.String(), "String()")
		})
	}
}
