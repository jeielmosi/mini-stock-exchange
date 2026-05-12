package dto_helper

import (
	"time"
	_ "time/tzdata"
)

func DateToEndOfDay(timestamp string) (time.Time, error) {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		return time.Time{}, err
	}

	bg, err := time.ParseInLocation(time.DateOnly, timestamp, loc)
	if err != nil {
		return bg, err
	}

	ed := time.Date(bg.Year(), bg.Month(), bg.Day()+1, 0, 0, 0, 0, loc).Add(-time.Nanosecond)
	return ed, nil
}
