package client

import (
	"reflect"
	"testing"
)

func mapHash(hash interface{}) RpcRequest {
	array := []interface{}{hash}
	return RpcRequest{
		Method: "GetTransaction",
		Params: array,
	}
}

func Test_makeRequests(t *testing.T) {
	type args struct {
		hashes   []interface{}
		perGroup int
	}
	tests := []struct {
		name string
		args args
		want []RpcRequests
	}{
		{
			name: "test group size 1",
			args: args{
				hashes: []interface{}{
					"0x1", "0x2", "0x3",
				},
				perGroup: 1,
			},
			want: []RpcRequests{
				{
					&RpcRequest{
						Method: "GetTransaction",
						Params: []interface{}{"0x1"},
					},
				},
				{
					&RpcRequest{
						Method: "GetTransaction",
						Params: []interface{}{"0x2"},
					},
				},
				{
					&RpcRequest{
						Method: "GetTransaction",
						Params: []interface{}{"0x3"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MakeBatchRequests(tt.args.hashes, tt.args.perGroup, mapHash); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makeBatchRequests() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_makeBatches(t *testing.T) {
	type args struct {
		hashes    []interface{}
		batchSize int
	}
	tests := []struct {
		name        string
		args        args
		wantBatches [][]interface{}
	}{
		{
			name: "Test batch size 4",
			args: args{
				hashes: []interface{}{
					"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11",
				},
				batchSize: 4,
			},
			wantBatches: [][]interface{}{
				{"1", "2", "3", "4"},
				{"5", "6", "7", "8"},
				{"9", "10", "11"},
			},
		},
		{
			name: "Test batch size 10",
			args: args{
				hashes: []interface{}{
					"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11",
				},
				batchSize: 10,
			},
			wantBatches: [][]interface{}{
				{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
				{"11"},
			},
		},
		{
			name: "Test batch size 11",
			args: args{
				hashes: []interface{}{
					"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11",
				},
				batchSize: 11,
			},
			wantBatches: [][]interface{}{
				{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotBatches := MakeBatches(tt.args.hashes, tt.args.batchSize); !reflect.DeepEqual(gotBatches, tt.wantBatches) {
				t.Errorf("makeBatches() = %v, want %v", gotBatches, tt.wantBatches)
			}
		})
	}
}
