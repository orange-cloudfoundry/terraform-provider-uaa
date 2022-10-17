package providertest

import (
	"github.com/terraform-providers/terraform-provider-uaa/uaa/provider"
	"testing"
)

func TestProvider(t *testing.T) {
	if err := provider.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = provider.Provider()
}
