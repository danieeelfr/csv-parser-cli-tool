package file_parser

import (
	"fmt"
	"testing"

	"github.com/danieeelfr/csv_extractor_cli/internal/controller/employee"
	"github.com/danieeelfr/csv_extractor_cli/internal/domain"
	"github.com/stretchr/testify/suite"
)

type FileParserTestSuite struct {
	suite.Suite
	parser   FileParser
	employee employee.Employee
}

func (suite *FileParserTestSuite) SetupSuite() {
	suite.parser = NewFileParser()
	suite.employee = employee.NewEmployee()
}

func (suite *FileParserTestSuite) TestSearchFiles() {
	out, err := suite.parser.SearchFiles(fmt.Sprintf(`%s`, "../../../testdata"))
	suite.Nil(err)
	suite.NotEmpty(out)
	suite.GreaterOrEqual(len(out), 4)
}

func (suite *FileParserTestSuite) TestSearchFilesWithInvalidPathShouldReturnError() {
	out, err := suite.parser.SearchFiles(fmt.Sprintf(`%s`, "./invalid"))
	suite.NotNil(err)
	suite.Nil(out)
}

func (suite *FileParserTestSuite) TestCreate() {
	data := make([]*domain.Employee, 0)
	e1 := domain.Employee{
		Id:     "123",
		Name:   "testerson",
		Email:  "tst@test.com",
		Salary: 1100.50,
	}
	e2 := domain.Employee{
		Id:     "321",
		Name:   "qualyeverson",
		Email:  "qualy@tester.com",
		Salary: 100.50,
	}
	data = append(data, &e1, &e2)
	s := suite.employee.ToCSVString(data)
	err := suite.parser.Create(s, "../../../testdata", "testcreate.csv")
	suite.Nil(err)
}

func (suite *FileParserTestSuite) TestCreateWithInvalidPathShouldReturnError() {
	data := make([]*domain.Employee, 0)
	e1 := domain.Employee{}
	e2 := domain.Employee{
		Id:     "321",
		Name:   "qualyeverson",
		Email:  "qualy@tester.com",
		Salary: 100.50,
	}
	data = append(data, &e1, &e2)
	s := suite.employee.ToCSVString(data)
	err := suite.parser.Create(s, "./invalid", "testcreate.csv")
	suite.NotNil(err)
}

func TestFileParserTestSuite(t *testing.T) {
	suite.Run(t, new(FileParserTestSuite))
}
