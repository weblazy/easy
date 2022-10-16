package fmetric

import "testing"

func TestSplitGrpcMethodName(t *testing.T) {
	type args struct {
		fullMethodName string
	}

	tests := []struct {
		name        string
		args        args
		wantService string
		wantMethod  string
	}{
		{
			"normal",
			args{fullMethodName: "/channel.quote.v1.QuoteAPI/QuotationGenerate"},
			"channel.quote.v1.QuoteAPI",
			"QuotationGenerate",
		},
		{
			"invalid",
			args{fullMethodName: "/channel.quote.v1.QuoteAPI"},
			unknown,
			unknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotService, gotMethod := SplitGrpcMethodName(tt.args.fullMethodName)
			if gotService != tt.wantService {
				t.Errorf("SplitGrpcMethodName() gotService = %v, want %v", gotService, tt.wantService)
			}
			if gotMethod != tt.wantMethod {
				t.Errorf("SplitGrpcMethodName() gotMethod = %v, want %v", gotMethod, tt.wantMethod)
			}
		})
	}
}
