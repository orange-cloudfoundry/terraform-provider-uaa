package uaa

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	uaaTestManager := NewUaaTestManager()

	log.Println("Running tests...")
	exitCode := m.Run()

	uaaTestManager.Destroy()

	os.Exit(exitCode)
}

type UaaTestManager struct {
	context      context.Context
	dbContainer  testcontainers.Container
	uaaContainer testcontainers.Container
}

func NewUaaTestManager() *UaaTestManager {
	context := context.Background()

	uaaTestManager := &UaaTestManager{context: context}
	uaaTestManager.PrepareEnvironment()
	uaaTestManager.PrepareDbContainer()
	uaaTestManager.PrepareUaaContainer()
	uaaTestManager.CreateTestIdentityZone()

	return uaaTestManager
}

func (uaaTestManager *UaaTestManager) PrepareEnvironment() {
	log.Println("Preparing test environment...")

	os.Setenv("DB_DATABASE", "uaa")
	os.Setenv("DB_USERNAME", "uaa")
	os.Setenv("DB_PASSWORD", "password")
	os.Setenv("UAA_CONFIG_PATH", "/uaa")
}

func (uaaTestManager *UaaTestManager) PrepareDbContainer() {
	log.Println("Preparing test db container...")
	//uaaTestManager.dbContainer = preparePostgresContainer(uaaTestManager.context)

	req := testcontainers.ContainerRequest{
		Image: "postgres:14.5-alpine",
		//ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       os.Getenv("DB_DATABASE"),
			"POSTGRES_USER":     os.Getenv("DB_USERNAME"),
			"POSTGRES_PASSWORD": os.Getenv("DB_PASSWORD"),
			"PGPASSWORD":        os.Getenv("DB_PASSWORD"),
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
	}

	uaaTestManager.dbContainer = uaaTestManager.PrepareContainer(req)

	ip, err := uaaTestManager.dbContainer.ContainerIP(uaaTestManager.context)
	if err != nil {
		log.Fatal("Unable to determine DB container IP")
	}
	os.Setenv("DB_HOST", ip)
}

func (uaaTestManager *UaaTestManager) PrepareUaaContainer() {
	log.Println("Preparing test uaa container...")

	configPath, err := filepath.Abs("../testresources/uaa.yml")
	if err != nil {
		log.Fatal("Unable to load UAA config file")
	}

	req := testcontainers.ContainerRequest{
		Image:        "cloudfoundry/uaa:76.0.0",
		ExposedPorts: []string{"8080/tcp"},
		WaitingFor:   wait.ForLog("Server startup in"),
		Env: map[string]string{
			"DB_HOST":         os.Getenv("DB_HOST"),
			"DB_DATABASE":     os.Getenv("DB_DATABASE"),
			"DB_USERNAME":     os.Getenv("DB_USERNAME"),
			"DB_PASSWORD":     os.Getenv("DB_PASSWORD"),
			"UAA_CONFIG_PATH": os.Getenv("UAA_CONFIG_PATH"),
		},
		Mounts: testcontainers.Mounts(
			testcontainers.BindMount(configPath, "/uaa/uaa.yml"),
		),
	}

	uaaTestManager.uaaContainer = uaaTestManager.PrepareContainer(req)

	endpoint, err := uaaTestManager.uaaContainer.Endpoint(uaaTestManager.context, "8080")
	if err != nil {
		log.Fatal("Unable to determine UAA endpoint URL")
	}
	endpoint = "http://" + strings.Split(endpoint, "://")[1]

	os.Setenv("UAA_AUTH_URL", endpoint)
	os.Setenv("UAA_LOGIN_URL", endpoint)
	os.Setenv("UAA_CLIENT_ID", "admin")
	os.Setenv("UAA_CLIENT_SECRET", "adminsecret")
	os.Setenv("TF_ACC", "1")
	os.Setenv("UAA_SKIP_SSL_VALIDATION", "1")
}

func (uaaTestManager *UaaTestManager) CreateTestIdentityZone() {
	// TODO: can we create the additional identity zone via `uaa.yml` instead?

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

func (uaaTestManager *UaaTestManager) PrepareContainer(req testcontainers.ContainerRequest) (container testcontainers.Container) {
	container, err := testcontainers.GenericContainer(
		uaaTestManager.context,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
	if err != nil {
		log.Fatal("Unable to start container using image: " + req.Image)
	}

	return container
}

func (uaaTestManager *UaaTestManager) Destroy() {
	log.Println("Terminating docker containers...")

	defer uaaTestManager.uaaContainer.Terminate(uaaTestManager.context)
	defer uaaTestManager.dbContainer.Terminate(uaaTestManager.context)
}
