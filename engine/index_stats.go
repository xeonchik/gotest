package engine

type indexStats struct {
	Hits int
}

func initStats() *indexStats {
	obj := &indexStats{Hits: 0}
	return obj
}

var StatsObj = initStats()
