package random

import (
	"testing"
)

func TestRandomizeFileName(t *testing.T) {
	fileName := "我喜欢你.jpg"
	want := "3b2b1b2b"
	if got := RandomizeFileName(fileName); got != want {
		t.Errorf("RandomizeFileName() = %v, want %v", got, want)
	}
}
