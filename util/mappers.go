package util

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func ToInterface(data []string) (res []interface{}) {
	for _, val := range data {
		res = append(res, val)
	}
	return
}

func ToStringsSlice(data interface{}) (res []string) {
	for _, val := range data.(*schema.Set).List() {
		res = append(res, val.(string))
	}
	return
}
