package xstring

import "testing"

func TestFunctionName(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FunctionName(tt.args.i); got != tt.want {
				t.Errorf("FunctionName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestObjectName(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ObjectName(tt.args.i); got != tt.want {
				t.Errorf("ObjectName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCallerName(t *testing.T) {
	type args struct {
		skip int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CallerName(tt.args.skip); got != tt.want {
				t.Errorf("CallerName() = %v, want %v", got, tt.want)
			}
		})
	}
}
