package util

import (
	"strconv"
	"testing"
)

func TestNormalizeAndValidatePhone(t *testing.T) {
	tests := []struct {
		phone string
		want  bool
	}{
		{"+77086358906", true},
		{"87086358906", true},
		{"77086358906", true},
		{"770863589067", true},
	}
	for ttI, tt := range tests {
		t.Run(strconv.Itoa(ttI), func(t *testing.T) {
			if got := NormalizeAndValidatePhone(&tt.phone); got != tt.want {
				t.Errorf("NormalizeAndValidatePhone() = %v, want %v", got, tt.want)
			}
		})
	}
}
