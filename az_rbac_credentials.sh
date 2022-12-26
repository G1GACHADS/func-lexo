#!/bin/bash
set -ex

az ad sp create-for-rbac --name $NAME --role contributor \
    --scopes /subscriptions/$SUBSCRIPTION_ID/resourceGroups/$RESOURCE_GROUP_NAME/providers/Microsoft.Web/sites/$AZURE_FUNCTIONAPP_NAME \
    --sdk-auth
