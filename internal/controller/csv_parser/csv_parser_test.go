package csv_parser

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type CSVParserTestSuite struct {
	suite.Suite
	parser CSVParser
}

func (suite *CSVParserTestSuite) SetupSuite() {
	suite.parser = NewCSVParser()
}

func (suite *CSVParserTestSuite) TestToCSVWithoutChunk() {
	data, err := suite.parser.ToCSV("../../../testdata", "../../../testdata", "roster1.csv")
	suite.Nil(err)
	suite.NotEmpty(data)
	suite.Equal(6, len(data))
}

func (suite *CSVParserTestSuite) TestToCSVWithChunk() {
	data, err := suite.parser.ToCSV("../../../testdata", "../../../testdata", "large1M.csv")
	suite.Nil(err)
	suite.NotEmpty(data)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(CSVParserTestSuite))
}
