package tailscale_test

import (
	"net/http"
	"testing"

	"github.com/davidsbond/terraform-provider-tailscale/internal/tailscale"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testTailnetKey = `
	resource "tailscale_tailnet_key" "example_key" {
		reusable = true
		ephemeral = true
	}
`

func TestProvider_TailscaleTailnetKey(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		PreCheck: func() {
			testServer.ResponseCode = http.StatusOK
			testServer.ResponseBody = tailscale.Key{
				ID:  "test",
				Key: "thisisatestkey",
			}
		},
		ProviderFactories: testProviderFactories(t),
		Steps: []resource.TestStep{
			testResourceCreated("tailscale_tailnet_key.example_key", testTailnetKey),
			testResourceDestroyed("tailscale_tailnet_key.example_key", testTailnetKey),
		},
	})
}
