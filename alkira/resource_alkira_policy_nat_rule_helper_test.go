package alkira

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestPolicyNatRuleExpandMatchArrays(t *testing.T) {
	expectedStrArr := []interface{}{"1", "2", "3"}
	expectedIntArr := []interface{}{1, 2, 3}

	m := makeMapPolicyNatRuleMatch(expectedStrArr, expectedIntArr)
	mArr := []interface{}{m}

	r := resourceAlkiraPolicyNatRule()
	f := schema.HashResource(r)
	s := schema.NewSet(f, mArr)

	actual := expandPolicyNatRuleMatch(s)

	require.Equal(t, convertTypeListToStringList(expectedStrArr), actual.SourcePrefixes)
	require.Equal(t, convertTypeListToIntList(expectedIntArr), actual.SourcePrefixListIds)
	require.Equal(t, convertTypeListToStringList(expectedStrArr), actual.SourcePortList)

	require.Equal(t, convertTypeListToStringList(expectedStrArr), actual.DestPrefixes)
	require.Equal(t, convertTypeListToIntList(expectedIntArr), actual.DestPrefixListIds)
	require.Equal(t, convertTypeListToStringList(expectedStrArr), actual.DestPortList)
}

func TestPolicyNatRuleExpandNatRuleActionArrays(t *testing.T) {
	expectedStrArr := []interface{}{"1", "2", "3"}
	expectedIntArr := []interface{}{1, 2, 3}

	m := makeMapPolicyNatRuleAction(expectedStrArr, expectedIntArr)
	mArr := []interface{}{m}

	r := resourceAlkiraPolicyNatRule()
	f := schema.HashResource(r)
	s := schema.NewSet(f, mArr)

	actual := expandPolicyNatRuleAction(s)

	//Src array validations
	require.Equal(t,
		convertTypeListToStringList(expectedStrArr),
		actual.SourceAddressTranslation.TranslatedPrefixes,
	)
	require.Equal(t,
		convertTypeListToIntList(expectedIntArr),
		actual.SourceAddressTranslation.TranslatedPrefixListIds,
	)

	//Dst array validations
	require.Equal(t,
		convertTypeListToStringList(expectedStrArr),
		actual.DestinationAddressTranslation.TranslatedPrefixes,
	)
	require.Equal(t,
		convertTypeListToStringList(expectedStrArr),
		actual.DestinationAddressTranslation.TranslatedPortList,
	)
	require.Equal(t,
		convertTypeListToIntList(expectedIntArr),
		actual.DestinationAddressTranslation.TranslatedPrefixListIds,
	)

}

//
// TEST HELPERS
//

func makeMapPolicyNatRuleAction(strArr, intArr []interface{}) map[string]interface{} {
	m := make(map[string]interface{})

	m["src_addr_translation_prefixes"] = strArr
	m["src_addr_translation_prefix_list_ids"] = intArr

	m["dst_addr_translation_prefixes"] = strArr
	m["dst_addr_translation_prefix_list_ids"] = intArr
	m["dst_addr_translation_ports"] = strArr

	return m
}

func makeMapPolicyNatRuleMatch(strArr, intArr []interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	m["src_prefixes"] = strArr
	m["src_prefix_list_ids"] = intArr
	m["src_ports"] = strArr
	m["dst_prefixes"] = strArr
	m["dst_prefix_list_ids"] = intArr
	m["dst_ports"] = strArr
	m["protocol"] = "ANY"

	return m
}
