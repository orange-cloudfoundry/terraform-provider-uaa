package uaa

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// getChangedValueString
func getChangedValueString(key string, updated *bool, d *schema.ResourceData) *string {

	if d.HasChange(key) {
		vv := d.Get(key).(string)
		*updated = *updated || true
		return &vv
	} else if v, ok := d.GetOk(key); ok {
		vv := v.(string)
		return &vv
	}
	return nil
}

// getChangedValueInt
func getChangedValueInt(key string, updated *bool, d *schema.ResourceData) *int {

	if d.HasChange(key) {
		vv := d.Get(key).(int)
		*updated = *updated || true
		return &vv
	} else if v, ok := d.GetOk(key); ok {
		vv := v.(int)
		return &vv
	}
	return nil
}

// getChangedValueBool
func getChangedValueBool(key string, updated *bool, d *schema.ResourceData) *bool {

	if d.HasChange(key) {
		vv := d.Get(key).(bool)
		*updated = *updated || true
		return &vv
	} else if v, ok := d.GetOk(key); ok {
		vv := v.(bool)
		return &vv
	}
	return nil
}

// getChangedValueStringList
func getChangedValueStringList(key string, updated *bool, d *schema.ResourceData) *[]string {
	var a []interface{}

	if d.HasChange(key) {
		a = d.Get(key).(*schema.Set).List()
		*updated = *updated || true
	} else if v, ok := d.GetOk(key); ok {
		a = v.(*schema.Set).List()
	}
	if a != nil {
		aa := []string{}
		for _, vv := range a {
			aa = append(aa, vv.(string))
		}
		return &aa
	}
	return nil
}

// getResourceChange -
func getResourceChange(key string, d *schema.ResourceData) (bool, string, string) {
	old, new := d.GetChange(key)
	return old != new, old.(string), new.(string)
}

// getListChanges -
func getListChanges(old interface{}, new interface{}) (remove []string, add []string) {

	var a bool

	for _, o := range old.(*schema.Set).List() {
		remove = append(remove, o.(string))
	}
	for _, n := range new.(*schema.Set).List() {
		nn := n.(string)
		a = true
		for i, r := range remove {
			if nn == r {
				remove = append(remove[:i], remove[i+1:]...)
				a = false
				break
			}
		}
		if a {
			add = append(add, nn)
		}
	}
	return
}
