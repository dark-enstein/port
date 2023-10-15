package auth

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"strconv"
	"testing"
)

var GlobalTestGeneration = 10

type UserTest struct {
	permissions TestPermission
	suite.Suite
}

// TestPermission defines sample permisions and their expected binary representstion
type TestPermission struct {
	permissions []Permission
	expected    []string
}

// InitPermissionToBinary initializes the TestPermission struct
func InitPermissionToBinary() TestPermission {
	var dump TestPermission
	for i := 0; i < GlobalTestGeneration; i++ {
		// gen attributes
	}
	return dump
}

func (s *UserTest) SetupTest() {
	fmt.Println("Starting tests...")

	fmt.Println("Tests startup complete...")
}

// TestIsIn tests that IsIn function works as expected
func (s *UserTest) TestPermission() {
	s.permissions = InitPermissionToBinary()
	if len(s.permissions.permissions) == 0 {
		s.permissions = TestPermission{
			permissions: []Permission{
				{
					Name:   "User",
					Create: true,
					Read:   true,
					Update: false,
					Delete: true,
				},
			},
			expected: []string{
				"1101",
			},
		}
	}

	for i := 0; i < len(s.permissions.permissions); i++ {
		binAct := *s.permissions.permissions[i].ToBinary()
		expBin, err := strconv.ParseInt(s.permissions.expected[i], 2, 64)
		if err != nil {
			return
		}
		s.Assert().Equal(expBin, binAct)
	}
}

//

func (s *UserTest) TearDownSuite() {
	fmt.Println("Commencing test cleanup")
	//err := cleanUpAfterCatTest()
	//s.Require().NoError(err)
	fmt.Println("All testing complete")
}

func TestUtilTest(t *testing.T) {
	suite.Run(t, new(UserTest))
}

//func cleanUpAfterCatTest() error {
//	err := cleanUpAfterTest()
//	// cat content
//	return err
//}
