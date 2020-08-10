package redisid

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func KnuthDurstenfeldShuffle(array []string) []string {
	for i := len(array) - 1; i >= 0; i-- {
		p := RandInt64(0, int64(i))
		a := array[i]
		array[i] = array[p]
		array[p] = a
	}
	return array
}

func RandInt64(min, max int64) int64 {
	randSource := rand.NewSource(time.Now().UnixNano())
	randCustom := rand.New(randSource)
	if min >= max || max == 0 {
		return max
	}
	return randCustom.Int63n(max-min) + min
}

func TestShuffle(t *testing.T) {
	//input :"ABCDEFGHIJKLMNOPRSTUVWXYZ0123456789"
	input := "ABCDEFGHJKLMNPRSTUVWXY0123456789"
	a := []rune(input)
	var inputArray []string
	for _, r := range a {
		inputArray = append(inputArray, string(r))
	}
	output := KnuthDurstenfeldShuffle(inputArray)
	assert.EqualValues(t, len(input), len(output))
	output = KnuthDurstenfeldShuffle(inputArray)
	assert.EqualValues(t, len(input), len(output))
	output = KnuthDurstenfeldShuffle(inputArray)
	assert.EqualValues(t, len(input), len(output))
	output = KnuthDurstenfeldShuffle(inputArray)
	assert.EqualValues(t, len(input), len(output))
	output = KnuthDurstenfeldShuffle(inputArray)
	assert.EqualValues(t, len(input), len(output))
	output = KnuthDurstenfeldShuffle(inputArray)
	assert.EqualValues(t, len(input), len(output))
	fmt.Println(strings.Join(output, ""))
}

func TestInt2String2(t *testing.T) {

	//test_944_17_14_13 52
	//test_944_17_14_14 52
	//test_944_17_14_15 52
	//SH052 3SYEKMV 944
	//SH052 3SYEKMV 944
	//SH052 3SYE89D
	//1597079654
	result := Int2String(1597051040)
	fmt.Println(result)
}
func TestInt2String(t *testing.T) {
	check := make(map[int64]string)
	check[0] = "D"
	check[1] = "3"
	check[2] = "E"
	check[3] = "K"
	check[4] = "8"
	check[5] = "9"
	check[30] = "L"
	check[31] = "M"
	check[32] = "3D"
	check[33] = "33"
	check[34] = "3E"
	check[35] = "3K"
	for k, v := range check {
		result := Int2String(k)
		assert.EqualValues(t, v, result)
	}
}

func TestHexEncodeToString(t *testing.T) {
	input := "SILLY_HAT_XU"
	expected := "53494C4C595F4841545F5855"
	output := hexEncodeToString(input)
	assert.EqualValues(t, expected, output)
}

func TestGetLifeCycleNumber(t *testing.T) {
	assert.EqualValues(t, 1, getLifeCycleNumber(Second))
	assert.EqualValues(t, 60, getLifeCycleNumber(Minute))
	assert.EqualValues(t, 60*60, getLifeCycleNumber(Hour))
}
