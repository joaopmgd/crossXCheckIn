package app

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/crossXCheckIn/config"
	"github.com/crossXCheckIn/model"
	log "github.com/sirupsen/logrus"
)

func MakeCheckInForEverybody(config *config.Config) {
	for _, user := range config.Users {
		config.Logger.Info("Loggin with ", user.Login)
		token, err := userLogin(user, config.Endpoints["login"], config.Requests["login"])
		if err != nil {
			config.Logger.Error(err, user.Login)
		}
		config.Logger.Info("Searching workouts.")
		workoutId, timeId, err := userGetWorkouts(token, user.Workout, user.Time, config.Endpoints["workoutlist"], config.Requests["workoutlist"])
		if err != nil {
			config.Logger.Error(err)
		}
		config.Logger.Info("Making check-in.")
		status, err := userMakeCheckIn(token, workoutId, timeId, config.Endpoints["checkin"], config.Requests["checkin"])
		if err != nil {
			config.Logger.WithFields(log.Fields{
				"timeId":    timeId,
				"workoutId": workoutId,
				"login":     user.Login,
			}).Error(err)
		}
		config.Logger.Info(user.Login, status)
	}
}

func userLogin(user config.User, endpoint config.Endpoint, request config.Request) (string, error) {
	form := url.Values{
		"email":    {user.Login},
		"password": {user.Password},
		"type":     {request.Body["type"]},
	}
	var userLogin model.UserLogin
	status := makeRequest(form, "POST", endpoint.Url, nil, request, &userLogin)
	if status != "200 OK" {
		return "", errors.New("It was not possible to Login")
	}
	return userLogin.Data.AccessToken, nil
}

func userGetWorkouts(token, workoutType, workoutHour string, endpoint config.Endpoint, request config.Request) (string, string, error) {
	days, err := strconv.Atoi(os.Getenv("DAYS_FROM_NOW"))
	if err != nil {
		return "", "", errors.New("DAYS_FROM_NOW must be an integer.")
	}
	t := time.Now().Add(time.Hour * 24 * time.Duration(days))
	query := make(map[string]string)
	query["date"] = t.Format("2006-01-02")
	request.Header["X-ACCESSTOKEN"] = token
	var dayWorkouts model.DayWorkouts
	status := makeRequest(nil, "GET", endpoint.Url, query, request, &dayWorkouts)
	if status != "200 OK" {
		return "", "", errors.New("It was not possible to request the Workouts")
	}
	return findWorkoutData(workoutHour, workoutType, dayWorkouts.Data)
}

func userMakeCheckIn(token, workoutId, timeId string, endpoint config.Endpoint, request config.Request) (string, error) {
	form := url.Values{
		"time_id":    {timeId},
		"workout_id": {workoutId},
	}
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
