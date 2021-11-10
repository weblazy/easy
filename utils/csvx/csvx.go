package csv

import (
	"encoding/csv"
	"log"
	"os"
)

type CSV struct {
	file   *os.File
	wirter *csv.Writer
	reader *csv.Reader
}

func NewCSV(path string) *CSV {
	//OpenFile读取文件，不存在时则创建，使用追加模式
	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Println("文件打开失败！")
		return nil
	}
	wirter := csv.NewWriter(file)
	wirter.Comma = ';'
	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.LazyQuotes = true
	return &CSV{
		file:   file,
		wirter: wirter,
		reader: reader,
	}
}

func (this *CSV) Close() {
	this.file.Close()
}

//按行写入csv
func (this *CSV) WriterCSV(str []string) {
	//写入一条数据，传入数据为切片(追加模式)
	err := this.wirter.Write(str)
	if err != nil {
		log.Println("WriterCsv写入文件失败")
	}
	this.wirter.Flush() //刷新，不刷新是无法写入的
}

//按行读取csv
func (this *CSV) ReadRow() ([]string, error) {
	return this.reader.Read()
}

//按行读取csv
func (this *CSV) ReadAll() ([][]string, error) {
	return this.reader.ReadAll()
}

//按行读取csv
func (this *CSV) Reset() ([][]string, error) {
	return this.reader.ReadAll()
}
