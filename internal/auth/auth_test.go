package auth_test

import (
	"testing"

	"github.com/tierant5/chirpy/internal/auth"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		password string
		want     string
		wantErr  bool
	}{
		{
			name:     "test1",
			password: "Pa$$word!",
			want:     "",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := auth.HashPassword(tt.password)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("HashPassword() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("HashPassword() succeeded unexpectedly")
			}
			if got == tt.password {
				t.Errorf("HashPassword() = %v, want %v", got, tt.want)
			}
			err := auth.CheckPasswordHash(tt.password, got)
			if err != nil {
				t.Errorf("CheckPasswordHash() = %v, failed unexpectedly", err)
			}
			err = auth.CheckPasswordHash("wrongPassword", got)
			if err == nil {
				t.Errorf("CheckPasswordHash() = %v, passed with incorrect password", err)
			}
			err = auth.CheckPasswordHash(tt.password, "")
			if err == nil {
				t.Errorf("CheckPasswordHash() = %v, passed with incorrect hash", err)
			}
		})
	}
}
