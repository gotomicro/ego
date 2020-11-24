package xmap

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestMergeStringMap(t *testing.T) {
	type args struct {
		dest map[string]interface{}
		src  map[string]interface{}
		tar  map[string]interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "二维测试",
			args: args{
				dest: map[string]interface{}{
					"2w": map[string]interface{}{
						"test":  "2wtd",
						"test1": "2wtd1",
					},
					"2wa": map[string]interface{}{
						"test":  "2wtd",
						"test1": "2wtd1",
					},
					"2wi": map[interface{}]interface{}{
						"test":  "2wtd",
						"test1": "2wtd1",
					},
				},
				src: map[string]interface{}{
					"2w": map[string]interface{}{
						"test":  "2wtds",
						"test1": "2wtd1s",
					},
					"2wb": map[string]interface{}{
						"test":  "2wtds",
						"test1": "2wtd1s",
					},
					"2wi": map[interface{}]interface{}{
						"test":  "2wtds",
						"test1": "2wtd1s",
					},
				},
				tar: map[string]interface{}{
					"2w": map[string]interface{}{
						"test":  "2wtds",
						"test1": "2wtd1s",
					},
					"2wb": map[string]interface{}{
						"test":  "2wtds",
						"test1": "2wtd1s",
					},
					"2wa": map[string]interface{}{
						"test":  "2wtd",
						"test1": "2wtd1",
					},
					"2wi": map[string]interface{}{
						"test":  "2wtds",
						"test1": "2wtd1s",
					},
				},
			},
		},
		{
			name: "一维测试",
			args: args{
				dest: map[string]interface{}{
					"1w":  "tt",
					"1wa": "mq",
				},
				src: map[string]interface{}{
					"1w":  "tts",
					"1wb": "bq",
				},
				tar: map[string]interface{}{
					"1w":  "tts",
					"1wa": "mq",
					"1wb": "bq",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MergeStringMap(tt.args.dest, tt.args.src)
			if !reflect.DeepEqual(tt.args.dest, tt.args.tar) {
				spew.Dump(tt.args.dest)
				t.FailNow()
			}
		})
	}
}
