package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources/v3"
	"github.com/sirupsen/logrus"
)

type AzureAdapter struct {
	ctx         context.Context
	groupClient *armresources.ResourceGroupsClient
}

func New(subscriptionId string) (*AzureAdapter, error) {
	credential, _ := azidentity.NewDefaultAzureCredential(nil)
	groupsClient, err := armresources.NewResourceGroupsClient(subscriptionId, credential, nil)

	return &AzureAdapter{
		ctx:         context.Background(),
		groupClient: groupsClient,
	}, err
}

func (a *AzureAdapter) RemoveResourceGroup(resourceGroupName string) error {
	result, err := a.groupClient.CheckExistence(a.ctx, resourceGroupName, &armresources.ResourceGroupsClientCheckExistenceOptions{})
	if err != nil {
		return err
	}

	if result.Success {
		poller, err := a.groupClient.BeginDelete(a.ctx, resourceGroupName, nil)
		if err != nil {
			return err
		}
		_, err = poller.PollUntilDone(a.ctx, nil)
		if err != nil {
			return err
		}

		logrus.Infof("Resourcegroup %s removed", resourceGroupName)
	} else {
		logrus.Infof("Resourcegroup %s does not exist. Remove skipped.", resourceGroupName)
	}

	return nil
}
