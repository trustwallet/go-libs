package address

import "testing"

func TestEIP55Checksum(t *testing.T) {
	type args struct {
		unchecksummed string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"test checksum 1", args{"checktest"}, "0xChecKTeSt", false},
		{"test checksum 2", args{"trustwallet"}, "0xtrUstWaLlET", false},
		{"test checksum number", args{"16345785d8a0000"}, "0x16345785d8A0000", false},
		{"test checksum hex", args{"fffdefefed"}, "0xFfFDEfeFeD", false},
		{"test checksum 3", args{"0x0000000000000000003731342d4f4e452d354639"}, "0x0000000000000000003731342d4f4E452d354639", false},
		{"test checksum 4", args{"0000000000000000003731342d4f4e452d354639"}, "0x0000000000000000003731342d4f4E452d354639", false},
		{"test checksum Ethereum address", args{"0x84a0d77c693adabe0ebc48f88b3ffff010577051"}, "0x84A0d77c693aDAbE0ebc48F88b3fFFF010577051", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EIP55Checksum(tt.args.unchecksummed)
			if (err != nil) != tt.wantErr {
				t.Errorf("EIP55Checksum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EIP55Checksum() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemove0x(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"remove 0x from addres", "0x158079ee67fce2f58472a96584a73c7ab9ac95c1", "158079ee67fce2f58472a96584a73c7ab9ac95c1"},
		{"remove 0x from hash", "0x230798fe22abff459b004675bf827a4089326a296fa4165d0c2ad27688e03e0c", "230798fe22abff459b004675bf827a4089326a296fa4165d0c2ad27688e03e0c"},
		{"remove 0x hex value", "0xfffdefefed", "fffdefefed"},
		{"remove 0x hex number", "0x16345785d8a0000", "16345785d8a0000"},
		{"remove hex without 0x", "trustwallet", "trustwallet"},
		{"remove hex number without 0x", "16345785d8a0000", "16345785d8a0000"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Remove0x(tt.input); got != tt.want {
				t.Errorf("Remove0x() = %v, want %v", got, tt.want)
			}
		})
	}
}