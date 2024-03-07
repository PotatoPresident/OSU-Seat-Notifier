package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type CourseSearchResponse struct {
	Success    bool `json:"success"`
	TotalCount int  `json:"totalCount"`
	Sections   []struct {
		ID                      int    `json:"id"`
		Term                    string `json:"term"`
		TermDesc                string `json:"termDesc"`
		CourseReferenceNumber   string `json:"courseReferenceNumber"`
		PartOfTerm              string `json:"partOfTerm"`
		CourseNumber            string `json:"courseNumber"`
		Subject                 string `json:"subject"`
		SubjectDescription      string `json:"subjectDescription"`
		SequenceNumber          string `json:"sequenceNumber"`
		CampusDescription       string `json:"campusDescription"`
		ScheduleTypeDescription string `json:"scheduleTypeDescription"`
		CourseTitle             string `json:"courseTitle"`
		CreditHours             any    `json:"creditHours"`
		MaximumEnrollment       int    `json:"maximumEnrollment"`
		Enrollment              int    `json:"enrollment"`
		SeatsAvailable          int    `json:"seatsAvailable"`
		WaitCapacity            int    `json:"waitCapacity"`
		WaitCount               int    `json:"waitCount"`
		WaitAvailable           int    `json:"waitAvailable"`
		CrossList               any    `json:"crossList"`
		CrossListCapacity       any    `json:"crossListCapacity"`
		CrossListCount          any    `json:"crossListCount"`
		CrossListAvailable      any    `json:"crossListAvailable"`
		CreditHourHigh          any    `json:"creditHourHigh"`
		CreditHourLow           int    `json:"creditHourLow"`
		CreditHourIndicator     any    `json:"creditHourIndicator"`
		OpenSection             bool   `json:"openSection"`
		LinkIdentifier          any    `json:"linkIdentifier"`
		IsSectionLinked         bool   `json:"isSectionLinked"`
		SubjectCourse           string `json:"subjectCourse"`
		Faculty                 []struct {
			BannerID              string `json:"bannerId"`
			Category              any    `json:"category"`
			Class                 string `json:"class"`
			CourseReferenceNumber string `json:"courseReferenceNumber"`
			DisplayName           string `json:"displayName"`
			EmailAddress          string `json:"emailAddress"`
			PrimaryIndicator      bool   `json:"primaryIndicator"`
			Term                  string `json:"term"`
		} `json:"faculty"`
		MeetingsFaculty []struct {
			Category              string `json:"category"`
			Class                 string `json:"class"`
			CourseReferenceNumber string `json:"courseReferenceNumber"`
			Faculty               []any  `json:"faculty"`
			MeetingTime           struct {
				BeginTime              string  `json:"beginTime"`
				Building               string  `json:"building"`
				BuildingDescription    string  `json:"buildingDescription"`
				Campus                 string  `json:"campus"`
				CampusDescription      string  `json:"campusDescription"`
				Category               string  `json:"category"`
				Class                  string  `json:"class"`
				CourseReferenceNumber  string  `json:"courseReferenceNumber"`
				CreditHourSession      float64 `json:"creditHourSession"`
				EndDate                string  `json:"endDate"`
				EndTime                string  `json:"endTime"`
				Friday                 bool    `json:"friday"`
				HoursWeek              float64 `json:"hoursWeek"`
				MeetingScheduleType    string  `json:"meetingScheduleType"`
				MeetingType            string  `json:"meetingType"`
				MeetingTypeDescription string  `json:"meetingTypeDescription"`
				Monday                 bool    `json:"monday"`
				Room                   string  `json:"room"`
				Saturday               bool    `json:"saturday"`
				StartDate              string  `json:"startDate"`
				Sunday                 bool    `json:"sunday"`
				Term                   string  `json:"term"`
				Thursday               bool    `json:"thursday"`
				Tuesday                bool    `json:"tuesday"`
				Wednesday              bool    `json:"wednesday"`
			} `json:"meetingTime"`
			Term string `json:"term"`
		} `json:"meetingsFaculty"`
		ReservedSeatSummary            any    `json:"reservedSeatSummary"`
		SectionAttributes              []any  `json:"sectionAttributes"`
		InstructionalMethod            string `json:"instructionalMethod"`
		InstructionalMethodDescription string `json:"instructionalMethodDescription"`
	} `json:"data"`
	PageOffset           int    `json:"pageOffset"`
	PageMaxSize          int    `json:"pageMaxSize"`
	SectionsFetchedCount int    `json:"sectionsFetchedCount"`
	PathMode             string `json:"pathMode"`
	SearchResultsConfigs []struct {
		Config   string `json:"config"`
		Display  string `json:"display"`
		Title    string `json:"title"`
		Required bool   `json:"required"`
		Width    string `json:"width"`
	} `json:"searchResultsConfigs"`
	ZtcEncodedImage string `json:"ztcEncodedImage"`
}

func GetCourse(subject string, courseNumber string) (CourseSearchResponse, error) {
	client := &http.Client{}

	data := url.Values{}
	data.Set("term", "202403")

	postReq, err := http.NewRequest("POST", "https://prodapps.isadm.oregonstate.edu/StudentRegistrationSsb/ssb/term/search?mode=search", strings.NewReader(data.Encode()))
	postReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(postReq)
	if err != nil {
		return CourseSearchResponse{}, err
	}

	cookies := strings.Join(res.Header.Values("Set-Cookie"), "; ")

	getUrl := "https://prodapps.isadm.oregonstate.edu/StudentRegistrationSsb/ssb/searchResults/searchResults?" +
		"&chk_open_only=true" +
		"&txt_term=202403" +
		"&pageMaxSize=500" +
		"&txt_campus=C" +
		"&txt_subject=" + subject +
		"&txt_courseNumber=" + courseNumber

	req, err := http.NewRequest("GET", getUrl, nil)
	if err != nil {
		return CourseSearchResponse{}, err
	}

	req.Header.Add("Cookie", cookies)

	resp, err := client.Do(req)
	if err != nil {
		return CourseSearchResponse{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return CourseSearchResponse{}, err
	}

	var result CourseSearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return CourseSearchResponse{}, err
	}

	return result, nil
}

func trackCourse(subject string, number string, duration time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()

	for range time.Tick(duration) {
		course, err := GetCourse(subject, number)
		if err != nil {
			log.Println("Failed to fetch course", subject, number, "from api")
			log.Println(err)
			continue
		}
		log.Println(subject + number)
		hasOpenSeats := false
		for _, section := range course.Sections {
			if section.SeatsAvailable > 0 {
				log.Printf("CRN %v (%v) has %v seats avalible", section.CourseReferenceNumber, section.ScheduleTypeDescription, section.SeatsAvailable)
				hasOpenSeats = true
			}
		}
		if !hasOpenSeats {
			log.Println("No open seats")
		}
	}
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(3)

	duration := 10 * time.Second

	go trackCourse("CS", "162", duration, &wg)
	go trackCourse("CS", "261", duration, &wg)
	go trackCourse("CS", "271", duration, &wg)

	wg.Wait()
}
