package manager

import (
	"reflect"
	"testing"

	"github.com/gotomicro/ego/core/econf"
)

func TestNewDataSource(t *testing.T) {
	type args struct {
		configAddr string
		watch      bool
	}
	tests := []struct {
		name    string
		args    args
		want    econf.DataSource
		want1   econf.Unmarshaller
		want2   econf.ConfigType
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := NewDataSource(tt.args.configAddr, tt.args.watch)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDataSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDataSource() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("NewDataSource() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("NewDataSource() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestRegister(t *testing.T) {
	type args struct {
		scheme  string
		creator econf.DataSource
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Register(tt.args.scheme, tt.args.creator)
		})
	}
}
