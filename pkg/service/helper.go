package service

import "time"

var (
	TimeLocation = time.UTC
)

// dateFormat 返回指定的时间格式
func dateFormat(t time.Time, layout string) string {
	return t.In(TimeLocation).Format(layout)
}
