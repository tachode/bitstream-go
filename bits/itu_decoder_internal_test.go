package bits

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	t.Run("MergeTwoSlices", func(t *testing.T) {
		oldSlice := []int{1, 2, 3}
		newSlice := []int{4, 5, 6}

		oldValue := reflect.ValueOf(&oldSlice).Elem()
		newValue := reflect.ValueOf(newSlice)

		merged := merge(oldValue, newValue).Interface().([]int)

		assert.Equal(t, []int{4, 5, 6}, merged)
	})

	t.Run("MergeNestedSlices", func(t *testing.T) {
		oldSlice := [][]int{{1, 2}, {3, 4}}
		newSlice := [][]int{nil, {7, 8}}

		oldValue := reflect.ValueOf(&oldSlice).Elem()
		newValue := reflect.ValueOf(newSlice)

		merged := merge(oldValue, newValue).Interface().([][]int)

		assert.Equal(t, [][]int{{1, 2}, {7, 8}}, merged)
	})

	t.Run("Merge3dSlices", func(t *testing.T) {
		oldSlice := [][][]int{{{1, 2}, {3, 4}}, {{5, 6}}}
		newSlice := [][][]int{{nil, {7, 8}}, nil, nil, {nil, {9, 10}}}

		oldValue := reflect.ValueOf(&oldSlice).Elem()
		newValue := reflect.ValueOf(newSlice)

		merged := merge(oldValue, newValue).Interface().([][][]int)

		assert.Equal(t, [][][]int{{{1, 2}, {7, 8}}, {{5, 6}}, nil, {nil, {9, 10}}}, merged)
	})

	t.Run("MergeWithEmptyOldSlice", func(t *testing.T) {
		oldSlice := []int{}
		newSlice := []int{4, 5, 6}

		oldValue := reflect.ValueOf(&oldSlice).Elem()
		newValue := reflect.ValueOf(newSlice)

		merged := merge(oldValue, newValue).Interface().([]int)

		assert.Equal(t, []int{4, 5, 6}, merged)
	})

	t.Run("MergeWithEmptyNewSlice", func(t *testing.T) {
		oldSlice := []int{1, 2, 3}
		newSlice := []int{}

		oldValue := reflect.ValueOf(&oldSlice).Elem()
		newValue := reflect.ValueOf(newSlice)

		merged := merge(oldValue, newValue).Interface().([]int)

		assert.Equal(t, []int{1, 2, 3}, merged)
	})

	t.Run("MergeNonSliceValues", func(t *testing.T) {
		oldValue := reflect.ValueOf(42)
		newValue := reflect.ValueOf(100)

		merged := merge(oldValue, newValue).Interface()

		assert.Equal(t, 100, merged)
	})
}
