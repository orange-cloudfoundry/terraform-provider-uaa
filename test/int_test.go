package test

//
//import (
//	"github.com/stretchr/testify/suite"
//	usertest "github.com/terraform-providers/terraform-provider-uaa/test/user"
//	"github.com/terraform-providers/terraform-provider-uaa/test/util"
//	"testing"
//)
//
//type UaaIntegrationTestSuite struct {
//	suite.Suite
//}
//
//func (suite *UaaIntegrationTestSuite) SetupSuite() {
//	util.IntegrationTestManager = util.NewIntegrationTestManager()
//}
//
//func (suite *UaaIntegrationTestSuite) TearDownSuite() {
//	util.IntegrationTestManager.Destroy()
//}
//
//func TestUaaIntegrationTestSuite(t *testing.T) {
//	suite.Run(t, new(UaaIntegrationTestSuite))
//}
//
//func (suite *UaaIntegrationTestSuite) TestUserDataSource() {
//	usertest.TestUserDataSource(suite.T())
//}
