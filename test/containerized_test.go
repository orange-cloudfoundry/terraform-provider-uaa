//go:build containerized

package test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// This function executes before the test suite begins execution
func (suite *integrationTestSuite) SetupSuite() {

	testManager := &integrationTestManager{context: context.Background()}
	testManager.prepareEnvironment()
	testManager.prepareDbContainer()
	testManager.prepareUaaContainer()
	testManager.createTestIdentityZone()

	// We would defer cleanup here if we didn't want to rely on garbage collection from testcontainer's reaper (ryuk)
	// defer testManager.dbContainer.Terminate(uaaTestManager.context)
	// defer testManager.uaaContainer.Terminate(uaaTestManager.context)
}

// This function executes after each test case
func (suite *integrationTestSuite) TearDownTest() {
	// Nothing to do; garbage collection will cleanup containers
}

// In order for 'go test' to run this suite, we need to create a normal test function and pass our suite to suite.Run
func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(integrationTestSuite))
}

// All methods that begin with "Test" are run as tests within a suite.

func (suite *integrationTestSuite) TestIntegrationTests() {
	cmd := exec.Command("go", "test", "-v", "./...")
	cmd.Env = os.Environ()
	stdout, err := cmd.Output()

	log.Print("\n", string(stdout))

	if err != nil {
		assert.Failf(suite.T(), "Error running integration tests", err.Error())
	}
}

// Private structures, methods, etc. for env setup
type integrationTestSuite struct {
	suite.Suite
}

type integrationTestManager struct {
	context      context.Context
	dbContainer  testcontainers.Container
	uaaContainer testcontainers.Container
}

func (uaaTestManager *integrationTestManager) prepareEnvironment() {
	log.Println("Preparing test environment...")

	os.Setenv("TF_ACC", "true")
	os.Setenv("DB_DATABASE", "uaa")
	os.Setenv("DB_USERNAME", "uaa")
	os.Setenv("DB_PASSWORD", "password")
	os.Setenv("UAA_CONFIG_PATH", "/uaa")
}

func (uaaTestManager *integrationTestManager) prepareDbContainer() {
	log.Println("Preparing test db container...")

	req := testcontainers.ContainerRequest{
		Image: "postgres:14.5-alpine",
		Env: map[string]string{
			"POSTGRES_DB":       os.Getenv("DB_DATABASE"),
			"POSTGRES_USER":     os.Getenv("DB_USERNAME"),
			"POSTGRES_PASSWORD": os.Getenv("DB_PASSWORD"),
			"PGPASSWORD":        os.Getenv("DB_PASSWORD"),
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
	}

	uaaTestManager.dbContainer = uaaTestManager.prepareContainer(req)

	ip, err := uaaTestManager.dbContainer.ContainerIP(uaaTestManager.context)
	if err != nil {
		log.Fatal("Unable to determine DB container IP")
	}
	os.Setenv("DB_HOST", ip)
}

func (uaaTestManager *integrationTestManager) prepareUaaContainer() {
	log.Println("Preparing test uaa container...")

	uaaConfigPath, _ := filepath.Abs("../test/resources/uaa.yml")
	req := testcontainers.ContainerRequest{
		Image:        "cloudfoundry/uaa:76.0.0",
		ExposedPorts: []string{"8080/tcp"},
		WaitingFor:   wait.ForLog("Server startup in").WithStartupTimeout(time.Minute * 3),
		Env: map[string]string{
			"DB_HOST":         os.Getenv("DB_HOST"),
			"DB_DATABASE":     os.Getenv("DB_DATABASE"),
			"DB_USERNAME":     os.Getenv("DB_USERNAME"),
			"DB_PASSWORD":     os.Getenv("DB_PASSWORD"),
			"UAA_CONFIG_PATH": os.Getenv("UAA_CONFIG_PATH"),
		},
		Mounts: testcontainers.Mounts(
			testcontainers.BindMount(uaaConfigPath, "/uaa/uaa.yml"),
		),
	}

	uaaTestManager.uaaContainer = uaaTestManager.prepareContainer(req)

	endpoint, err := uaaTestManager.uaaContainer.Endpoint(uaaTestManager.context, "8080")
	if err != nil {
		log.Fatal("Unable to determine UAA endpoint URL")
	}
	endpoint = "http://" + strings.Split(endpoint, "://")[1]

	os.Setenv("UAA_AUTH_URL", endpoint)
	os.Setenv("UAA_LOGIN_URL", endpoint)
	os.Setenv("UAA_CLIENT_ID", "admin")
	os.Setenv("UAA_CLIENT_SECRET", "adminsecret")
	os.Setenv("UAA_SKIP_SSL_VALIDATION", "1")
}

func (uaaTestManager *integrationTestManager) createTestIdentityZone() {
	log.Println("Creating additional test identity zone...")

	cmd := []string{
		"psql",
		"-U", os.Getenv("DB_USERNAME"),
		"-d", os.Getenv("DB_DATABASE"),
		"-c", "insert into identity_zone (id, name, subdomain) values ('test-zone', 'Test Zone', 'testzone');",
	}

	_, _, err := uaaTestManager.dbContainer.Exec(uaaTestManager.context, cmd)

	if err != nil {
		log.Println("Unable to create test zone; Expect test failures for tests that depend on that zone.")
	}
}

func (uaaTestManager *integrationTestManager) prepareContainer(req testcontainers.ContainerRequest) (container testcontainers.Container) {
	container, err := testcontainers.GenericContainer(
		uaaTestManager.context,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
	if err != nil {
		log.Println(fmt.Sprintf("Error: %s", err))
		log.Fatal("Unable to start container using image: " + req.Image)
	}

	return container
}
