package engine

import (
	"fmt"
	"github.com/tidwall/btree"
)

type FlatIndex struct {
	Tree *btree.BTree
}

type FlatItem struct {
	Value uint64
}

func byFlatVal(a, b interface{}) bool {
	i1, i2 := a.(*FlatItem), b.(*FlatItem)
	return i1.Value < i2.Value
}

type MultiIndex struct {
	Tree *btree.BTree
}

type MultiItem struct {
	Keys     *FlatIndex
	IdxValue int
}

func byIdxVal(a, b interface{}) bool {
	i1, i2 := a.(*MultiItem), b.(*MultiItem)
	return i1.IdxValue < i2.IdxValue
}

// Add / indexValue - Record
/// PrimaryKey - PK Link
func (idx *MultiIndex) Add(indexValue int, key uint64) {
	// check that exists
	item := idx.Tree.Get(&MultiItem{
		IdxValue: indexValue,
	})

	if item == nil {
		item := &MultiItem{
			Keys:     &FlatIndex{Tree: btree.New(byFlatVal)},
			IdxValue: indexValue,
		}
		item.Keys.Tree.Set(&FlatItem{Value: key})
		idx.Tree.Set(item)
	} else {
		it := item.(*MultiItem)
		flatItem := &FlatItem{Value: key}
		it.Keys.Tree.Set(flatItem)
	}
}

func (idx *MultiIndex) GetTree() *btree.BTree {
	return idx.Tree
}

func (idx *MultiIndex) Print() {
	point := &MultiItem{
		IdxValue: 0,
	}

	idx.Tree.Ascend(point, func(item interface{}) bool {
		it := item.(*MultiItem)
		fmt.Println("Value: ", it.IdxValue)
		it.Keys.Tree.Ascend(nil, func(item2 interface{}) bool {
			fmt.Print(item2.(*FlatItem).Value, " ")
			return true
		})
		return true
	})
}

func (idx *MultiIndex) Get(key int) *MultiItem {
	StatsObj.Hits += 1
	return idx.Tree.Get(&MultiItem{
		IdxValue: key,
	}).(*MultiItem)
}

type IndexType interface {
	Add()
}

type BTIndex struct {
	Tree *btree.BTree
}

type FloatIndex struct {
	BTIndex
}

func (idx *FloatIndex) A() {
	panic("implement me")
}

type FloatItem struct {
	IdxValue float64
	Key      uint64
}

func byFloatVal(a, b interface{}) bool {
	i1, i2 := a.(*FloatItem), b.(*FloatItem)
	return i1.IdxValue < i2.IdxValue
}

func (idx *FloatIndex) Add(value float64, key uint64) {
	item := FloatItem{
		IdxValue: value,
		Key:      key,
	}

	idx.Tree.Set(&item)
}

func CreateFloatIndex() *FloatIndex {
	return &FloatIndex{
		BTIndex{Tree: btree.New(byFloatVal)},
	}
}

func CreateMulti() *MultiIndex {
	return &MultiIndex{
		Tree: btree.New(byIdxVal),
	}
}

func (idx *MultiIndex) AddArray(arr []int, key uint64) {
	for _, element := range arr {
		idx.Add(element, key)
	}
}
