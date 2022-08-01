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
	expecedName := "testName"
	expectedReservePublicIPs := true
	expectedCidrs := []string{"10.255.254.0/24", "10.255.255.0/24"}

	d.Set("asn", expectedAsn)
	d.Set("name", expecedName)
	d.Set("reserve_public_ips", expectedReservePublicIPs)
	d.Set("cidrs", expectedCidrs)

	s, err := generateSegmentRequest(d)

	require.NoError(t, err)
	require.Equal(t, s.Asn, expectedAsn)
	require.Equal(t, s.Name, expecedName)
	require.Equal(t, s.ReservePublicIPsForUserAndSiteConnectivity, expectedReservePublicIPs)
	require.Equal(t, s.IpBlock, "") // should be empty becuase IpBlocks can be used even for signle CIDR values
	require.Equal(t, len(s.IpBlocks.Values), len(expectedCidrs))

	// Test with single CIDR
	expectedCidr := "10.255.255.0/24"
	expectedCidrs = []string{expectedCidr}
	d.Set("cidrs", expectedCidrs)

	s, err = generateSegmentRequest(d)
	require.NoError(t, err)
	require.Equal(t, s.IpBlock, "")
	require.Equal(t, len(s.IpBlocks.Values), 1)
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

	setCidrsSegmentRead(d, s)

	c := convertTypeListToStringList(d.Get("cidrs").([]interface{}))

	require.Equal(t, len(c), len(expectedValues))
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

	setCidrsSegmentRead(d, s)

	c := convertTypeListToStringList(d.Get("cidrs").([]interface{}))

	require.Equal(t, len(c), len(expectedValues))
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

	setCidrsSegmentRead(d, s)

	c := convertTypeListToStringList(d.Get("cidrs").([]interface{}))

	require.Equal(t, len(c), len(expectedValues)+1)
}
