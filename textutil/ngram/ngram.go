package ngram

type NGramGroup []rune

type NGram struct {
	N      int
	Groups []NGramGroup
}

func Parse(n int, raw string) *NGram {
	runeraw := []rune(raw)
	lenRuneraw := len(runeraw)
	s := make([]NGramGroup, n*n+((lenRuneraw+1)*n))
	s[0] = make(NGramGroup, n)

	x := 0
	for i := 0; i < lenRuneraw; i++ {
		for j := 0; j < n && (i+j) < lenRuneraw; j++ {
			s[x][j] = runeraw[i+j]
		}
		if i+n >= lenRuneraw {
			break
		}

		if i+1 < lenRuneraw {
			x++
			s[x] = make([]rune, n)
		}
	}

	return &NGram{
		N:      n,
		Groups: s[:x+1],
	}
}
