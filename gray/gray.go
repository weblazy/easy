package gray

import (
	"context"
	"math/rand"
	"strconv"
	"strings"
)

// GrayConfig 基于用户ID灰度开关配置
type GrayConfig struct {
	TailMemberIDs []string // 按照MemberID 尾号灰度
	MemberIDs     []string // 精准匹配 MemberID
	Percent       int64    // memberID 未命中灰度策略,根据比例灰度判断是否命中 范围0~100
}

type GraySwitch struct {
	config *GrayConfig
}

func NewGraySwitch(c *GrayConfig) *GraySwitch {
	return &GraySwitch{
		config: c,
	}
}

// IsGray 返回true表示命中灰度策略
func (g *GraySwitch) IsGray(ctx context.Context, memberID string) bool {
	// 100%全量
	if g.config.Percent == 100 {
		return true
	}

	// 判断灰度策略是否命中
	if memberID != "" {
		// 判断精准匹配
		for _, v := range g.config.MemberIDs {
			if v == memberID {
				return true
			}
		}
		// 判断尾号匹配
		for _, v := range g.config.TailMemberIDs {
			if strings.HasSuffix(memberID, v) {
				return true
			}
		}
	}

	// memberID 未命中灰度策略,根据比例灰度判断是否命中
	if g.config.Percent != 0 {
		if rand.Int63n(100) < g.config.Percent {
			return true
		}
	}

	return false
}

// IsGray 返回true表示命中灰度策略
func (g *GraySwitch) IsGrayByMemberIDInt(ctx context.Context, MemberIDInt int64) bool {
	// 100%全量
	if g.config.Percent == 100 {
		return true
	}
	memberID := strconv.FormatInt(MemberIDInt, 10)
	// 判断灰度策略是否命中
	if memberID != "" {
		// 判断精准匹配
		for _, v := range g.config.MemberIDs {
			if v == memberID {
				return true
			}
		}
		// 判断尾号匹配
		for _, v := range g.config.TailMemberIDs {
			if strings.HasSuffix(memberID, v) {
				return true
			}
		}
	}

	// memberID 未命中灰度策略,根据比例灰度判断是否命中
	if g.config.Percent != 0 {
		remainder := MemberIDInt % 100
		if remainder < g.config.Percent {
			return true
		}
	}

	return false
}
