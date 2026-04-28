package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDaysToMask(t *testing.T) {
	mask := DaysToMask([]int{1, 2, 7})
	assert.True(t, mask&Monday != 0)
	assert.True(t, mask&Tuesday != 0)
	assert.True(t, mask&Sunday != 0)
	assert.False(t, mask&Wednesday != 0)
}

func TestMaskToDays(t *testing.T) {
	mask := Monday | Tuesday | Sunday
	days := MaskToDays(mask)
	assert.Contains(t, days, 1)
	assert.Contains(t, days, 2)
	assert.Contains(t, days, 7)
	assert.Len(t, days, 3)
}

func TestWeekdayToMask(t *testing.T) {
	assert.Equal(t, Monday, WeekdayToMask(time.Monday))
	assert.Equal(t, Tuesday, WeekdayToMask(time.Tuesday))
	assert.Equal(t, Wednesday, WeekdayToMask(time.Wednesday))
	assert.Equal(t, Thursday, WeekdayToMask(time.Thursday))
	assert.Equal(t, Friday, WeekdayToMask(time.Friday))
	assert.Equal(t, Saturday, WeekdayToMask(time.Saturday))
	assert.Equal(t, Sunday, WeekdayToMask(time.Sunday))
}
