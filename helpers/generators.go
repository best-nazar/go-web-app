// helpers/generators.go

package helpers

import (
	"encoding/base64"
	"errors"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

// Generates session token with attached data
func GenerateSessionToken(data string, ip string) string {
	a := strconv.FormatInt(time.Now().UnixNano(), 10)
	b := []byte(a + "==" + data + "==" + ip)
	token := base64.StdEncoding.EncodeToString(b)

	return token
}

// Gets data from token
func RecoverSessionToken(token string, ip string) (int64, string, error) {
	decode := base64.NewDecoder(base64.StdEncoding, strings.NewReader(token))
	data, err := io.ReadAll(decode)

	if err != nil {
		log.Fatal(err, data)
		return 0, "", err
	} 

	contArr := strings.Split(string(data), "==")

	if len(contArr) != 3 || contArr[2] != ip {
		err := errors.New("token recovery error")

		return 0, "", err
	} else {
		datetime, err := strconv.ParseInt(contArr[0], 10, 64)
		userID := contArr[1]

		return datetime, userID, err
	}	
}