package alkira

import (
	"fmt"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/stretchr/testify/require"
)

func TestSegmentGenerateSegmentRequest(t *testing.T) {
	r := resourceAlkiraSegment()
	d := r.TestResourceData()

	// Test with multiple CIDR
	expectedAsn := 2
	expectedName := "testName"
	expectedReservePublicIPs := true
	expectedCidrs := []string{"10.255.254.0/24", "10.255.255.0/24"}

	d.Set("asn", expectedAsn)
	d.Set("name", expectedName)
	d.Set("reserve_public_ips", expectedReservePublicIPs)
	d.Set("cidrs", expectedCidrs)

	s, err := generateSegmentRequest(d)

	require.NoError(t, err)
	require.Equal(t, expectedAsn, s.Asn)
	require.Equal(t, expectedName, s.Name)
	require.Equal(t, expectedReservePublicIPs, s.ReservePublicIPsForUserAndSiteConnectivity)
	require.Equal(t, "", s.IpBlock) // should be empty because IpBlocks can be used even for single CIDR values
	require.Equal(t, len(expectedCidrs), len(s.IpBlocks.Values))

	// Test with single CIDR
	expectedCidr := "10.255.255.0/24"
	expectedCidrs = []string{expectedCidr}
	d.Set("cidrs", expectedCidrs)

	s, err = generateSegmentRequest(d)
	require.NoError(t, err)
	require.Equal(t, "", s.IpBlock)
	require.Equal(t, 1, len(s.IpBlocks.Values))
}

func TestSegmentSetCidrSegmentReadEmptyIpBlock(t *testing.T) {
	r := resourceAlkiraSegment()
	d := r.TestResourceData()

	expectedValues := []string{"a", "b", "c"}
	s := alkira.Segment{
		IpBlock: "",
		IpBlocks: alkira.SegmentIpBlocks{
			Values: expectedValues,
		},
	}

	setCidrsSegmentRead(d, &s)

	c := convertTypeListToStringList(d.Get("cidrs").([]interface{}))

	require.Equal(t, len(expectedValues), len(c))
	fmt.Println(c)
}

func TestSegmentSetCidrSegmentReadIpBlockContainedIpBlocks(t *testing.T) {
	r := resourceAlkiraSegment()
	d := r.TestResourceData()

	expectedIpBlock := "a"
	expectedValues := []string{expectedIpBlock, "b", "c"}
	s := alkira.Segment{
		IpBlock: expectedIpBlock,
		IpBlocks: alkira.SegmentIpBlocks{
			Values: expectedValues,
		},
	}

	setCidrsSegmentRead(d, &s)

	c := convertTypeListToStringList(d.Get("cidrs").([]interface{}))

	require.Equal(t, len(expectedValues), len(c))
	fmt.Println(c)
}

func TestSetCidrSegmentReadIpBlockAndIpBlocksPopulated(t *testing.T) {
	r := resourceAlkiraSegment()
	d := r.TestResourceData()

	expectedIpBlock := "d"
	expectedValues := []string{"a", "b", "c"}
	s := alkira.Segment{
		IpBlock: expectedIpBlock,
		IpBlocks: alkira.SegmentIpBlocks{
			Values: expectedValues,
		},
	}

	setCidrsSegmentRead(d, &s)

	c := convertTypeListToStringList(d.Get("cidrs").([]interface{}))

	require.Equal(t, len(expectedValues)+1, len(c))
}
