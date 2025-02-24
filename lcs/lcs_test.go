package lcs_test

import (
	"fmt"
	"lcs"
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLCSMatch(t *testing.T) {
	threshold := float32(0.7)

	t.Run("get lcs string size", func(t *testing.T) {
		s1 := "キャノン"
		s2 := "キヤノン"
		match, _ := lcs.LCSMatch(s1, s2, threshold)
		assert.Equal(t, true, match)
	})

	t.Run("get lcs empty", func(t *testing.T) {
		s1 := "aaaaa"
		s2 := ""
		match, _ := lcs.LCSMatch(s1, s2, threshold)
		assert.Equal(t, false, match)
	})

	t.Run("get lcs different size", func(t *testing.T) {
		s1 := "axayaaaaaz"
		s2 := "bbxbybz"
		match, _ := lcs.LCSMatch(s1, s2, threshold)
		assert.Equal(t, false, match)
	})

	t.Run("get lcs match address", func(t *testing.T) {
		s1 := "麻布台ヒルズ"
		s2 := "〒106-0041東京都港区麻布台1丁目3-1麻布台ヒルズ森JPタワー 23F"
		match, _ := lcs.LCSMatch(s1, s2, threshold)
		assert.Equal(t, true, match)

		s1 = "106-0041	あ 東京都港区麻布台 麻布台ヒルズ	森JPタワー"
		s2 = "〒106-0041東京都港区麻布台1丁目3-1麻布台ヒルズ森JPタワー 23F"
		match, _ = lcs.LCSMatch(s1, s2, threshold)
		assert.Equal(t, true, match)
	})

	t.Run("get lcs match address", func(t *testing.T) {
		s1 := "セルフィスタ渋谷"
		s2 := "インドア ゴルフレッスンスタジオ渋谷"
		match, _ := lcs.LCSMatch(s1, s2, threshold)

		// 偽陽性
		// 「インドア ゴルフレッスンスタジオ渋谷」に出現する文字がキーワード塊を無視してマッチしている
		assert.Equal(t, true, match)
	})
}

func TestLCS(t *testing.T) {
	t.Run("get lcs string size", func(t *testing.T) {
		s1 := "キャノン"
		s2 := "キヤノン"
		assert.Equal(t, int16(3), lcs.Lcs(s1, s2))
	})

	t.Run("get lcs empty", func(t *testing.T) {
		s1 := "aaaaa"
		s2 := ""
		assert.Equal(t, int16(0), lcs.Lcs(s1, s2))
		assert.Equal(t, int16(0), lcs.Lcs(s2, s1))
	})

	t.Run("get lcs different size", func(t *testing.T) {
		s1 := "axayaaaaaz"
		s2 := "bbxbybz"
		assert.Equal(t, int16(3), lcs.Lcs(s1, s2))
		assert.Equal(t, int16(3), lcs.Lcs(s2, s1))
	})
}

func TestSmithWaterman(t *testing.T) {
	l := lcs.LocalAlignment{
		MatchScore:   1,
		UnmatchScore: 1,
		GapPenarty:   1,
	}

	t.Run("same string", func(t *testing.T) {
		s1 := "京都駅"
		s2 := "京都駅"
		lcs, maxLcs := lcs.SmithWaterman(s1, s2, l)
		assert.Equal(t, int16(3), maxLcs)
		assert.Equal(t, int16(3), lcs)
	})

	t.Run("prefix string", func(t *testing.T) {
		s1 := "京都駅"
		s2 := "JR京都駅"
		lcs, maxLcs := lcs.SmithWaterman(s1, s2, l)
		assert.Equal(t, int16(3), maxLcs)
		assert.Equal(t, int16(3), lcs)
	})

	t.Run("prefix string", func(t *testing.T) {
		s1 := "京都駅"
		s2 := "京都駅西"
		lcs, maxLcs := lcs.SmithWaterman(s1, s2, l)
		assert.Equal(t, int16(3), maxLcs)
		assert.Equal(t, int16(2), lcs)
	})

	t.Run("prefix string", func(t *testing.T) {
		s1 := "京都駅"
		s2 := "梅小路京都西駅"
		lcs, maxLcs := lcs.SmithWaterman(s1, s2, l)
		assert.Equal(t, int16(2), maxLcs)
		assert.Equal(t, int16(2), lcs)
	})

	t.Run("prefix string", func(t *testing.T) {
		s1 := "京都駅"
		s2 := "レグゼスタ京都駅西"
		lcs, maxLcs := lcs.SmithWaterman(s1, s2, l)
		assert.Equal(t, int16(3), maxLcs)
		assert.Equal(t, int16(2), lcs)
	})

	t.Run("prefix string", func(t *testing.T) {
		s1 := "京都駅"
		s2 := "京都駅西ビル"
		lcs, maxLcs := lcs.SmithWaterman(s1, s2, l)
		assert.Equal(t, int16(3), maxLcs)
		assert.Equal(t, int16(0), lcs)
	})
}

func TestPrefixContains(t *testing.T) {
	t.Run("prefix string", func(t *testing.T) {
		s1 := "渋谷駅"
		s2 := "渋谷駅前病院"

		match, lcs := lcs.PrefixContains(s1, s2)

		assert.Equal(t, true, match)
		assert.Equal(t, lcs, float32(0.5))
	})

	t.Run("suffix string", func(t *testing.T) {
		s1 := "渋谷駅"
		s2 := "パークハウス渋谷駅前"

		match, lcs := lcs.PrefixContains(s1, s2)

		assert.Equal(t, false, match)
		assert.Equal(t, lcs, float32(0))
	})

	t.Run("prefix string", func(t *testing.T) {
		s1 := "渋谷駅"
		s2 := "渋谷駅前"

		match, lcs := lcs.PrefixContains(s1, s2)

		assert.Equal(t, true, match)
		assert.Equal(t, lcs, float32(0.75))
	})
}

func TestContainsEvaluate(t *testing.T) {
	t.Run("same string", func(t *testing.T) {
		s1 := "渋谷駅"
		s2 := "渋谷駅"

		point := lcs.ContainsEvaluate(s1, s2)

		assert.Equal(t, 0, point)
	})

	t.Run("prefix string", func(t *testing.T) {
		s1 := "渋谷駅"
		s2 := "渋谷駅病院"

		point := lcs.ContainsEvaluate(s1, s2)

		assert.Equal(t, 4, point)
	})

	t.Run("suffix string", func(t *testing.T) {
		s1 := "渋谷駅"
		s2 := "JR渋谷駅"

		point := lcs.ContainsEvaluate(s1, s2)

		assert.Equal(t, 2, point)
	})

	t.Run("not conatinas string", func(t *testing.T) {
		s1 := "あああ"
		s2 := "JR渋谷駅"

		point := lcs.ContainsEvaluate(s1, s2)

		assert.Equal(t, math.MaxInt, point)
	})

	t.Run("prefix string", func(t *testing.T) {
		s1 := "仙川駅"
		s2 := "ファミリーマート仙川駅西店"

		fmt.Println(strings.Index(s2, s1))

		point := lcs.ContainsEvaluate(s1, s2)

		assert.Equal(t, 12, point)
	})
}

func BenchmarkLCS(b *testing.B) {
	s1 := "キャノン"
	s2 := "キヤノン"

	b.Run("get lcs different size", func(bSub *testing.B) {
		bSub.ResetTimer()
		for i := 0; i < bSub.N; i++ {
			lcs.Lcs(s1, s2)
		}
	})

	b.Run("contains different size", func(bSub *testing.B) {
		bSub.ResetTimer()
		for i := 0; i < bSub.N; i++ {
			strings.Contains(s1, s2)
		}
	})

	s1 = "麻布台ヒルズ"
	s2 = "〒106-0041東京都港区麻布台1丁目3-1麻布台ヒルズ森JPタワー 23F"

	b.Run("get lcs address different size", func(bSub *testing.B) {
		bSub.ResetTimer()
		for i := 0; i < bSub.N; i++ {
			lcs.Lcs(s1, s2)
		}
	})

	b.Run("contains address different size", func(bSub *testing.B) {
		bSub.ResetTimer()
		for i := 0; i < bSub.N; i++ {
			strings.Contains(s1, s2)
		}
	})

	s1 = "YbS2WJYnNAmkCukdNO6gcRhhucImjB3ptoPrjmpxQknb22CxgHafHfUU0WGxN5tSab5HTlsaqEjBoPKLr7Gtudiwb7qxIlfEDWRFCDldYAykH8RLEvlcu7EqntKshBH4lCpuLNs6AxWJUXBoBQ0uQN6qoYFj3jn7OI6unE33okuwuAHioCQlw7SrWta7U8q7KKmihjKq03NPQvdv4s3au8LvFCHt1BAwA3KIlHD2406xIxRtvIeb40bAp1huQwyhMYbd3GpURkaIEMEO8mDoQD1nXgnVAnqX7RfRvqYIel7IRwsHsJiNcxuMJgc0xY6qFjPtqWD156nSAlV4JlhI8EfpnidmnkfpshNN7yXsQdVIxLgfTVdKi2U0DEVGJjnqwpYf6dUYufgfWUSUhCO8bF8ErA4vEcy6rb7e3aC3WbFlU2LKKKuSQJVjQfMlVtFPipIsPig66CApVihmd68SMI73VcXcqtEypqDnrLfJlUcjV6okbTNG"
	s2 = "BDugtL6sEGuauPamYIWKGf1pL73XPXlpUxoI6Nv6YQLAHh4l2THkfUI3SPqkKJNWoOveK25MPDLpdlAh5SYImy1hw5Pv3CNPCOqoWNr5BAEFeUwl0M12rs5LqH4i1jHFoEQyMt6Ng65YubGfPPtkF6KSJiSibRc64TejUrM8cwUTfhJcqv4cQTDF7ALBmmWsfKkJBJiBWdcdgNqxC8M5J6wIVRUTMpRjrKMhbQWJMT3yJJOWHpEeSmBPt2FWRAUUpeEs1dOQwEwKDpQ6RfuI85HTT2104bhlLFlRG8yexUhxyfnhi3GSjfmTv8DVNeYETWdx44tBsi0U3OTH1ntW4EgJNbN1GDSChIPxOMLgDjCphaEVqBpKkF3A74ExFm5FqvTDvlqx5dEK2QCXeXBkO8TwqTsQU7QYoy8SdRGr7RLvobK5LQkGToWf5VlCD0PYBaPjWlak6lMlnogKeJW008OFpKpr4b1JyE0OKQxTeX60Sj4S6e0G"

	b.Run("get lcs long different size", func(bSub *testing.B) {
		bSub.ResetTimer()
		for i := 0; i < bSub.N; i++ {
			lcs.Lcs(s1, s2)
		}
	})

	b.Run("contains long different size", func(bSub *testing.B) {
		bSub.ResetTimer()
		for i := 0; i < bSub.N; i++ {
			strings.Contains(s1, s2)
		}
	})
}
