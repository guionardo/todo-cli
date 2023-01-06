package utils

import (
	"fmt"
	"time"
)

func singPlural(value int, singular, plural string) string {
	var template string
	if value == 1 {
		template = "%d " + singular
	} else {
		template = "%d " + plural
	}

	return fmt.Sprintf(template, value)
}

func DurationString(d time.Duration) string {
	if d.Hours() >= 24*365 {
		years := int(d.Hours() / 24 / 365)
		return singPlural(years, "year", "years")
	}
	if d.Hours() >= 24*30 {
		months := int(d.Hours() / 24 / 30)
		return singPlural(months, "month", "months")
	}
	if d.Hours() >= 24 {
		return singPlural(int(d.Hours()/24), "day", "days")
	}
	if d.Hours() >= 1 {
		return singPlural(int(d.Hours()), "hour", "hours")
	}
	d = d.Round(time.Minute * 15)
	return singPlural(int(d.Minutes()), "minute", "minutes")
}
