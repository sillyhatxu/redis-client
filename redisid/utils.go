package redisid

import (
	"encoding/hex"
	"strings"
)

const BaseString = "D3EK89VRPCNWX1US76GY245A0TJFBHLM"
const BaseStringLength = int64(len(BaseString))

func Int2String(seq int64) (shortURL string) {
	var charSeq []rune
	if seq != 0 {
		for seq != 0 {
			mod := seq % BaseStringLength
			div := seq / BaseStringLength
			charSeq = append(charSeq, rune(BaseString[mod]))
			seq = div
		}
	} else {
		charSeq = append(charSeq, rune(BaseString[seq]))
	}

	tmpShortURL := string(charSeq)
	shortURL = reverse(tmpShortURL)
	return
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func getLifeCycleNumber(lifeCycle LifeCycleType) int64 {
	if lifeCycle == Minute {
		return 60
	} else if lifeCycle == Hour {
		return 60 * 60
	} else {
		return 1
	}
}

func hexEncodeToString(s string) string {
	return strings.ToUpper(hex.EncodeToString([]byte(s)))
}
