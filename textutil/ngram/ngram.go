package ngram

type GramGroup []rune

type Gram struct {
	N      int
	Groups []GramGroup
}

func Parse(n int, raw string) *Gram {
	runeraw := []rune(raw)
	lenRuneraw := len(runeraw)
	s := make([]GramGroup, n*n+((lenRuneraw+1)*n))
	s[0] = make(GramGroup, n)

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

	return &Gram{
		N:      n,
		Groups: s[:x+1],
	}
}
