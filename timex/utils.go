package timex

import "time"

const (
	TimeLayout = "2006-01-02 15:04:05"
	DateLayout = "2006-01-02"
)

// @desc 获取某一天的0点时间
// @auth liuguoqiang 2020-04-27
// @param
// @return
func ZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.Local)
}

func ZeroTimeWithLocation(d time.Time, loc *time.Location) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, loc)
}

// @desc 返回一个月的开始时间和结束时间
// @auth liuguoqiang 2020-04-27
// @param
// @return
func MonthRange(timeStamp int64) (int64, int64) {
	d := time.Unix(timeStamp, 0)
	d = d.AddDate(0, 0, -d.Day()+1)
	start := ZeroTime(d)
	end := start.AddDate(0, 1, 0)
	return start.Unix(), end.Unix()
}

func MonthRangeWthiLocation(timeStamp int64, loc *time.Location) (int64, int64) {
	d := time.Unix(timeStamp, 0)
	d = d.AddDate(0, 0, -d.Day()+1)
	start := ZeroTimeWithLocation(d, loc)
	end := start.AddDate(0, 1, 0)
	return start.Unix(), end.Unix()
}
