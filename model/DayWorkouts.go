package model

type DayWorkouts struct {
	Success bool          `json:"success"`
	Data    []WorkoutData `json:"data"`
}

type WorkoutData struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Workout Workout `json:"workout"`
}

type Workout struct {
	ID               int              `json:"id"`
	Date             string           `json:"date"`
	Created          string           `json:"created"`
	Tecniques        []Tecnique       `json:"tecniques"`
	Warmups          []Tecnique       `json:"warmups"`
	Wods             []Wod            `json:"wods"`
	Hours            []WorkoutHours   `json:"hours"`
	Results          []string         `json:"results"`
	ResultCategories ResultCategories `json:"result_categories"`
}

type Tecnique struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Exercises string `json:"exercises"`
}

type Wod struct {
	ID                 string        `json:"id"`
	Name               string        `json:"name"`
	Exercises          string        `json:"exercises"`
	Result_types       []ResultTypes `json:"result_types"`
	Result_sent        string        `json:"result_sent"`
	User_result        string        `json:"user_result"`
	Result_observation string        `json:"result_observation"`
}

type WorkoutHours struct {
	ID            int    `json:"id"`
	StartTime     string `json:"startTime"`
	EndTime       string `json:"endTime"`
	Coach         User   `json:"coach"`
	Limit         int    `json:"limit"`
	Users         int    `json:"users"`
	Users_list    []User `json:"users_list"`
	Did_checkin   bool   `json:"did_checkin"`
	Allow_checkin bool   `json:"allow_checkin"`
	Waiting_list  bool   `json:"waiting_list"`
	Allow_cancel  bool   `json:"allow_cancel"`
	Gympass       bool   `json:"gympass"`
}

type ResultTypes struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Result string `json:"result"`
}

type User struct {
	Id          int     `json:"id"`
	FacebookID  string  `json:"facebookID"`
	Email       string  `json:"email"`
	Name        string  `json:"name"`
	Last_name   string  `json:"last_name"`
	Privacy     int     `json:"privacy"`
	Birthday    string  `json:"birthday"`
	Gender      string  `json:"gender"`
	Weight      float64 `json:"weight"`
	City        string  `json:"city"`
	AccessToken string  `json:"accessToken"`
	Role        int     `json:"role"`
	Picture     string  `json:"picture"`
	Phone       string  `json:"phone"`
	User_status int     `json:"user_status"`
	Gympass     bool    `json:"gympass"`
}

type ResultCategories struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
