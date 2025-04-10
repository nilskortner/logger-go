package timezone

import (
	"fmt"
	"time"
)

func ToBytes(t time.Time) []byte {
	year := fmt.Sprintf("%04d", t.Year())
	month := fmt.Sprintf("%02d", t.Month())
	day := fmt.Sprintf("%02d", t.Day())
	hour := fmt.Sprintf("%02d", t.Hour())
	minute := fmt.Sprintf("%02d", t.Minute())
	second := fmt.Sprintf("%02d", t.Second())
	millis := fmt.Sprintf("%03d", t.Nanosecond()/1e6)

	dateString := fmt.Sprintf("%s-%s-%s %s:%s:%s.%s", year, month, day, hour, minute, second, millis)

	return []byte(dateString)
}
