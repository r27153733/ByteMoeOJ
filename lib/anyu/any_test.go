package anyu

import (
	"reflect"
	"testing"
)

func TestAnyToPointer(t *testing.T) {
	t.Run("int value", func(t *testing.T) {
		var input any = 42
		ptr := AnyToPointer[int](input)
		if *ptr != 42 {
			t.Errorf("AnyToPointer() = %v, want %v", *ptr, 42)
		}
	})

	t.Run("string value", func(t *testing.T) {
		var input any = "hello"
		ptr := AnyToPointer[string](input)
		if *ptr != "hello" {
			t.Errorf("AnyToPointer() = %v, want %v", *ptr, "hello")
		}
	})

	t.Run("struct value", func(t *testing.T) {
		type TestStruct struct{ Name string }
		expected := TestStruct{"test"}
		var input any = expected
		ptr := AnyToPointer[TestStruct](input)
		if !reflect.DeepEqual(*ptr, expected) {
			t.Errorf("AnyToPointer() = %v, want %v", *ptr, expected)
		}
	})
}

func TestPointerToValueAny(t *testing.T) {
	t.Run("primitive types", func(t *testing.T) {
		// Test with int
		intVal := 42
		intAny := PointerToValueAny(&intVal)
		if intAny.(int) != intVal {
			t.Errorf("PointerToValueAny() int = %v, want %v", intAny, intVal)
		}

		// Test with string
		strVal := "hello"
		strAny := PointerToValueAny(&strVal)
		if strAny.(string) != strVal {
			t.Errorf("PointerToValueAny() string = %v, want %v", strAny, strVal)
		}

		// Test with bool
		boolVal := true
		boolAny := PointerToValueAny(&boolVal)
		if boolAny.(bool) != boolVal {
			t.Errorf("PointerToValueAny() bool = %v, want %v", boolAny, boolVal)
		}
	})

	t.Run("struct types", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		person := Person{Name: "Alice", Age: 30}
		personAny := PointerToValueAny(&person)

		if !reflect.DeepEqual(personAny.(Person), person) {
			t.Errorf("PointerToValueAny() struct = %v, want %v", personAny, person)
		}
	})

	t.Run("modify through pointer", func(t *testing.T) {
		// Create a value and convert to any
		intVal := 42
		intAny := PointerToValueAny(&intVal)

		// Modify the original value
		intVal = 100

		// Check if the any value reflects the change
		if intAny.(int) != 100 {
			t.Errorf("PointerToValueAny() did not reflect changes to original value, got %v, want %v", intAny, 100)
		}
	})

	t.Run("slice types", func(t *testing.T) {
		slice := []int{1, 2, 3}
		sliceAny := PointerToValueAny(&slice)

		if !reflect.DeepEqual(sliceAny.([]int), slice) {
			t.Errorf("PointerToValueAny() slice = %v, want %v", sliceAny, slice)
		}

		// Modify the original slice
		slice[0] = 99

		// Check if the any value reflects the change
		if sliceAny.([]int)[0] != 99 {
			t.Errorf("PointerToValueAny() did not reflect changes to slice, got %v, want %v", sliceAny.([]int)[0], 99)
		}
	})
}

func TestRoundTrip(t *testing.T) {
	t.Run("PointerToValueAny -> AnyToPointer with specific types", func(t *testing.T) {
		// Test with string
		t.Run("string", func(t *testing.T) {
			original := "test string"
			anyVal := PointerToValueAny(&original)
			ptr := AnyToPointer[string](anyVal)

			if *ptr != original {
				t.Errorf("Round trip failed, got %v, want %v", *ptr, original)
			}

			// Modify through pointer
			*ptr = "modified"

			// Check if change is reflected
			modifiedAny := anyVal.(string)
			if modifiedAny != "modified" {
				t.Errorf("Changes not reflected, any value is %v, want %v", modifiedAny, "modified")
			}
		})

		// Test with int
		t.Run("int", func(t *testing.T) {
			original := 42
			anyVal := PointerToValueAny(&original)
			ptr := AnyToPointer[int](anyVal)

			if *ptr != original {
				t.Errorf("Round trip failed, got %v, want %v", *ptr, original)
			}

			// Modify through pointer
			*ptr = 100

			// Check if change is reflected
			modifiedAny := anyVal.(int)
			if modifiedAny != 100 {
				t.Errorf("Changes not reflected, any value is %v, want %v", modifiedAny, 100)
			}
		})

		// Test with struct
		t.Run("struct", func(t *testing.T) {
			type Person struct {
				Name string
				Age  int
			}

			original := Person{Name: "Alice", Age: 30}
			anyVal := PointerToValueAny(&original)
			ptr := AnyToPointer[Person](anyVal)

			if !reflect.DeepEqual(*ptr, original) {
				t.Errorf("Round trip failed, got %v, want %v", *ptr, original)
			}

			// Modify through pointer
			ptr.Age = 31

			// Check if change is reflected
			modifiedAny := anyVal.(Person)
			expected := Person{Name: "Alice", Age: 31}
			if !reflect.DeepEqual(modifiedAny, expected) {
				t.Errorf("Changes not reflected, any value is %v, want %v", modifiedAny, expected)
			}
		})
	})

	t.Run("avoid heap allocation", func(t *testing.T) {
		// This test is mostly illustrative, as we can't directly test memory allocation patterns
		// But we can verify the functionality

		type Large struct {
			Data [1024]int
		}

		large := Large{}
		for i := range large.Data {
			large.Data[i] = i
		}

		// Convert to any using our function
		anyVal := PointerToValueAny(&large)

		// Access and verify some values
		result := anyVal.(Large)
		for i := 0; i < 10; i++ {
			if result.Data[i] != i {
				t.Errorf("Data mismatch at index %d, got %v, want %v", i, result.Data[i], i)
			}
		}
	})
}

type Large struct {
	Data [64]int
}

var l = new(Large)

//go:noinline
func safePointerToValueAny[T any](t *T) (res any) {
	res = *t
	return
}

// Benchmarks to compare standard interface conversion vs our functions
func BenchmarkStandardConversion(b *testing.B) {
	b.ReportAllocs()

	large := l
	for i := range large.Data {
		large.Data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x := safePointerToValueAny(large)
		_ = x.(Large)
	}
}

func BenchmarkPointerToValueAny(b *testing.B) {
	b.ReportAllocs()

	large := l
	for i := range large.Data {
		large.Data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x := PointerToValueAny(large)
		_ = x.(Large)
	}
}
