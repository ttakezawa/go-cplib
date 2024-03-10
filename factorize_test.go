// Originated from https://qiita.com/Kiri8128/items/eca965fe86ea5f4cbb98
// Verify https://algo-method.com/tasks/553
package cplib

import (
	"reflect"
	"testing"
)

func Test_factorize(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want map[int]int
	}{
		{
			name: "example1",
			args: args{12},
			want: map[int]int{2: 2, 3: 1},
		},
		{
			name: "example2",
			args: args{341550054645379},
			want: map[int]int{341550054645379: 1},
		},
		{
			name: "example3",
			args: args{100},
			want: map[int]int{2: 2, 5: 2},
		},
		{
			name: "example4",
			args: args{347484690041206937},
			want: map[int]int{381727069: 1, 910296173: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Factorize(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("factorize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_is_prime(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "example1",
			args: args{2},
			want: true,
		},
		{
			name: "example2",
			args: args{341550054645379},
			want: true,
		},
		{
			name: "example3",
			args: args{347484690041206937},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPrime(tt.args.n); got != tt.want {
				t.Errorf("is_prime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_find_factor_rho(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "example1",
			args: args{347484690041206937},
			want: 381727069,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := _findFactorRho(tt.args.n); got != tt.want {
				t.Errorf("find_factor_rho() = %v, want %v", got, tt.want)
			}
		})
	}
}
