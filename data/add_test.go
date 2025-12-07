package data

import (
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{
			name:     "Positive numbers",
			a:        5,
			b:        3,
			expected: 8,
		},
		{
			name:     "Zero addition",
			a:        0,
			b:        0,
			expected: 0,
		},
		{
			name:     "Zero plus positive",
			a:        0,
			b:        7,
			expected: 7,
		},
		{
			name:     "Positive plus zero",
			a:        9,
			b:        0,
			expected: 9,
		},
		{
			name:     "Negative numbers",
			a:        -5,
			b:        -3,
			expected: -8,
		},
		{
			name:     "Mixed positive and negative",
			a:        10,
			b:        -4,
			expected: 6,
		},
		{
			name:     "Negative plus positive",
			a:        -7,
			b:        3,
			expected: -4,
		},
		{
			name:     "Large numbers",
			a:        1000000,
			b:        2000000,
			expected: 3000000,
		},
		{
			name:     "Minimum int32 values",
			a:        -2147483648,
			b:        -1,
			expected: -2147483649,
		},
		{
			name:     "Maximum int32 values",
			a:        2147483647,
			b:        1,
			expected: 2147483648,
		},
		{
			name:     "Opposite numbers equal zero",
			a:        42,
			b:        -42,
			expected: 0,
		},
		{
			name:     "Single digit",
			a:        1,
			b:        1,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Add(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestAddCommutative(t *testing.T) {
	testCases := []struct {
		a int
		b int
	}{
		{1, 2},
		{-1, 5},
		{0, 100},
		{-50, -25},
		{999, -999},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Commutative_%d_%d", tc.a, tc.b), func(t *testing.T) {
			result1 := Add(tc.a, tc.b)
			result2 := Add(tc.b, tc.a)
			if result1 != result2 {
				t.Errorf("Add(%d, %d) = %d but Add(%d, %d) = %d; commutative property violated", 
					tc.a, tc.b, result1, tc.b, tc.a, result2)
			}
		})
	}
}

func TestAddAssociativeProperty(t *testing.T) {
	testCases := []struct {
		a int
		b int
		c int
	}{
		{1, 2, 3},
		{-1, 5, -2},
		{0, 100, 50},
		{-50, -25, 75},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Associative_%d_%d_%d", tc.a, tc.b, tc.c), func(t *testing.T) {
			result1 := Add(Add(tc.a, tc.b), tc.c)
			result2 := Add(tc.a, Add(tc.b, tc.c))
			if result1 != result2 {
				t.Errorf("Add(Add(%d, %d), %d) = %d but Add(%d, Add(%d, %d)) = %d; associative property violated", 
					tc.a, tc.b, tc.c, result1, tc.a, tc.b, tc.c, result2)
			}
		})
	}
}

func TestAddIdentityElement(t *testing.T) {
	testValues := []int{
		0, 1, -1, 42, -100, 2147483647, -2147483648,
	}

	for _, val := range testValues {
		t.Run(fmt.Sprintf("Identity_%d", val), func(t *testing.T) {
			result1 := Add(val, 0)
			result2 := Add(0, val)
			if result1 != val || result2 != val {
				t.Errorf("Add(%d, 0) = %d, Add(0, %d) = %d; expected %d", 
					val, result1, val, result2, val)
			}
		})
	}
}

// Benchmark test for performance measurement
func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add(100, 200)
	}
}

func BenchmarkAddWithNegativeNumbers(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add(-50, 75)
	}
}

func BenchmarkAddWithZero(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add(0, 12345)
	}
}