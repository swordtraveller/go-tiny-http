package http

import "testing"

func Test_getKeyValue(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		// TODO: Add test cases.
		{"Content-Length", args{input: "Content-Length:20"}, "Content-Length", "20"},
		{"host", args{input: "host: www.google.com"}, "host", "www.google.com"},
		{"User-Agent", args{input: "User-Agent:  curl/7.68.0"}, "User-Agent", "curl/7.68.0"},
		{"Accept", args{input: "Accept: */*"}, "Accept", "*/*"},
		{"Content-Type", args{input: "Content-Type: application/x-www-form-urlencoded"}, "Content-Type", "application/x-www-form-urlencoded"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getKeyValue(tt.args.input)
			if got != tt.want {
				t.Errorf("getKeyValue() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getKeyValue() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
