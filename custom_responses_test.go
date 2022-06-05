package main

import (
	"testing"
)

func TestGetCustomResponseFolder(t *testing.T) {
	type args struct {
		folderName string
	}
	tests := []struct {
		name    string
		args    args
		wantCrf *CustomResponseFolder
		wantErr bool
	}{
		{
			name: "custom_responses",
			args: args{folderName: "custom_responses"},
			wantCrf: &CustomResponseFolder{
				Path: "custom_responses",
				RoutesData: map[string]CustomResponseData{
					"/api": {
						Path:        "/api",
						Content:     []byte(`{"message":"Hello World"}`),
						Method:      "GET",
						StatusCode:  200,
						ContentType: "application/json",
					},
				},				
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCrf, err := GetCustomResponseFolder(tt.args.folderName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCustomResponseFolder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(gotCrf.Path) == 0 {
				t.Errorf("GetCustomResponseFolder() gotCrf.Path = %v, want %v", gotCrf.Path, tt.wantCrf.Path)
			}
		})
	}
}
