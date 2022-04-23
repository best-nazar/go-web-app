package helpers

import (
	"log"
	"time"
	"github.com/golang-module/carbon/v2"
)

const dateFormat = "d-m-Y"

// Converting sting date-time with Time-zone designator to Time struct
// "2014-11-12T11:45:26.371Z"
func StringToTimeFull (datetime string) time.Time {
	t, err := time.Parse(time.RFC3339, datetime)

	if err != nil {
		log.Fatal(err)
	}
	
	return t
}

// Converting sting date-time to Time struct
// (see const dateFormat)
// https://github.com/golang-module/carbon
func StringToTimestamp (datetime string) int64 {
	return carbon.ParseByFormat(datetime, dateFormat).Timestamp()
}

//Converting timestamp to string date (see const dateFormat)
// https://github.com/golang-module/carbon
func TimestampToSting (timestamp int64) string {
	return carbon.CreateFromTimestamp(timestamp).ToFormatString(dateFormat)
}