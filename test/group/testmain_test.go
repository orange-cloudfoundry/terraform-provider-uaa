package group

import (
	"github.com/terraform-providers/terraform-provider-uaa/test/util"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	uaaConfigPath, _ := filepath.Abs("../../test/resources/uaa.yml")
	util.RunIntegrationTests(m, uaaConfigPath)
}
