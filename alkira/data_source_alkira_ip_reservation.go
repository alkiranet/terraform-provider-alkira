package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraIpReservation() *schema.Resource {
	return &schema.Resource{
		Description: "This data source allows to retrieve an existing IP Reservation by its name.",
		Read:        dataSourceAlkiraIpReservationRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the IP Reservation.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceAlkiraIpReservationRead(d *schema.ResourceData, m interface{}) error {

	// INIT
	api := alkira.NewIPReservation(m.(*alkira.AlkiraClient))

	reservation, _, err := api.GetByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(string(reservation.Id))

	return nil
}
