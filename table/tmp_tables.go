package table

import "godoc/engine"

type TemporaryDataSet struct {
	Records []*engine.DataRecord
	Length  uint32
}

func CreateTempTable() *TemporaryDataSet {
	return &TemporaryDataSet{}
}

func (table *TemporaryDataSet) Add(rec *engine.DataRecord) {
	table.Records = append(table.Records, rec)
}
