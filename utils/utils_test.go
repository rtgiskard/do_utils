package main

import (
	"math/rand"
	"testing"
	"time"
)

func TestInSlice(t *testing.T) {

	{
		s, c, want := []bool{true, false}, true, true
		get := InSlice(s, c)
		if want != get {
			t.Errorf("InSlice(%v, %v) want %t, got %t", s, c, want, get)
		}
	}

	{
		s, c, want := []int{1, 2, 3, 4}, 4, true
		get := InSlice(s, c)
		if want != get {
			t.Errorf("InSlice(%v, %v) want %t, got %t", s, c, want, get)
		}
	}

	{
		s, c, want := []float32{1.1, 2.1, 3.0, 0.4}, 4, false
		get := InSlice(s, c)
		if want != get {
			t.Errorf("InSlice(%v, %v) want %t, got %t", s, c, want, get)
		}
	}

	{
		s, c, want := []string{"a", "ab", "cc"}, "c", false
		get := InSlice(s, c)
		if want != get {
			t.Errorf("InSlice(%v, %v) want %t, got %t", s, c, want, get)
		}
	}
}

func TestGenRanStr(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 1024; i++ {
		for s := 1; s < 5; s++ {
			if len(GenRandStr(i, s)) != i {
				t.Errorf("GenRandStr length not match: (%d, %d)", i, s)
			}
		}
	}

	t.Log(GenRandStr(20, 1))
	t.Log(GenRandStr(20, 2))
	t.Log(GenRandStr(20, 3))
	t.Log(GenRandStr(20, 4))
	t.Log(GenRandStr(20, 5))
}

func TestReprBitsLen(t *testing.T) {
	type Test struct {
		arg  uint64
		want int
	}
	tests := []Test{
		{0, 1},
		{1, 1},
		{2, 2},
		{63, 6},
		{64, 7},
		{127, 7},
	}

	for _, test := range tests {
		get := ReprBitsLen(test.arg)
		if get != test.want {
			t.Errorf("ReprBits(%d) want: %d, get: %d", test.arg, test.want, get)
		}
	}
}

func BenchmarkReprBitsLen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ReprBitsLen(uint64(123456))
	}
}

func BenchmarkGenRandStr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenRandStr(20, 5)
	}
}
