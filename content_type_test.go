package main

import "testing"

func TestGetFileContentType(t *testing.T) {
	type args struct {
		content []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "json",
			args: args{content: []byte(`{"message":"Hello World"}`)},
			want: "application/json",
		},
		{
			name: "xml",
			args: args{content: []byte(`<xml><message>Hello World</message></xml>`)},
			want: "application/xml",
		},
		{
			name: "yaml",
			args: args{content: []byte(`message: Hello World`)},
			want: "application/yaml",
		},
		{
			name: "text",
			args: args{content: []byte(`Hello World`)},
			want: "text/plain",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFileContentType(tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileContentType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetFileContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}
