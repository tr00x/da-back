package utils

import (
	"math/rand"
	"mime/multipart"
	"strconv"
	"time"
)

func IsPDF(file *multipart.FileHeader) bool {
	f, err := file.Open()

	if err != nil {
		return false
	}

	defer f.Close()
	// PDF files start with "%PDF"
	buf := make([]byte, 4)
	_, err = f.Read(buf)

	if err != nil {
		return false
	}

	return string(buf) == "%PDF"
}

// return the current time in GMT + 0 timezone
func GMTTime() time.Time {
	return time.Now().In(time.FixedZone("GMT", 0))
}

func CheckLastIDLimit(lastID, limit, typeOfQuery string) (int, int) {
	lastIDInt, limitInt := 0, 0
	if typeOfQuery == "chat" {
		lastIDInt, limitInt = 9999, 50
	} else {
		lastIDInt, limitInt = 0, 50
	}

	if lastID == "" || limit == "" || len(lastID) > 9 {
		return lastIDInt, limitInt
	}

	limitInt, _ = strconv.Atoi(limit)
	lastIDInt, _ = strconv.Atoi(lastID)

	return lastIDInt, limitInt
}

// check gmt + 0
func CheckGMTTime(t time.Time) bool {
	_, offset := t.Zone()
	return offset == 0
}

func RandomOTP() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(900000) + 100000
}

func RandomUsername() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)

	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}
