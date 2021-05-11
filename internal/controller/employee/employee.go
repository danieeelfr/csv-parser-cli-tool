package employee

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/danieeelfr/csv_extractor_cli/internal/domain"
)

type Employee interface {
	FetchFromStringCSV(data [][]string) []*domain.Employee
	Dedupe(data []*domain.Employee) []*domain.Employee
	SplitValidAndInvalidRecords(data []*domain.Employee) (valid []*domain.Employee, invalid []*domain.Employee)
	ToCSVString(data []*domain.Employee) [][]string
}

type employee struct {
}

func NewEmployee() Employee {
	return new(employee)
}

func (ref *employee) FetchFromStringCSV(data [][]string) []*domain.Employee {
	var employee *domain.Employee
	var employees []*domain.Employee

	idx_id, idx_email, idx_salary, idx_name := getColumnsIndexes(data[0])

	for i, line := range data {
		if i == 0 {
			continue
		}

		employee = new(domain.Employee)
		if line[idx_name[0]] == line[idx_name[1]] {
			employee.Name = strings.TrimSpace(line[idx_name[0]])
		} else {
			employee.Name = strings.TrimSpace(fmt.Sprintf("%s %s", line[idx_name[0]], line[idx_name[1]]))
		}

		employee.Id = strings.TrimSpace(line[idx_id])
		employee.Email = strings.TrimSpace(line[idx_email])

		s := strings.ReplaceAll(strings.ReplaceAll(line[idx_salary], "$", ""), ",", "")
		sal, err := strconv.ParseFloat(s, 64)
		if err != nil {
			sal = 0
		}
		employee.Salary = sal

		// TODO: bind the aditional fields
		employees = append(employees, employee)
	}

	return employees
}

func (ref *employee) Dedupe(data []*domain.Employee) []*domain.Employee {
	out := make([]*domain.Employee, 0)
	unique := make(map[string]int, 0)
	for i, e := range data {
		if isDuplicated(unique, e.Email) {
			continue
		}
		out = append(out, e)
		unique[e.Email] = i
	}

	return out
}

func (ref *employee) SplitValidAndInvalidRecords(data []*domain.Employee) ([]*domain.Employee, []*domain.Employee) {

	validRecords, invalidRecords := make([]*domain.Employee, 0), make([]*domain.Employee, 0)

	for _, e := range data {
		if e.IsValid() {
			validRecords = append(validRecords, e)
		} else {
			invalidRecords = append(invalidRecords, e)
		}
	}

	return validRecords, invalidRecords
}

func (ref *employee) ToCSVString(data []*domain.Employee) [][]string {

	out := make([][]string, 0)
	if data == nil || len(data) == 0 {
		return out
	}

	out = append(out, []string{"Id", "Name", "Email", "Salary"})

	for _, l := range data {
		if l.Email == "" && l.Id == "" && l.Name == "" {
			continue
		}
		out = append(out, l.ToString())
	}

	return out
}

func getColumnsIndexes(line []string) (id, email, salary int, name []int) {
	var idx_id, idx_email, idx_salary int
	var idx_name = make([]int, 2)

	for i := 0; i < len(line); i++ {
		idx_name = idx_name[:2]

		col_name := strings.ToLower(strings.TrimSpace(line[i]))

		switch col_name {

		case "email", "e-mail", "e mail":
			idx_email = i
		case "wage", "salary", "rate":
			idx_salary = i
		case "number", "employeenumber", "emp id", "id":
			idx_id = i
		default:
			firstNameMatched, err := regexp.MatchString(`(f.\s(name)|first)`, col_name)
			checkError(err)
			if firstNameMatched {
				idx_name[0] = i
				continue
			}

			lastNameMatched, err := regexp.MatchString(`(l.\s(name)|last)`, col_name)
			checkError(err)
			if lastNameMatched {
				idx_name[1] = i
				continue
			}

			fullNameMatched, err := regexp.MatchString(`(name)`, col_name)
			checkError(err)
			if fullNameMatched {
				idx_name[0] = i
				idx_name[1] = i
				continue
			}

			// fmt.Println(fmt.Sprintf("column [%s] not binded!", col_name))
		}
	}

	return idx_id, idx_email, idx_salary, idx_name
}

func isDuplicated(data map[string]int, email string) bool {
	// for _, v := range data {
	// 	if v == email {
	// 		return true
	// 	}
	// }

	if _, exists := data[email]; exists {
		return true
	}

	return false
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
}
