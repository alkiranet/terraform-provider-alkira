// Copyright (C) 2026 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// BillingAccount is one of the tenant's Lumen billing account numbers
// (BANs).
type BillingAccount struct {
	CustomerNumber       string `json:"customerNumber"`
	BillingAccountNumber string `json:"billingAccountNumber"`
	BillingAccountName   string `json:"billingAccountName"`
}

// MultiCloudGatewayOnRampRegion is a Lumen on-ramp of a CXP's cloud provider
// region.
type MultiCloudGatewayOnRampRegion struct {
	Name           string `json:"name"`
	ProviderRegion string `json:"providerRegion"`
}

// billingAccountListResponse is one page of the billing-account list:
// always paginated, as multiCloudGatewayListResponse is for the MCGW list.
type billingAccountListResponse struct {
	Data       []BillingAccount `json:"data"`
	Pagination Pagination       `json:"pagination"`
}

// multiCloudGatewayOnRampRegionListResponse is one page of the MCGW
// on-ramp-region list.
type multiCloudGatewayOnRampRegionListResponse struct {
	Data       []MultiCloudGatewayOnRampRegion `json:"data"`
	Pagination Pagination                      `json:"pagination"`
}

// GetBillingAccounts lists the tenant's Lumen billing account numbers
// (BANs), walking the offset/limit pages as GetMultiCloudGatewayByName does
// for the MCGW list. The list isn't tenant-network-scoped: every network of
// the tenant sees the same rows.
func GetBillingAccounts(ac *AlkiraClient) ([]BillingAccount, error) {
	baseUri := fmt.Sprintf("%s/billing-accounts", ac.URI)

	var accounts []BillingAccount
	offset := 0

	for {
		uri := fmt.Sprintf("%s?offset=%d&limit=%d", baseUri, offset, multiCloudGatewayPageLimit)

		data, _, err := ac.get(uri)
		if err != nil {
			return nil, fmt.Errorf("GetBillingAccounts: failed to get page at offset %d: %w", offset, err)
		}

		var page billingAccountListResponse
		if err := json.Unmarshal(data, &page); err != nil {
			return nil, fmt.Errorf("GetBillingAccounts: failed to unmarshal page at offset %d: %w", offset, err)
		}

		accounts = append(accounts, page.Data...)
		offset += len(page.Data)

		if len(page.Data) == 0 || offset >= page.Pagination.Hits {
			break
		}
	}

	return accounts, nil
}

// GetMultiCloudGatewayOnRampRegions lists the Lumen on-ramp regions of cxp's
// cloud provider region, walking the offset/limit pages as
// GetMultiCloudGatewayByName does for the MCGW list.
func GetMultiCloudGatewayOnRampRegions(ac *AlkiraClient, cxp string) ([]MultiCloudGatewayOnRampRegion, error) {
	if len(cxp) == 0 {
		return nil, fmt.Errorf("GetMultiCloudGatewayOnRampRegions: invalid cxp")
	}

	baseUri := fmt.Sprintf("%s/tenantnetworks/%s/multi-cloud-gateway-on-ramp-regions", ac.URI, ac.TenantNetworkId)

	var regions []MultiCloudGatewayOnRampRegion
	offset := 0

	for {
		uri := fmt.Sprintf("%s?cxp=%s&offset=%d&limit=%d", baseUri, url.QueryEscape(cxp), offset, multiCloudGatewayPageLimit)

		data, _, err := ac.get(uri)
		if err != nil {
			return nil, fmt.Errorf("GetMultiCloudGatewayOnRampRegions: failed to get page at offset %d: %w", offset, err)
		}

		var page multiCloudGatewayOnRampRegionListResponse
		if err := json.Unmarshal(data, &page); err != nil {
			return nil, fmt.Errorf("GetMultiCloudGatewayOnRampRegions: failed to unmarshal page at offset %d: %w", offset, err)
		}

		regions = append(regions, page.Data...)
		offset += len(page.Data)

		if len(page.Data) == 0 || offset >= page.Pagination.Hits {
			break
		}
	}

	return regions, nil
}
