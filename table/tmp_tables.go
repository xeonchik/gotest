package table

import (
	"github.com/tidwall/btree"
	"godoc/engine"
)

type TemporaryDataSet struct {
	Tree    *btree.BTree
	Records map[uint16]*engine.PKItem
	Keys    []uint16
	Length  uint32
	i       int
}

func CreateTempTable(len int) *TemporaryDataSet {
	return &TemporaryDataSet{
		Records: map[uint16]*engine.PKItem{},
		Keys:    make([]uint16, len),
		i:       0,
	}
}

func (tmp *TemporaryDataSet) AddPK(pk *engine.PKItem) {
	tmp.Records[uint16(pk.PrimaryKey)] = pk
	tmp.Keys = append(tmp.Keys, uint16(pk.PrimaryKey))
}

func (tmp *TemporaryDataSet) Read() *engine.PKItem {
	length := len(tmp.Keys)
	pointer := tmp.i

	if pointer >= length {
		return nil
	}

	tmp.i++

	pk := tmp.Keys[pointer]
	return tmp.Records[pk]
}
