package engine

import (
	"fmt"
	"github.com/tidwall/btree"
)

type PKIndex struct {
	Tree *btree.BTree
}

type PKItem struct {
	Record     *DataRecord
	Locator    DataRowLocator
	PrimaryKey uint64
}

func byKey(a, b interface{}) bool {
	i1, i2 := a.(*PKItem), b.(*PKItem)
	return i1.PrimaryKey < i2.PrimaryKey
}

func (idx *PKIndex) Get(key uint64) *PKItem {
	StatsObj.Hits += 1
	item := idx.Tree.Get(&PKItem{
		PrimaryKey: key,
	}).(*PKItem)
	return item
}

func (idx *PKIndex) Load(record *DataRecord, locator DataRowLocator, key uint64) {
	item := PKItem{
		Record:     record,
		PrimaryKey: key,
		Locator:    locator,
	}
	idx.Tree.Load(&item)
}

func (idx *PKIndex) Add(record *DataRecord, locator DataRowLocator, key uint64) {
	item := PKItem{
		Record:     record,
		PrimaryKey: key,
		Locator:    locator,
	}

	if idx.Tree.Get(&item) != nil {
		panic("PK already exists")
	}

	idx.Tree.Set(&item)
}

func (idx *PKIndex) Print() {
	point := &PKItem{
		PrimaryKey: 0,
	}

	idx.Tree.Ascend(point, func(item interface{}) bool {
		it := item.(*PKItem)
		fmt.Println(it.Record.ID)
		return true
	})
}

func CreatePKIndex() *PKIndex {
	idx := &PKIndex{
		Tree: btree.New(byKey),
	}
	return idx
}
