package main

import (
	"fmt"
	"github.com/Stroby241/UntisAPI"
	ics "github.com/arran4/golang-ical"
	"log"
	"math/rand"
	"os"
	"time"
)

var user *UntisAPI.User

const querryDays int = 90

type scedule struct {
	from time.Time
	till time.Time
}
type date struct {
	year  int
	month time.Month
	day   int
}

var timetable map[date]scedule

func main() {

	log.Printf("Loading Untis Data")

	user = UntisAPI.NewUser(
		"maarten8",
		"behn500",
		"TBZ Mitte Bremen",
		"https://tipo.webuntis.com")
	err := user.Login()
	if err != nil {
		log.Fatal(err)
	}

	personalId, err := user.GetPersonId("Maarten", "Behn", false)
	if err != nil {
		log.Fatal(err)
	}

	startTime := time.Now()
	endTime := startTime.AddDate(0, 0, querryDays)
	untisStartDate := UntisAPI.ToUntisDate(startTime)
	untisEndDate := UntisAPI.ToUntisDate(endTime)

	periods, err := user.GetTimeTable(personalId, 5, untisStartDate, untisEndDate)
	if err != nil {
		log.Fatal(err)
	}

	user.Logout()

	timetable = map[date]scedule{}

	for _, period := range periods {

		periodStartTime := UntisAPI.ToGoTime(period.StartTime)
		periodEndTime := UntisAPI.ToGoTime(period.EndTime)

		periodDate := date{}
		periodDate.year, periodDate.month, periodDate.day = UntisAPI.ToGoDate(period.Date).Date()

		log.Printf("\rPassing peroid on day %d.%d.%d", periodDate.day, periodDate.month, periodDate.year)

		day := timetable[periodDate]
		if day.from.IsZero() {
			day = scedule{periodStartTime, periodEndTime}
		} else {
			if day.from.Unix() > periodStartTime.Unix() {
				day.from = periodStartTime
			}
			if day.till.Unix() < periodEndTime.Unix() {
				day.till = periodEndTime
			}
		}
		timetable[periodDate] = day
	}

	cal := ics.NewCalendar()

	i := 0
	for date, day := range timetable {

		log.Printf("\rCreateing calendar event for %d.%d.%d", date.day, date.month, date.year)

		event := cal.AddEvent(fmt.Sprintf("%d", rand.Int()))
		i++

		event.SetCreatedTime(time.Now())
		event.SetModifiedAt(time.Now())
		event.SetDtStampTime(time.Now())
		event.SetSummary("Schule")

		event.SetStartAt(time.Date(date.year, date.month, date.day, day.from.Hour(), day.from.Minute(), 0, 0, time.Local))
		event.SetEndAt(time.Date(date.year, date.month, date.day, day.till.Hour(), day.till.Minute(), 0, 0, time.Local))
	}

	t := cal.Serialize()

	f, err := os.Create("untis.ics")
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.WriteString(t)
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Done")
}
