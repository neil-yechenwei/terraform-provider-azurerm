package bot_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/botservice/mgmt/2021-03-01/botservice"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/bot/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type BotChannelTelegramResource struct {
}

func testAccBotChannelTelegram_basic(t *testing.T) {
	skipTelegramChannel(t)

	data := acceptance.BuildTestData(t, "azurerm_bot_channel_telegram", "test")
	r := BotChannelTelegramResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccBotChannelTelegram_requiresImport(t *testing.T) {
	skipTelegramChannel(t)

	data := acceptance.BuildTestData(t, "azurerm_bot_channel_telegram", "test")
	r := BotChannelTelegramResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func testAccBotChannelTelegram_update(t *testing.T) {
	skipTelegramChannel(t)

	data := acceptance.BuildTestData(t, "azurerm_bot_channel_telegram", "test")
	r := BotChannelTelegramResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t BotChannelTelegramResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.BotChannelID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Bot.ChannelClient.Get(ctx, id.ResourceGroup, id.BotServiceName, string(botservice.ChannelNameTelegramChannel))
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (BotChannelTelegramResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_telegram" "test" {
  bot_name                      = azurerm_bot_channels_registration.test.name
  location                      = azurerm_bot_channels_registration.test.location
  resource_group_name           = azurerm_resource_group.test.name
  telegram_channel_access_token = "%s"
}
`, BotChannelsRegistrationResource{}.basicConfig(data), os.Getenv("ARM_TEST_TELEGRAM_CHANNEL_ACCESS_TOKEN"))
}

func (r BotChannelTelegramResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_telegram" "test" {
  bot_name                      = azurerm_bot_channel_telegram.test.bot_name
  location                      = azurerm_bot_channel_telegram.test.location
  resource_group_name           = azurerm_bot_channel_telegram.test.resource_group_name
  telegram_channel_access_token = azurerm_bot_channel_telegram.test.telegram_channel_access_token
}
`, r.basic(data))
}

func (BotChannelTelegramResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_bot_channel_telegram" "test" {
  bot_name                      = azurerm_bot_channels_registration.test.name
  location                      = azurerm_bot_channels_registration.test.location
  resource_group_name           = azurerm_resource_group.test.name
  telegram_channel_access_token = "%s"
}
`, BotChannelsRegistrationResource{}.basicConfig(data), os.Getenv("ARM_TEST_TELEGRAM_CHANNEL_ACCESS_TOKEN2"))
}

func skipTelegramChannel(t *testing.T) {
	if os.Getenv("ARM_TEST_TELEGRAM_CHANNEL_ACCESS_TOKEN") == "" || os.Getenv("ARM_TEST_TELEGRAM_CHANNEL_ACCESS_TOKEN2") == "" {
		t.Skip("Skipping as one of `ARM_TEST_TELEGRAM_CHANNEL_ACCESS_TOKEN`, `ARM_TEST_TELEGRAM_CHANNEL_ACCESS_TOKEN2` was not specified")
	}
}
