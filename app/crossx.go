package app

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/crossXCheckIn/config"
	"github.com/crossXCheckIn/model"
	log "github.com/sirupsen/logrus"
)

func MakeCheckInForEverybody(config *config.Config) {
	config.Logger.Info("Searching workouts.")
	workoutId, timeId, err := userGetWorkouts(config.Users[0], config.Endpoints["workoutlist"], config.Requests["workoutlist"])
	if err != nil {
		config.Logger.WithFields(log.Fields{
			"user.Workout":     config.Users[0].Workout,
			"user.WeekdayHour": config.Users[0].WeekdayHour,
			"user.WeekendHour": config.Users[0].WeekendHour,
			"login":            config.Users[0].Login,
		}).Error(err)
		return
	}
	config.Logger.Info("Starting check-in.")
	for _, user := range config.Users {
		status, err := userMakeCheckIn(user.Token, workoutId, timeId, config.Endpoints["checkin"], config.Requests["checkin"])
		if err != nil {
			config.Logger.WithFields(log.Fields{
				"login":     user.Login,
				"timeId":    timeId,
				"workoutId": workoutId,
				"status":    status,
			}).Error(err)
			continue
		}
		config.Logger.WithFields(log.Fields{
			"login":       user.Login,
			"status":      status,
			"WeekdayHour": user.WeekdayHour,
			"WeekendHour": user.WeekendHour,
		}).Info("Check-in done in 3 days.")
	}
}

func userLogin(user config.User, endpoint config.Endpoint, request config.Request) (string, error) {
	form := url.Values{
		"email":    {user.Login},
		"password": {user.Password},
		"type":     {request.Body["type"]},
	}
	parameters := url.Values{}
	parameters.Add("email", user.Login)
	parameters.Add("password", user.Password)
	parameters.Add("type", request.Body["type"])
	request.Header["Content-Length"] = strconv.Itoa(len(parameters.Encode()))
	var userLogin model.UserLogin
	status := makeRequest(form, "POST", endpoint.Url, nil, request, &userLogin)
	if status != "200 OK" {
		return "", errors.New("It was not possible to Login")
	}
	return userLogin.Data.AccessToken, nil
}

func userGetWorkouts(user config.User, endpoint config.Endpoint, request config.Request) (string, string, error) {
	t := time.Now().Add(time.Hour * 24 * time.Duration(user.DaysAhead))
	query := make(map[string]string)
	query["date"] = t.Format("2006-01-02")
	request.Header["X-ACCESSTOKEN"] = user.Token
	var dayWorkouts model.DayWorkouts
	status := makeRequest(nil, "GET", endpoint.Url, query, request, &dayWorkouts)
	if status != "200 OK" {
		return "", "", errors.New("It was not possible to request the Workouts")
	}
	var checkinHour string
	if t.Weekday() == 6 {
		checkinHour = user.WeekendHour
	} else if t.Weekday() == 0 {
		return "", "", errors.New("Sunday does not have Crossfit")
	} else {
		checkinHour = user.WeekdayHour
	}
	return findWorkoutData(checkinHour, user.Workout, dayWorkouts.Data)
}

func userMakeCheckIn(token, workoutId, timeId string, endpoint config.Endpoint, request config.Request) (string, error) {
	form := url.Values{
		"time_id":    {timeId},
		"workout_id": {workoutId},
	}
	parameters := url.Values{}
	parameters.Add("time_id", timeId)
	parameters.Add("workout_id", workoutId)
	request.Header["Content-Length"] = strconv.Itoa(len(parameters.Encode()))
	request.Header["X-ACCESSTOKEN"] = token
	status := makeRequest(form, "POST", endpoint.Url, nil, request, nil)
	if status != "204 No Content" {
		return "", errors.New("It was not possible to Check-in")
	}
	return status, nil
}

func makeRequest(form url.Values, method, endpoint string, query map[string]string, request config.Request, target interface{}) string {
	client := &http.Client{}
	req, err := http.NewRequest(method, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		log.Fatal(err)
	}

	if query != nil {
		q := req.URL.Query()
		for key, value := range query {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	req.PostForm = form
	for key, value := range request.Header {
		req.Header.Add(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if target != nil {
		json.NewDecoder(resp.Body).Decode(target)
	}

	return resp.Status
}

func findWorkoutData(workoutStartHour, workoutType string, workoutsData []model.WorkoutData) (string, string, error) {
	var workoutId, timeId string
	for _, workoutData := range workoutsData {
		if workoutData.Name == workoutType {
			workoutId = strconv.Itoa(workoutData.Workout.Id)
			for _, workoutHour := range workoutData.Workout.Hours {
				if workoutHour.StartTime == workoutStartHour {
					timeId = strconv.Itoa(workoutHour.Id)
					return workoutId, timeId, nil
				}
			}
		}
	}
	return workoutId, timeId, errors.New("No workout found with the specified config")
}

// TODO: add gympass call
