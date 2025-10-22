package levenshtein

func LevenshteinDistance(a string, b string) int {
	n := len(a)
	m := len(b)

	p := make([]int, m+1)
	c := make([]int, m+1)

	for j := range m {
		p[j] = j
	}
	for i := range n {
		c[0] = i + 1

		for j := range m {
			sc := 0
			if a[i] != b[j] {
				sc = 1
			}
			c[j+1] = min(
				p[j+1]+1, // deletion
				c[j]+1,   // insertion
				p[j]+sc,  // substitution
			)
		}
		copy(p, c)
	}
	return p[m]
}
