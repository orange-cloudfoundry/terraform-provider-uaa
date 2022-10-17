package util

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/provider"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/uaaapi"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"strings"
	"testing"
)

// TODO: figure out how to get uaa config path dependably and avoid passing it in
// TODO: figure out how to use same containers when running tests for multiple packages

var IntegrationTestManager *UaaIntegrationTestManager

func RunIntegrationTests(m *testing.M, uaaConfigPath string) {
	IntegrationTestManager = newIntegrationTestManager(uaaConfigPath)

	log.Println("Running tests...")
	exitCode := m.Run()

	IntegrationTestManager.destroy()

	os.Exit(exitCode)
}

type UaaIntegrationTestManager struct {
	context           context.Context
	dbContainer       testcontainers.Container
	uaaContainer      testcontainers.Container
	uaaProvider       *schema.Provider
	ProviderFactories map[string]func() (*schema.Provider, error)
}

func newIntegrationTestManager(uaaConfigPath string) *UaaIntegrationTestManager {
	ctx := context.Background()

	uaaTestManager := &UaaIntegrationTestManager{context: ctx}
	uaaTestManager.prepareEnvironment()
	uaaTestManager.prepareDbContainer()
	defer uaaTestManager.dbContainer.Terminate(ctx)
	uaaTestManager.prepareUaaContainer(uaaConfigPath)
	defer uaaTestManager.uaaContainer.Terminate(ctx)
	uaaTestManager.createTestIdentityZone()
	uaaTestManager.prepareProviderFactories()

	return uaaTestManager
}

func (uaaTestManager *UaaIntegrationTestManager) prepareEnvironment() {
	log.Println("Preparing test environment...")

	os.Setenv("DB_DATABASE", "uaa")
	os.Setenv("DB_USERNAME", "uaa")
	os.Setenv("DB_PASSWORD", "password")
	os.Setenv("UAA_CONFIG_PATH", "/uaa")
}

func (uaaTestManager *UaaIntegrationTestManager) prepareDbContainer() {
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

func (uaaTestManager *UaaIntegrationTestManager) prepareUaaContainer(uaaConfigPath string) {
	log.Println("Preparing test uaa container...")

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
	os.Setenv("TF_ACC", "1")
	os.Setenv("UAA_SKIP_SSL_VALIDATION", "1")
}

func (uaaTestManager *UaaIntegrationTestManager) createTestIdentityZone() {
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

func (uaaTestManager *UaaIntegrationTestManager) prepareContainer(req testcontainers.ContainerRequest) (container testcontainers.Container) {
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

func (uaaTestManager *UaaIntegrationTestManager) prepareProviderFactories() {
	log.Println("Preparing provider factories...")

	uaaTestManager.uaaProvider = provider.Provider()

	uaaTestManager.ProviderFactories = map[string]func() (*schema.Provider, error){
		"uaa": func() (*schema.Provider, error) {
			return uaaTestManager.uaaProvider, nil
		},
	}
}

func (uaaTestManager *UaaIntegrationTestManager) UaaSession() *uaaapi.Session {
	return uaaTestManager.uaaProvider.Meta().(*uaaapi.Session)
}
