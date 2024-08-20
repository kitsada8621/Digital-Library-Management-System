package utils

import (
	"time"
)

func StringToDate(dateString string) (*time.Time, error) {
	layout := "2006-01-02"
	date, err := time.Parse(layout, dateString)
	if err != nil {
		return nil, err
	}

	return &date, nil
}
