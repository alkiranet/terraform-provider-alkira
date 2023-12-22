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
	e := alkira.EgressAction{}

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
		if v, ok := actionValue["src_addr_translation_match_and_invalidate"].(bool); ok {
			//
			// This flag is only available when TranslationType is not
			// "NONE". Otherwise, API will fail with
			// validation. However, the default value is "true".
			//
			if st.TranslationType != "NONE" {
				matchAndInvalidate := new(bool)
				*matchAndInvalidate = v
				st.MatchAndInvalidate = matchAndInvalidate
			}
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
		if v, ok := actionValue["dst_addr_translation_list_policy_fqdn_id"].(int); ok {
			dt.TranslatedPolicyFqdnListId = v
		}
		if v, ok := actionValue["dst_addr_translation_advertise_to_connector"].(bool); ok {
			//
			// This flag is only available when TranslationType is not
			// "NONE". Otherwise, API will fail with
			// validation. However, the default value is "true".
			//
			if dt.TranslationType != "NONE" {
				t := new(bool)
				*t = v
				dt.AdvertiseToConnector = t
			}
		}
		if v, ok := actionValue["egress_type"].(string); ok {
			e.IpType = v
		}

		if v, ok := actionValue["src_addr_translation_routing_track_prefixes"].([]interface{}); ok {
			list := convertTypeListToStringList(v)
			if len(list) > 0 {
				st.RoutingOptions.TrackPrefixes = list
			}
		}
		if v, ok := actionValue["src_addr_translation_routing_track_prefix_list_ids"].([]interface{}); ok {
			list := convertTypeListToIntList(v)
			if len(list) > 0 {
				st.RoutingOptions.TrackPrefixListIds = list
			}
		}
		if v, ok := actionValue["src_addr_translation_routing_track_invalidate_prefixes"].(bool); ok {
			invalidateRoutingTrackPrefixes := new(bool)
			*invalidateRoutingTrackPrefixes = v
			st.RoutingOptions.InvalidateRoutingTrackPrefixes = invalidateRoutingTrackPrefixes
		}

		//
		// This field has a fixed value based on the translation type.
		//
		if st.TranslationType == "STATIC_IP" {
			st.Bidirectional = func() *bool { b := true; return &b }()
		} else if st.TranslationType == "DYNAMIC_IP" {
			st.Bidirectional = func() *bool { b := false; return &b }()
		} else {
		}

		if dt.TranslationType == "STATIC_IP" {
			dt.Bidirectional = func() *bool { b := true; return &b }()
		} else if dt.TranslationType == "STATIC_IP_AND_PORT" {
			dt.Bidirectional = func() *bool { b := true; return &b }()
		} else if dt.TranslationType == "STATIC_PORT" {
			dt.Bidirectional = func() *bool { b := true; return &b }()
		} else {
		}
	}

	if len(st.RoutingOptions.TrackPrefixes) == 0 && len(st.RoutingOptions.TrackPrefixListIds) == 0 {
		st.RoutingOptions.InvalidateRoutingTrackPrefixes = nil
	}
	if st.RoutingOptions.InvalidateRoutingTrackPrefixes != nil {
		st.MatchAndInvalidate = nil
	}

	action := alkira.NatRuleAction{
		SourceAddressTranslation:      st,
		DestinationAddressTranslation: dt,
		Egress:                        e,
	}

	return &action
}

func setNatRuleActionOptions(a alkira.NatRuleAction, d *schema.ResourceData) {
	var action []map[string]interface{}

	in := map[string]interface{}{
		"src_addr_translation_routing_track_prefixes":            a.SourceAddressTranslation.RoutingOptions.TrackPrefixes,
		"src_addr_translation_routing_track_prefix_list_ids":     a.SourceAddressTranslation.RoutingOptions.TrackPrefixListIds,
		"src_addr_translation_routing_track_invalidate_prefixes": a.SourceAddressTranslation.RoutingOptions.InvalidateRoutingTrackPrefixes,
	}

	action = append(action, in)

	d.Set("action", action)
}
