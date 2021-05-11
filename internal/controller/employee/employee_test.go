package employee

import (
	"testing"

	"github.com/danieeelfr/csv_extractor_cli/internal/domain"
	"github.com/stretchr/testify/suite"
)

type EmployeeTestSuite struct {
	suite.Suite
	emp Employee
}

func (suite *EmployeeTestSuite) SetupSuite() {
	suite.emp = NewEmployee()
}

func (suite *EmployeeTestSuite) TestToCSVString() {
	data := make([]*domain.Employee, 0)
	e1 := domain.Employee{
		Id:     "50",
		Name:   "everson",
		Email:  "etv@tester.com",
		Salary: 1000,
	}
	e2 := domain.Employee{
		Id:     "1",
		Name:   "qualyeverson",
		Email:  "qualy@tester.com",
		Salary: 2000,
	}
	data = append(data, &e1, &e2)
	out := suite.emp.ToCSVString(data)
	suite.NotEmpty(out)
	suite.Equal(3, len(out))
}

func (suite *EmployeeTestSuite) TestToCSVStringRecordWithoutAllMandatoryFieldsShouldBeIgnored() {
	data := make([]*domain.Employee, 0)
	e1 := domain.Employee{}
	e2 := domain.Employee{}
	e3 := domain.Employee{
		Id:     "2",
		Name:   "qualyeterson",
		Email:  "qualyter@tester.com",
		Salary: 500,
	}
	e4 := domain.Employee{}
	data = append(data, &e1, &e2, &e3, &e4)
	out := suite.emp.ToCSVString(data)
	suite.NotEmpty(out)
	suite.Equal(2, len(out))
}
func (suite *EmployeeTestSuite) TestToCSVStringWithoutRecordsShouldReturnEmpty() {
	data := make([]*domain.Employee, 0)
	out := suite.emp.ToCSVString(data)
	suite.Empty(out)
	suite.Equal(0, len(out))
}

func (suite *EmployeeTestSuite) TestSplitValidAndInvalidRecords() {
	data := make([]*domain.Employee, 0)

	invalid := domain.Employee{
		Name: "testerrr",
	}
	valid := domain.Employee{
		Id:     "150",
		Name:   "qualytester",
		Email:  "tester@tester.com",
		Salary: 8000,
	}
	data = append(data, &valid, &invalid)
	v, i := suite.emp.SplitValidAndInvalidRecords(data)
	suite.NotEmpty(v, i)
	suite.Equal(1, len(v))
	suite.Equal(1, len(i))
}

func (suite *EmployeeTestSuite) TestDedupe() {
	data := make([]*domain.Employee, 0)
	e1 := domain.Employee{
		Id:     "23",
		Name:   "qualyeterson",
		Email:  "repeated@rep.com",
		Salary: 2000,
	}
	e2 := domain.Employee{
		Id:     "2150",
		Name:   "qtname",
		Email:  "repeated@rep.com",
		Salary: 500,
	}
	e3 := domain.Employee{
		Id:     "2150",
		Name:   "qtname",
		Email:  "repeated2@rep.com",
		Salary: 500,
	}
	e4 := domain.Employee{
		Id:     "2150",
		Name:   "qtname",
		Email:  "repeated2@rep.com",
		Salary: 500,
	}
	data = append(data, &e1, &e2, &e3, &e4)
	deduplicated := suite.emp.Dedupe(data)
	suite.NotEmpty(deduplicated)
	suite.Equal(2, len(deduplicated))
}

func (suite *EmployeeTestSuite) TestFetchFromStringCSV() {
	data := make([][]string, 0)
	data = append(data, []string{"Id", "Name", "Email", "Wage"})
	data = append(data, []string{"1", "Clouderson", "clouderson@cloud.com", "5000.50"})

	out := suite.emp.FetchFromStringCSV(data)
	suite.NotEmpty(out)
	suite.Equal(1, len(out))
	
	emp := out[0]
	suite.Equal("1", emp.Id)
	suite.Equal("Clouderson", emp.Name)
	suite.Equal("clouderson@cloud.com", emp.Email)
	suite.Equal(5000.50, emp.Salary)
}

func TestEmployeeTestSuite(t *testing.T) {
	suite.Run(t, new(EmployeeTestSuite))
}
