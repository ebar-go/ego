package array

import (
	"reflect"
	"testing"
)

func TestInt2Interface(t *testing.T) {
	type args struct {
		items []int
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			name:"test",
			args: args{items:[]int{1,2,3}},
			want: []interface{}{1,2,3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int2Interface(tt.args.items); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int2Interface() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt2String(t *testing.T) {
	type args struct {
		items []int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name:"test",
			args: args{items:[]int{1,2,3}},
			want: []string{"1","2","3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int2String(tt.args.items); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int2String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUniqueInt(t *testing.T) {
	type args struct {
		items []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name:"test",
			args: args{items:[]int{1,2,3,3}},
			want: []int{1,2,3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UniqueInt(tt.args.items); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UniqueInt() = %v, want %v", got, tt.want)
			}
		})
	}
}