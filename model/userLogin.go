package model

type UserLogin struct {
	Success bool          `json:"success"`
	Data    DataUserLogin `json:"data"`
}

type DataUserLogin struct {
	Id          string       `json:"id"`
	FacebookID  string       `json:"facebookID"`
	Email       string       `json:"email"`
	Name        string       `json:"name"`
	LastName    string       `json:"last_name"`
	Privacy     int          `json:"privacy"`
	Birthday    string       `json:"birthday"`
	Gender      string       `json:"gender"`
	Weight      float64      `json:"weight"`
	City        string       `json:"city"`
	AccessToken string       `json:"accessToken"`
	Role        int          `json:"role"`
	Picture     string       `json:"picture"`
	Phone       string       `json:"phone"`
	Box         BoxUserLogin `json:"box"`
	Rg          string       `json:"rg"`
	Address     string       `json:"address"`
	Cpf         string       `json:"cpf"`
}

type BoxUserLogin struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Address1  string `json:"address1"`
	Address2  string `json:"address2"`
	Telephone string `json:"telephone"`
	Website   string `json:"website"`
	LogoUrl   string `json:"logoUrl"`
	Affiliate bool   `json:"affiliate"`
}
