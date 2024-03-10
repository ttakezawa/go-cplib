package cplib

import "testing"

func TestLCM(t *testing.T) {
	type args struct {
		xs []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "example1",
			args: args{[]int{2}},
			want: 2,
		},
		{
			name: "example2",
			args: args{[]int{2, 3}},
			want: 6,
		},
		{
			name: "example3",
			args: args{[]int{2, 3, 5}},
			want: 30,
		},
		{
			name: "example4",
			args: args{[]int{2, 3, 5, 7}},
			want: 210,
		},
		{
			name: "example5",
			args: args{[]int{2, 3, 3, 7, 4, 3, 7}},
			want: 84,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LCM(tt.args.xs...); got != tt.want {
				t.Errorf("LCM() = %v, want %v", got, tt.want)
			}
		})
	}
}
