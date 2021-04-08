package main

import "testing"

func Test_yamlUnmarshal(t *testing.T) {
	type args struct {
		buf  []byte
		dest interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := yaml_Unmarshal(tt.args.buf, tt.args.dest); (err != nil) != tt.wantErr {
				t.Errorf("yamlUnmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
