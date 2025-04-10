package timezone

import "time"

var ZONE_ID *time.Location

func init() {
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		println(err)
	}
	ZONE_ID = loc
}
