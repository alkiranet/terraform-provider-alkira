package alkira

import (
	"log"
	"strings"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandPolicyNatRuleMatch(in *schema.Set) *alkira.NatRuleMatch {
	if in == nil || in.Len() == 0 || in.Len() > 1 {
		log.Printf("[ERROR] invalid match section (%d)", in.Len())
		return nil
	}

	match := alkira.NatRuleMatch{}

	for _, m := range in.List() {

		matchValue := m.(map[string]interface{})

		if v, ok := matchValue["src_prefixes"].([]interface{}); ok {
			match.SourcePrefixes = convertTypeListToStringList(v)
		}
		if v, ok := matchValue["src_prefix_list_ids"].([]interface{}); ok {
			match.SourcePrefixListIds = convertTypeListToIntList(v)
		}
		if v, ok := matchValue["src_ports"].([]interface{}); ok {
			match.SourcePortList = convertTypeListToStringList(v)
		}
		if v, ok := matchValue["dst_prefixes"].([]interface{}); ok {
			match.DestPrefixes = convertTypeListToStringList(v)
		}
		if v, ok := matchValue["dst_prefix_list_ids"].([]interface{}); ok {
			match.DestPrefixListIds = convertTypeListToIntList(v)
		}
		if v, ok := matchValue["dst_ports"].([]interface{}); ok {
			match.DestPortList = convertTypeListToStringList(v)
		}
		if v, ok := matchValue["protocol"].(string); ok {
			match.Protocol = strings.ToLower(v)
		}
	}

	return &match
}

// expandPolicyNatRuleAction expand "action" section
func expandPolicyNatRuleAction(in *schema.Set) *alkira.NatRuleAction {
	if in == nil || in.Len() == 0 || in.Len() > 1 {
		log.Printf("[ERROR] invalid action section (%d)", in.Len())
		return nil
	}

	st := alkira.NatRuleActionSrcTranslation{}
	dt := alkira.NatRuleActionDstTranslation{}

	for _, m := range in.List() {

		actionValue := m.(map[string]interface{})

		if v, ok := actionValue["src_addr_translation_type"].(string); ok {
			st.TranslationType = v
		}
		if v, ok := actionValue["src_addr_translation_prefixes"].([]interface{}); ok {
			st.TranslatedPrefixes = convertTypeListToStringList(v)
		}
		if v, ok := actionValue["src_addr_translation_prefix_list_ids"].([]interface{}); ok {
			st.TranslatedPrefixListIds = convertTypeListToIntList(v)
		}
		if v, ok := actionValue["src_addr_translation_bidirectional"].(bool); ok {
			st.Bidirectional = v
		}
		if v, ok := actionValue["src_addr_translation_match_and_invalidate"].(bool); ok {
			st.MatchAndInvalidate = v
		}
		if v, ok := actionValue["dst_addr_translation_type"].(string); ok {
			dt.TranslationType = v
		}
		if v, ok := actionValue["dst_addr_translation_prefixes"].([]interface{}); ok {
			dt.TranslatedPrefixes = convertTypeListToStringList(v)
		}
		if v, ok := actionValue["dst_addr_translation_prefix_list_ids"].([]interface{}); ok {
			dt.TranslatedPrefixListIds = convertTypeListToIntList(v)
		}
		if v, ok := actionValue["dst_addr_translation_ports"].([]interface{}); ok {
			dt.TranslatedPortList = convertTypeListToStringList(v)
		}
		if v, ok := actionValue["dst_addr_translation_bidirectional"].(bool); ok {
			dt.Bidirectional = v
		}
		if v, ok := actionValue["dst_addr_translation_advertise_to_connector"].(bool); ok {
			dt.AdvertiseToConnector = v
		}
	}

	action := alkira.NatRuleAction{
		SourceAddressTranslation:      st,
		DestinationAddressTranslation: dt,
	}

	return &action
}
