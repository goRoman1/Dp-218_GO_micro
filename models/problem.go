package models

import "time"

// ProblemType - entity for types of problems user reported
type ProblemType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Problem - entity for problem representation in the system
type Problem struct {
	ID           int         `json:"id"`
	User         User        `json:"user"`
	Type         ProblemType `json:"type"`
	DateReported time.Time   `json:"date_reported"`
	Description  string      `json:"description"`
	IsSolved     bool        `json:"is_solved"`
}

// ProblemList - struct for list of problems
type ProblemList struct {
	Problems []Problem `json:"accounts"`
}

// Solution - entity for solution representation in the system
type Solution struct {
	Problem     Problem   `json:"problem"`
	DateSolved  time.Time `json:"date_solved"`
	Description string    `json:"description"`
}
