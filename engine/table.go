package engine

import "errors"

type Table struct {
	PrimaryIndex *PKIndex
	Indexes      map[string]IndexDefinition
}

type IndexDefinition struct {
	index   interface{}
	indexer indexer
}

type indexer func(index interface{}, rec *DataRecord)

var tableSpace = make(map[string]*Table)

func InitTableSpace(name string) *Table {
	primaryIdx := CreatePKIndex()
	table := &Table{
		PrimaryIndex: primaryIdx,
		Indexes:      make(map[string]IndexDefinition),
	}
	tableSpace[name] = table
	return table
}

func (tbl *Table) AddIndexer(name string, index interface{}, indexFn indexer) error {
	_, exist := tbl.Indexes[name]
	if exist {
		return errors.New("idx with name " + name + " already exist")
	}

	indexDef := IndexDefinition{
		index:   index,
		indexer: indexFn,
	}
	tbl.Indexes[name] = indexDef

	return nil
}

func (tbl *Table) GetByPK(pk uint64) *DataRecord {
	pkItem := tbl.PrimaryIndex.Get(pk)
	return pkItem.Record
}

func (tbl *Table) ReadPageRecords(page *Page) {
	var offset int64 = 0
	for {
		entity := page.ReadDataRecord(offset)
		if entity == nil {
			break
		}
		offset += int64(entity.Locator.Size)

		// add to PK index
		tbl.PrimaryIndex.Add(entity.Record, entity.Locator, entity.Record.ID)
		tbl.RecordAddToIndexes(entity.Record)
	}
}

func (tbl *Table) RecordAddToIndexes(record *DataRecord) {
	for _, indexDefinition := range tbl.Indexes {
		indexDefinition.indexer(indexDefinition.index, record)
	}
}
