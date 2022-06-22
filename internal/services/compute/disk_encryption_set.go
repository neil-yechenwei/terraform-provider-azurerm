package compute

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/sdk/2022-03-02/diskencryptionsets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/sdk/2022-03-02/disks"
)

// retrieveDiskEncryptionSetEncryptionType returns encryption type of the disk encryption set
func retrieveDiskEncryptionSetEncryptionType(ctx context.Context, client *diskencryptionsets.DiskEncryptionSetsClient, diskEncryptionSetId string) (*disks.EncryptionType, error) {
	diskEncryptionSet, err := diskencryptionsets.ParseDiskEncryptionSetID(diskEncryptionSetId)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, *diskEncryptionSet)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *diskEncryptionSet, err)
	}

	var encryptionType *disks.EncryptionType
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil && props.EncryptionType != nil {
			v := disks.EncryptionType(*props.EncryptionType)
			encryptionType = &v
		}

		if encryptionType == nil {
			return nil, fmt.Errorf("retrieving %s: EncryptionType was nil", *diskEncryptionSet)
		}
	}

	return encryptionType, nil
}
