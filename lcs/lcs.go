package lcs

import (
	"math"
	"strings"
	"unicode/utf8"
)

// Lcs 最長共通部分列のサイズを返す O(NM)
// https://www.cs.t-kougei.ac.jp/SSys/LCS.htm
func Lcs(s, t string) int16 {
	runeS := []rune(s)
	runeT := []rune(t)

	n, m := len(runeS), len(runeT)
	dp := make([][]int16, n+1)
	for i := 0; i < len(dp); i++ {
		dp[i] = make([]int16, m+1)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			// 一致した
			// SiとTjを採用する前の状態からの遷移
			if runeS[i] == runeT[j] {
				dp[i+1][j+1] = max(dp[i+1][j+1], dp[i][j]+1)
			} else {
				// 一致していない
				// dp[i][j+1]: Siを採用する前の状態からの遷移
				// dp[i+1][j]: Tjを採用する前の状態からの遷移
				dp[i+1][j+1] = max(dp[i][j+1], dp[i+1][j])
			}
		}
	}

	return dp[n][m]
}

type LocalAlignment struct {
	MatchScore   int16
	UnmatchScore int16
	GapPenarty   int16
}

func SmithWaterman(s, t string, a LocalAlignment) (lcs, maxLcs int16) {
	runeS := []rune(s)
	runeT := []rune(t)

	n, m := len(runeS), len(runeT)
	dp := make([][]int16, n+1)
	for i := 0; i < len(dp); i++ {
		dp[i] = make([]int16, m+1)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {

			var diagonalScore int16
			// 一致した
			if runeS[i] == runeT[j] {
				diagonalScore = dp[i][j] + a.MatchScore
			} else {
				// 一致していない
				diagonalScore = dp[i][j] - a.UnmatchScore
			}

			dp[i+1][j+1] = max(
				0,
				diagonalScore,
				dp[i][j+1]-a.GapPenarty, // 縦方向の遷移
				dp[i+1][j]-a.GapPenarty, // 横方向の遷移
			)

			//　部分一致のLCSを評価するため最大のLCSを取得する
			maxLcs = max(maxLcs, dp[i+1][j+1])
		}
	}

	return dp[n][m], maxLcs
}

//nolint:gochecknoglobals
var alignment = LocalAlignment{
	MatchScore:   1,
	UnmatchScore: 1,
	GapPenarty:   1,
}

func PrefixContains(substr, s string) (result bool, match float32) {
	result = strings.HasPrefix(s, substr)
	if result {
		prefixSize := float32(len([]rune(substr)))
		baseSize := float32(len([]rune(s)))

		match = prefixSize / baseSize

	}

	return
}

// ContainsEvaluate 部分一致率を評価する. 完全一致 => プレフィックス一致 => 部分一致の順にマッチ度を返却する
func ContainsEvaluate(substr, s string) int {
	result := strings.Contains(s, substr)
	if !result {
		return math.MaxInt
	}

	r1 := []rune(s)
	r2 := []rune(substr)

	idx := strings.Index(s, substr)

	prefixLen := utf8.RuneCountInString(s[:idx]) // 先頭の文字数

	suffixLen := len(r1) - (prefixLen + len(r2)) // 末尾の文字数

	return prefixLen + suffixLen*2
}

func SmithWatermanMatch(substr, s string, threshold float32) (bool, float32) {
	runeS := []rune(s)

	// 空文字はスコアに無関係なので削除
	s = strings.ReplaceAll(s, " ", "")
	substr = strings.ReplaceAll(substr, " ", "")

	_, maxLCS := SmithWaterman(substr, s, alignment)

	match := float32(maxLCS) / float32(len(runeS))

	return threshold <= match, match
}

func LCSMatch(substr, s string, threshold float32) (bool, float32) {
	runeS := []rune(substr)

	size := Lcs(substr, s)

	match := float32(size) / float32(len(runeS))

	return threshold <= match, match
}
