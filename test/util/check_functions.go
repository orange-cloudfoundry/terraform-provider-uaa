package util

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestCheckResourceSet(ref string, attr string, values []string) resource.TestCheckFunc {

	lTests := make([]resource.TestCheckFunc, 0)

	for i, cVal := range values {
		lKey := fmt.Sprintf("%s.%d", attr, i)
		lTests = append(lTests, resource.TestCheckResourceAttr(ref, lKey, cVal))
	}

	return resource.ComposeTestCheckFunc(lTests...)
}
