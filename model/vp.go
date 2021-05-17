package model

import "time"

type Line struct {
	Unit          string
	Ticket        string
	Identity      string
	Matricula     string
	UseType       string
	CheckIn       time.Time
	CheckOut      time.Time
	Duration      int64
	PaidValue     float64
	PaymentMethod string
	Table         string
}
