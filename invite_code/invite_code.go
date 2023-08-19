package invite_code

import (
	"math/rand"
	"time"
)

// 自定义进制(0,1没有加入,容易与o,l混淆)，数组顺序可进行调整增加反推难度，A用来补位因此此数组不包含A，共31个字符。
var BaseByte = []byte{'H', 'V', 'E', '8', 'S', '2', 'D', 'Z', 'X', '9', 'C', '7', 'P',
	'5', 'I', 'K', '3', 'M', 'J', 'U', 'F', 'R', '4', 'W', 'Y', 'L', 'T', 'N', '6', 'B', 'G', 'Q'}

// 邀请码长度
var CodeLenth = 8

// A补位字符，不能与自定义重复
var SuffixByte byte = 'A'

// 默认邀请码生成器
var DefaultInviteCodeHandler = NewInviteCodeHandler(BaseByte, CodeLenth, SuffixByte)

type InviteCodeHandler struct {
	BaseByte   []byte
	baseLength int
	CodeLength int
	SuffixByte byte
}

func NewInviteCodeHandler(baseByte []byte, codeLength int, suffixByte byte) *InviteCodeHandler {
	return &InviteCodeHandler{
		baseLength: len(baseByte),
		BaseByte:   baseByte,
		CodeLength: codeLength,
		SuffixByte: suffixByte,
	}
}

func (c *InviteCodeHandler) IdToCode(id int) string {
	buf := make([]byte, c.baseLength)
	charPos := c.baseLength
	for id/c.baseLength > 0 {
		index := id % c.baseLength
		charPos--
		buf[charPos] = c.BaseByte[index]
		id /= c.baseLength
	}
	charPos--
	buf[charPos] = c.BaseByte[id%c.baseLength]
	// 将字符数组转化为字符串
	result := buf[charPos:]
	// 长度不足指定长度则随机补全
	length := len(result)
	if length < int(c.CodeLength) {
		result = append(result, c.SuffixByte)
		now := time.Now().UnixNano()
		rand.Seed(now)
		// 去除SuffixByte本身占位之后需要补齐的位数
		for i := 0; i < c.CodeLength-length-1; i++ {
			randomNum := rand.Intn(c.baseLength)
			result = append(result, c.BaseByte[randomNum])
		}
	}
	return string(result)
}

func (c *InviteCodeHandler) CodeToId(code string) int {
	var result int
	for i := 0; i < len(code); i++ {
		var index int
		for j := 0; j < c.baseLength; j++ {
			if code[i] == c.BaseByte[j] {
				index = j
				break
			}
		}
		if code[i] == c.SuffixByte {
			break
		}

		if i > 0 {
			result = result*c.baseLength + index
		} else {
			result = index
		}
	}
	return result
}
