package amlfilesystems

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AmlFilesystemsClient struct {
	Client *resourcemanager.Client
}

func NewAmlFilesystemsClientWithBaseURI(api environments.Api) (*AmlFilesystemsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "amlfilesystems", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AmlFilesystemsClient: %+v", err)
	}

	return &AmlFilesystemsClient{
		Client: client,
	}, nil
}