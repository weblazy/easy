package bitmap

import (
	"encoding/base64"
	"errors"
	"strings"
	"sync"
)

const (
	bitSize = 8
)

var (
	bitmask     = []byte{1, 1 << 1, 1 << 2, 1 << 3, 1 << 4, 1 << 5, 1 << 6, 1 << 7}
	outRangeErr = errors.New("num out of range")
)

type bitmap struct {
	bits     []byte // byte字节内容
	count    uint64 // 已填入数字的数量
	byteSize uint64 // 容量byte字节数
	capacity uint64 // 容量bit数
	sync.RWMutex
}

// @desc NewBitmap
// @auth liuguoqiang 2020-12-27
// @param
// @return
func NewBitmap(maxNum uint64) *bitmap {
	byteSize := (maxNum + 7) >> 3
	return &bitmap{
		bits:     make([]byte, byteSize),
		count:    0,
		byteSize: byteSize,
		capacity: maxNum,
	}
}

// @desc 填入数字
// @auth liuguoqiang 2020-12-27
// @param
// @return
func (this *bitmap) Set(num uint64) error {
	if num > (this.capacity - 1) {
		return outRangeErr
	}
	byteIndex, bitPos := this.offset(num)
	// 1 左移 bitPos 位 进行 按位或 (置为 1)
	this.Lock()
	defer this.Unlock()
	this.bits[byteIndex] |= bitmask[bitPos]
	this.count++
	return nil
}

// @desc 清除填入的数字
// @auth liuguoqiang 2020-12-27
// @param
// @return
func (this *bitmap) Reset(num uint64) error {
	if num > (this.capacity - 1) {
		return outRangeErr
	}
	byteIndex, bitPos := this.offset(num)
	// 重置为空位 (重置为 0)
	this.Lock()
	defer this.Unlock()
	this.bits[byteIndex] &= ^bitmask[bitPos]
	this.count--
	return nil
}

// @desc 数字是否在位图中
// @auth liuguoqiang 2020-12-27
// @param
// @return
func (this *bitmap) Has(num uint64) bool {
	if num > (this.capacity - 1) {
		return false
	}
	byteIndex := num >> 3
	bitPos := num % bitSize
	// 右移 bitPos 位 和 1 进行 按位与
	this.RLock()
	defer this.RUnlock()
	return !(this.bits[byteIndex]&bitmask[bitPos] == 0)
}

// @desc  获取某个数字的的字节位置和bit位置
// @auth liuguoqiang 2020-12-27
// @param
// @return
func (this *bitmap) offset(num uint64) (uint64, byte) {
	byteIndex := num >> 3         // 字节索引
	bitPos := byte(num % bitSize) // bit位置
	return byteIndex, bitPos
}

// @desc 位图的容量
// @auth liuguoqiang 2020-12-27
// @param
// @return
func (this *bitmap) Capacity() uint64 {
	return this.capacity
}

// @desc 是否空位图
// @auth liuguoqiang 2020-12-27
// @param
// @return
func (this *bitmap) IsEmpty() bool {
	return this.count == 0
}

// @desc 是否已填满
// @auth liuguoqiang 2020-12-27
// @param
// @return
func (this *bitmap) IsFully() bool {
	return this.count == this.capacity
}

// @desc 已填入的数字个数
// @auth liuguoqiang 2020-12-27
// @param
// @return
func (this *bitmap) Count() uint64 {
	return this.count
}

// @desc 获取填入的数字切片
// @auth liuguoqiang 2020-12-27
// @param
// @return
func (this *bitmap) GetTotal() []uint64 {
	data := make([]uint64, 0)
	var lenth uint64 = 0
	if this.count == lenth {
		return data
	}
	this.RLock()
	defer this.RUnlock()
	for byteIndex := uint64(0); byteIndex < this.byteSize; byteIndex++ {
		if this.bits[byteIndex] == 0 {
			continue
		}
		for bitPos := range bitmask {
			if !(this.bits[byteIndex]&bitmask[bitPos] == 0) {
				data = append(data, byteIndex<<3+uint64(bitPos))
				lenth = uint64(len(data))
				//元素已经找满
				if this.count == lenth {
					return data
				}
			}
		}
	}
	return data
}

// @desc 获取字节数组格式内容
// @auth liuguoqiang 2020-12-27
// @param
// @return
func (this *bitmap) Bytes() []byte {
	return this.bits
}

// @desc 获取base64格式内容
// @auth liuguoqiang 2020-12-27
// @param
// @return
func (this *bitmap) Base64() string {
	return base64.StdEncoding.EncodeToString(this.bits)
}

// @desc 获取二进制格式内容
// @auth liuguoqiang 2020-12-27
// @param
// @return
func (this *bitmap) String() string {
	var sb strings.Builder
	byteIndex := this.byteSize
	this.RLock()
	defer this.RUnlock()
	for ; byteIndex >= 0; byteIndex-- {
		sb.WriteString(byteToBinaryString(this.bits[byteIndex]))
		sb.WriteString(" ")
	}
	return sb.String()
}

func byteToBinaryString(data byte) string {
	var sb strings.Builder
	for bitPos := 0; bitPos < bitSize; bitPos++ {
		if (bitmask[7-bitPos] & data) == 0 {
			sb.WriteString("0")
		} else {
			sb.WriteString("1")
		}
	}
	return sb.String()
}
