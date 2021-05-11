package domain

import (
	"fmt"
	"strconv"
)

type Summary struct {
	Files         int    `json:"files"`
	Records       int    `json:"records"`
	Valid         int    `json:"valid"`
	Invalid       int    `json:"invalid"`
	Duplicated    int    `json:"duplicated"`
	ExecutionTime string `json:"execution_time"`
}

func (ref *Summary) ToString() []string {
	return []string{strconv.Itoa(ref.Files), strconv.Itoa(ref.Records), strconv.Itoa(ref.Duplicated), strconv.Itoa(ref.Valid), strconv.Itoa(ref.Invalid), ref.ExecutionTime}
}

type Employee struct {
	Id              string            `csv:"Number"`
	Name            string            `csv:"Name"`
	Email           string            `csv:"Email"`
	Salary          float64           `csv:"Wage"`
	AditionalFields map[string]string `csv:"-"`
}

func (ref *Employee) String() string {
	return fmt.Sprintf("%s,%s,%s,%f", ref.Id, ref.Name, ref.Email, ref.Salary)
}

func (ref *Employee) ToString() []string {
	return []string{ref.Id, ref.Name, ref.Email, fmt.Sprintf("%.2f", ref.Salary)}
}

func (ref *Employee) IsValid() bool {
	return ref.Email != "" && ref.Id != "" && ref.Name != "" && ref.Salary > 0
}
