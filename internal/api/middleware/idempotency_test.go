package middleware

import (
	"net/http"
	"reflect"
	"testing"
)

func TestKeyAndURL(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "Should create key properly",
			args: args{
				r: &http.Request{
					RequestURI: "https://some-url.com",
					Header: http.Header{
						"Idempotency-Id": []string{"some-id"},
					},
				},
			},
			want: struct {
				Key string
				URL string
			}{
				Key: "some-id",
				URL: "https://some-url.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KeyAndURL(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeyAndURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
