package generator

import (
	"sort"
	"time"
)

// TODO This might be more optimal instead of a hashmap sent to nutritionist.
// TODO Create tests if so.

type Pair struct {
	Key   time.Time
	Value []Entry
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Key.Before(p[j].Key) }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func sortByTime(m map[time.Time][]Entry) PairList {
	pl := make(PairList, len(m))
	i := 0
	for k, v := range m {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(pl)
	return pl
}
