package alkira

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/k0kubun/pp"
)

func TestExpandPolicyNatRule(t *testing.T) {
	m := make(map[string]interface{})
	m["src_addr_translation_type"] = "STATIC_IP"
	m["src_addr_translation_prefixes"] = []string{"10.10.10.10/32"}
	m["src_addr_translation_bidirectional"] = true
	m["src_addr_translation_match_and_invalidate"] = true
	mArr := []interface{}{m}

	r := resourceAlkiraPolicyNatRule()
	f := schema.HashResource(r)
	s := schema.NewSet(f, mArr)

	n := expandPolicyNatRuleAction(s)

	pp.Println(n.SourceAddressTranslation)
}

func TestExpandPolicyNatMatch(t *testing.T) {
	m := make(map[string]interface{})
	m["src_prefixes"] = []string{"10.10.10.10/32"}
	m["dst_prefixes"] = []string{"any"}
	m["protocol"] = "any"
	mArr := []interface{}{m}

	r := resourceAlkiraPolicyNatRule()
	f := schema.HashResource(r)
	s := schema.NewSet(f, mArr)

	n := expandPolicyNatRuleMatch(s)

	pp.Println(n)
}
