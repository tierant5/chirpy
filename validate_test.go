package main

import "testing"

func Test_replaceProfanity(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		msg  string
		want string
	}{
		{
			name: "test1",
			msg:  "This is a kerfuffle opinion",
			want: "This is a **** opinion",
		},
		{
			name: "test2",
			msg:  "This is a kerfuffle! opinion",
			want: "This is a kerfuffle! opinion",
		},
		{
			name: "test3",
			msg:  "This is a KERFUFFLE opinion",
			want: "This is a **** opinion",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := replaceProfanity(tt.msg)
			if got != tt.want {
				t.Errorf("replaceProfanity() = %v, want %v", got, tt.want)
			}
		})
	}
}
