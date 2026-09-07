package models

type HospitalResponse struct {
	Status     bool        `json:"status"`
	Message    string      `json:"message"`
	Pagination Pagination  `json:"pagination"`
	Data       interface{} `json:"data"`
}

type Pagination struct {
	TotalNumberOfPages int       `json:"totalNumberOfPages"`
	CurrentPage        int       `json:"currentPage"`
	Next               *NextPage `json:"next"`
}

type NextPage struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}
