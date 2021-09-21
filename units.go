package main

import (
	"UnitsAPI"
	"log"
	"time"
)

var user *UntisAPI.User

const querryDays int = 20

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
	user = UntisAPI.NewUser(
		"maarten8",
		"behn500",
		"TBZ Mitte Bremen",
		"https://tipo.webuntis.com")
	err := user.Login()
	if err != nil {
		log.Fatal(err)
		return
	}

	personalId, err := user.GetPersonId("Maarten", "Behn", false)
	if err != nil {
		log.Fatal(err)
		return
	}

	startTime := time.Now()
	endTime := startTime.AddDate(0, 0, querryDays)
	untisStartDate := UntisAPI.ToUntisDate(startTime)
	untisEndDate := UntisAPI.ToUntisDate(endTime)

	//teacherList, err := user.GetTeachers()
	//if err != nil{ log.Fatal(err); return}

	roomList, err := user.GetRooms()
	if err != nil {
		log.Fatal(err)
		return
	}

	periods, err := user.GetTimeTable(personalId, 5, untisStartDate, untisEndDate)
	if err != nil {
		log.Fatal(err)
		return
	}

	timetable = map[date]scedule{}

	for _, period := range periods {

		/*
			var teachers []*UntisAPI.Teacher
			for _, t := range teacherList {
				for _, t2 := range period.Teacher {
					if t.Id == t2 {
						teachers = append(teachers, &t)
					}
				}
			}
		*/

		var rooms []*UntisAPI.Room
		for _, r := range roomList {
			for _, r2 := range period.Rooms {
				if r.Id == r2 {
					rooms = append(rooms, &r)
				}
			}
		}

		periodStartTime := UntisAPI.ToGoTime(period.StartTime)
		periodEndTime := UntisAPI.ToGoTime(period.EndTime)

		periodDate := date{}
		periodDate.year, periodDate.month, periodDate.day = UntisAPI.ToGoDate(period.Date).Date()

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

	/*
		log.Printf("Date: %02d.%02d.%04d From: %d:%02d Till: %d:%02d\n",
			periodDate.day, periodDate.month, periodDate.year,
			periodStartTime.Hour(), periodStartTime.Minute(),
			periodEndTime.Hour(), periodEndTime.Minute())
	*/

	user.Logout()
}
