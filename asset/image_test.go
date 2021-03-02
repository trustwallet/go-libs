package asset

import "testing"

func TestGetImageURL(t *testing.T) {
	type args struct {
		endpoint string
		asset    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Test coin",
			args{
				endpoint: "https://assets.com",
				asset:    "c60",
			},
			"https://assets.com/blockchains/ethereum/info/logo.png",
		},
		{
			"Test coin",
			args{
				endpoint: "https://assets.com",
				asset:    "c60_t123",
			},
			"https://assets.com/blockchains/ethereum/assets/123/logo.png",
		},
		{
			"Test invalid coin",
			args{
				endpoint: "https://assets.com",
				asset:    "c123",
			},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetImageURL(tt.args.endpoint, tt.args.asset); got != tt.want {
				t.Errorf("GetImageURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
