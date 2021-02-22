package engine

type indexStats struct {
	Hits int
}

func InitStats() *indexStats {
	obj := &indexStats{Hits: 0}
	return obj
}

var StatsObj = InitStats()
