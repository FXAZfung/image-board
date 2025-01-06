package utils

import (
	"testing"
)

func TestGenerateShortLink(t *testing.T) {
	input := "我喜欢你.jpg"
	want := "3b2b1b2b"
	if got := GenerateShortLink(input); got != want {
		t.Errorf("GenerateShortLink() = %v, want %v", got, want)
	}
}
