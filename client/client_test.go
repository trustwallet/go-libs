package client

import (
	"testing"
)

func TestRequest_GetBase(t *testing.T) {
	type fields struct {
		baseUrl string
	}
	tests := []struct {
		name   string
		fields fields
		path   string
		want   string
	}{
		{
			name: "Test base url ends with /, path starts with /",
			fields: fields{
				baseUrl: "https://api.example.com/",
			},
			path: "/v1/account/0x32Be343B94f860124dC4fEe278FDCBD38C102D88",
			want: "https://api.example.com/v1/account/0x32Be343B94f860124dC4fEe278FDCBD38C102D88",
		},
		{
			name: "Test only base url ends with /",
			fields: fields{
				baseUrl: "https://api.example.com/",
			},
			path: "v1/account/0x32Be343B94f860124dC4fEe278FDCBD38C102D88",
			want: "https://api.example.com/v1/account/0x32Be343B94f860124dC4fEe278FDCBD38C102D88",
		},
		{
			name: "Test only path starts with /",
			fields: fields{
				baseUrl: "https://api.example.com",
			},
			path: "/v1/account/0x32Be343B94f860124dC4fEe278FDCBD38C102D88",
			want: "https://api.example.com/v1/account/0x32Be343B94f860124dC4fEe278FDCBD38C102D88",
		},
		{
			name: "Test none /",
			fields: fields{
				baseUrl: "https://api.example.com",
			},
			path: "v1/account/0x32Be343B94f860124dC4fEe278FDCBD38C102D88",
			want: "https://api.example.com/v1/account/0x32Be343B94f860124dC4fEe278FDCBD38C102D88",
		},
		{
			name: "Test empty path",
			fields: fields{
				baseUrl: "https://api.example.com/",
			},
			path: "",
			want: "https://api.example.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := InitClient(tt.fields.baseUrl, nil)
			if got := r.GetBase(tt.path); got != tt.want {
				t.Errorf("Request.GetBase() = %v, want %v", got, tt.want)
			}
		})
	}
}
