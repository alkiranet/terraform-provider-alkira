This release introduces new resource to manage Azure ExpressRoute
connector and several new data sources with various enhancements and
bug fixes.

~> **DEPRECATION** All deprecated standalone `credential` resources
have been removed.  Those credentials have been integrated directly
into connector or service resources.  Please migrate your config when
upgrading to this version.
* `alkira_credential_cisco_sdwan`
* `alkira_credential_checkpoint`
* `alkira_credential_fortinet`
* `alkira_credential_fortinet_instance`
* `alkira_credential_pan`
* `alkira_credential_pan_instance`

## Resources & Data Sources

#### resource `alkira_connector_azure_expressroute` (**NEW**)

* New connector support of Azure ExpressRoute.

#### resource `alkira_credential_ssh_key_pair` (**NEW**)

* New resource for storing SSH key pair credential.  For now, it could
  be used along with `alkira_connector_cisco_sdwan`.

#### ￼ resource `alkira_connector_aws_vpc`

* Add `failover_cxps` to support additional CXPs as failover.
* Update documentation.

#### resource `alkira_connector_azure_vnet`

* Add `failover_cxps` to support additional CXPs as failover.
* Update documentation.

#### resource `alkira_connector_cisco_sdwan`

* Maintain instance order of `alkira_connector_cisco_sdwan` .
* Set default value for `custom_asn`.
* `username` and `password` are added in `vedge` block for entering credentials.

#### resource `alkira_connector_gcp_vpc`

* Add `failover_cxps` to support additional CXPs as failover.

#### resource `alkira_connector_internet_exit`

* `size` is removed and it’s not supported in this connector.
* `byoip_id` is added for supporting BYOIP.

#### resource `alkira_connector_ipsec`

* ￼ Add `availability` to dynamic routing of `alkira_connector_ipsec`.

#### resource `alkira_connector_oci_vcn`

* Add `failover_cxps` to support additional CXPs as failover.
* Remove deprecated `primary`.
* Update documentation.

#### resource `policy_nat_rule`

* Add the default value of `action` block of `policy_nat_rule`.
* Update documentation.

#### resource `alkira_service_fortinet`

* Add `license_key_path` to pass license key from a  file.
* Make  `license_key` and `license_key_path` optional for `PAY_AS_YOU_GO` license.

#### resource `alkira_service_infoblox`

* Only `BYOL` is allowed for license type.
* Update documentation.

#### resource `alkira_service_pan`

* Mark `pan_username` and `pan_password` required in `alkira_service_pan`.
* Fix the random string in `alkira_serivce_pan` instance name.


## Enhancments

* Add unique request ID for every single API calls for easier debugging.
