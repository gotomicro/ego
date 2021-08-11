package xtime

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDuration(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "",
			args: args{str: "1s"},
			want: time.Second,
		},
		{
			name: "",
			args: args{str: "2m"},
			want: 2 * time.Minute,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Duration(tt.args.str); got != tt.want {
				t.Errorf("Duration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseInLocation(t *testing.T) {
	type args struct {
		layout string
		value  string
		tzEnv  string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "NOT set TZ && time value NOT include zone abbreviations",
			args: args{layout: "2006-01-02 15:04:05", value: "2019-12-31 00:00:00"},
			want: 1577721600 + 8*3600,
		},
		{
			name: "NOT set TZ && time value include zone abbreviations",
			args: args{layout: "2006-01-02T15:04:05Z07:00", value: "2019-12-31T00:00:00+08:00"},
			want: 1577721600,
		},
		{
			name: "set TZ && time value NOT include zone abbreviations",
			args: args{layout: "2006-01-02 15:04:05", value: "2019-12-31 00:00:00", tzEnv: "Asia/Shanghai"},
			want: 1577721600,
		},
		{
			name: "set TZ && time value include zone abbreviations",
			args: args{layout: "2006-01-02T15:04:05Z07:00", value: "2019-12-31T00:00:00+08:00", tzEnv: "Asia/Shanghai"},
			want: 1577721600,
		},
		{
			name: "set TZ && time value include zone abbreviations && value zone abbreviations against TZ", // will ignore tzEnv
			args: args{layout: "2006-01-02T15:04:05Z07:00", value: "2019-12-31T00:00:00+00:00", tzEnv: "Asia/Shanghai"},
			want: 1577721600 + 8*3600,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.Setenv("TZ", tt.args.tzEnv)
			defer os.Unsetenv("TZ")
			assert.NoError(t, err)

			got, err := ParseInLocation(tt.args.layout, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Unix(), tt.want) {
				t.Errorf("Parse() got = %v, want %v", got.Unix(), tt.want)
			}
		})
	}
}
