package dwa

import (
	"fmt"
	"testing"
	"time"
)

func ExampleNewDWA() {
	d := NewDWA(5, 10*time.Millisecond)
	// only 5 items fit in sample
	d.Add(1, 2, 3, 4, 5, 6, 7)
	// average of 3, 4, 5, 6, 7
	fmt.Println(d.Value())
	time.Sleep(15 * time.Millisecond)
	// first element decayed
	// average of 0, 4, 5, 6, 7
	fmt.Println(d.Value())
	// Output:
	//
	// 5
	// 4.4
}

func TestAdd(t *testing.T) {
	d := NewDWA(5, 0)
	if d.Value() != 0 {
		t.Fatal("expected value is 0, got", d.Value())
	}
	d.Add(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	// sum(6,7,8,9,10)/5
	if d.Value() != float64(8) {
		t.Fatal("expected value is 8, got", d.Value())
	}
}

func TestValueDecay1(t *testing.T) {
	d := NewDWA(5, 10*time.Millisecond)
	d.Add(0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	time.Sleep(15 * time.Millisecond)
	if v := d.Value(); v != 6.8 {
		t.Fatalf("want 6.8, got %v, values are %v", v, d.values)
	}
	time.Sleep(22 * time.Millisecond)
	if v := d.Value(); v != 3.8 {
		t.Fatalf("want 3.8, got %v, values are %v", v, d.values)
	}
	time.Sleep(60 * time.Millisecond)
	if v := d.Value(); v != 0 {
		t.Fatalf("want 0, got %v, values are %v", v, d.values)
	}
}

func TestAddInvalidCall1(t *testing.T) {
	defer func() {
		if x := recover(); x == nil {
			t.Fatal("call with invalid params should have panicked, but it does not")
		}
	}()
	NewDWA(0, time.Second)
}

func TestAddInvalidCall2(t *testing.T) {
	defer func() {
		if x := recover(); x == nil {
			t.Fatal("call with invalid params should have panicked, but it does not")
		}
	}()
	NewDWA(3, -time.Second)
}

func BenchmarkAdd(b *testing.B) {
	d := NewDWA(100, time.Second)
	for i := 0; i < b.N; i++ {
		d.Add(42)
	}
}

func BenchmarkValue100(b *testing.B) {
	d := populatedDWA(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Value()
	}
}

func BenchmarkValue1000(b *testing.B) {
	d := populatedDWA(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Value()
	}
}

func BenchmarkValue10000(b *testing.B) {
	d := populatedDWA(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Value()
	}
}
func populatedDWA(n int) *DWA {
	d := NewDWA(n, 0)
	samples := make([]int64, n)
	for i := 0; i < n; i++ {
		samples[i] = int64(i)
	}
	d.Add(samples...)
	return d
}
