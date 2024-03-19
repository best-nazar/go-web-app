package helpers

import (
	"log"
	"time"
	"unicode"

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
func TimestampToSting (timestamp time.Time) string {
	return carbon.CreateFromTimestamp(timestamp.Unix()).ToFormatString(dateFormat)
}

func Capitalize(str string) string {
	runes := []rune(str)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// Converts array of interface{} to slice
// Example of usege is reading arrayy from yaml file: appConfig.ImageConfig["extentions"]
func InterfaceArray(data interface{}) []string {
	var arrays []string
	exts := data.([]interface{})
	
	for _, val := range exts {
		arrays = append(arrays, val.(string))
	}

	return arrays
}