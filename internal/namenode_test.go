package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNameNodeDataAddDirData(t *testing.T) {
	assert := assert.New(t)

	dataset := "mnist"
	dirnode := NewNode("train/ddd/", 999, true)
	dirnode1 := NewNode("test/sd1/", 999, true)
	dirnode2 := NewNode("train/ddd/mmd3/", 999, true)

	var mydata *NameNodeData = new(NameNodeData)
	mydata.AddDataSetData(dataset, dirnode)
	mydata.AddDataSetData(dataset, dirnode1)
	mydata.AddDataSetData(dataset, dirnode2)

	num := mydata.DataSetDirNum(dataset)
	assert.Equal(num, 3, "they should be equal")
}

func TestNameNodeDataAddFileData(t *testing.T) {
	assert := assert.New(t)

	dataset := "mnist"
	filenode := NewNode("train/ddd/a.jpg", 999, false)
	filenode1 := NewNode("test/sd1/b.jpg", 999, false)
	filenode2 := NewNode("train/ddd/mmd3/c.jpg", 999, false)

	var mydata *NameNodeData = new(NameNodeData)
	mydata.AddDataSetData(dataset, filenode)
	mydata.AddDataSetData(dataset, filenode1)
	mydata.AddDataSetData(dataset, filenode2)

	num := mydata.DataSetFileNum(dataset)
	assert.Equal(num, 3, "they should be equal")
}

func TestNameNodeDataAddData(t *testing.T) {
	assert := assert.New(t)

	dataset := "mnist"
	dirnode := NewNode("train/ddd/", 999, true)
	dirnode1 := NewNode("test/sd1/", 999, true)
	dirnode2 := NewNode("train/ddd/mmd3/", 999, true)
	filenode := NewNode("train/ddd/a.jpg", 999, false)
	filenode1 := NewNode("test/sd1/b.jpg", 999, false)
	filenode2 := NewNode("train/ddd/mmd3/c.jpg", 999, false)
	filenode3 := NewNode("train/ddd/mmd3/d.jpg", 999, false)

	var mydata *NameNodeData = new(NameNodeData)
	mydata.AddDataSetData(dataset, dirnode)
	mydata.AddDataSetData(dataset, dirnode1)
	mydata.AddDataSetData(dataset, dirnode2)
	mydata.AddDataSetData(dataset, filenode)
	mydata.AddDataSetData(dataset, filenode1)
	mydata.AddDataSetData(dataset, filenode2)
	mydata.AddDataSetData(dataset, filenode3)

	num := mydata.DataSetFileNum(dataset)
	assert.Equal(num, 4, "they should be equal")
	num = mydata.DataSetDirNum(dataset)
	assert.Equal(num, 3, "they should be equal")
	num = mydata.DataSetNum(dataset)
	assert.Equal(num, 7, "they should be equal")
}
