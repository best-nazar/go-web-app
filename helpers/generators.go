// helpers/generators.go

package helpers

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

// Generates session token with attached data
func GenerateSessionToken(data string) string {
	a := strconv.FormatInt(time.Now().UnixNano(), 10)
	b := []byte(a + "==" + data)
	token := base64.StdEncoding.EncodeToString(b)

	return token
}

// Gets data from token
func RecoverSessionToken(token string) (int64, int64, error) {
	decode := base64.NewDecoder(base64.StdEncoding, strings.NewReader(token))
	data, err := ioutil.ReadAll(decode)

	if err != nil {
		log.Fatal(err, data)
		return 0, 0, err
	} 

	contArr := strings.Split(string(data), "==")

	if len(contArr) != 2 {
		err := errors.New("token recovery error")

		return 0, 0, err
	} else {
		datetime, err := strconv.ParseInt(contArr[0], 10, 64)
		userID, _ := strconv.ParseInt(contArr[1], 10, 64)

		return datetime, userID, err
	}	
}