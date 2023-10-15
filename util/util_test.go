package util

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

var GlobalTestGeneration = 10

type UtilTest struct {
	isintable []IsInTable
	suite.Suite
}

type IsInTable struct {
	bee      string
	hive     []string
	expected bool
}

// generate some hashes with writes to the datastore
func InitIsInTables() []IsInTable {
	var dump []IsInTable
	for i := 0; i < GlobalTestGeneration; i++ {
		// gen attributes
	}
	return dump
}

func (s *UtilTest) SetupTest() {
	fmt.Println("Starting tests...")

	fmt.Println("Tests startup complete...")
}

// TestIsIn tests that IsIn function works as expected
func (s *UtilTest) TestIsIn() {
	s.isintable = InitIsInTables()
	if len(s.isintable) == 0 {
		s.isintable = []IsInTable{
			{
				bee:      "demi",
				hive:     []string{"demi", "sansa", "kebble", "juice"},
				expected: true,
			},
		}
	}

	for i := 0; i < len(s.isintable); i++ {
		s.Assert().Equal(s.isintable[i].expected, IsIn(s.isintable[i].bee, s.isintable[i].hive))
	}
}

//

func (s *UtilTest) TearDownSuite() {
	fmt.Println("Commencing test cleanup")
	//err := cleanUpAfterCatTest()
	//s.Require().NoError(err)
	fmt.Println("All testing complete")
}

func TestUtilTest(t *testing.T) {
	suite.Run(t, new(UtilTest))
}

//func cleanUpAfterCatTest() error {
//	err := cleanUpAfterTest()
//	// cat content
//	return err
//}
