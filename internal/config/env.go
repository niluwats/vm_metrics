package config

import (
	"log"
	"os"
	"strconv"
)

// ParseEnvironment loads a sibling `.env` file then looks through all environment
// variables to set global configuration.
func ParseEnvironment() error {
	os.Setenv("AZURE_GROUP_NAME", "vmtest")
	os.Setenv("AZURE_BASE_GROUP_NAME", "vmtest")
	os.Setenv("AZURE_LOCATION_DEFAULT", "eastus")
	// os.Setenv("AZURE_USE_DEVICEFLOW", "true")
	// os.Setenv("AZURE_SAMPLES_KEEP_RESOURCES", "0")
	os.Setenv("AZURE_CLIENT_ID", "98dae7e1-1f6f-45b2-b61a-176308cded01")
	os.Setenv("AZURE_CLIENT_SECRET", "Cqi64W-yGzxiVu_PWzxdTimYAy9YD4tQ8i")
	os.Setenv("AZURE_TENANT_ID", "5d88d8f9-6ecf-43be-9867-21ce46f07b74")
	os.Setenv("AZURE_SUBSCRIPTION_ID", "dcb8d550-31da-4e1b-afd6-b41939ab339c")
	// AZURE_GROUP_NAME and `config.GroupName()` are deprecated.
	// Use AZURE_BASE_GROUP_NAME and `config.GenerateGroupName()` instead.
	Group_Name = os.Getenv("AZURE_GROUP_NAME")
	BaseGroup_Name = os.Getenv("AZURE_BASE_GROUP_NAME")

	Location_Default = os.Getenv("AZURE_LOCATION_DEFAULT")

	var err error
	UseDevice_Flow, err = strconv.ParseBool(os.Getenv("AZURE_USE_DEVICEFLOW"))
	if err != nil {
		log.Printf("invalid value specified for AZURE_USE_DEVICEFLOW, disabling\n")
		UseDevice_Flow = false
	}
	Keep_Resources, err = strconv.ParseBool(os.Getenv("AZURE_SAMPLES_KEEP_RESOURCES"))
	if err != nil {
		log.Printf("invalid value specified for AZURE_SAMPLES_KEEP_RESOURCES, discarding\n")
		Keep_Resources = false
	}

	// these must be provided by environment
	// clientID
	// clientID = os.Getenv("AZURE_CLIENT_ID")
	Client_ID = os.Getenv("AZURE_CLIENT_ID")

	// clientSecret
	Client_Secret = os.Getenv("AZURE_CLIENT_SECRET")

	// tenantID (AAD)
	Tenant_ID = os.Getenv("AZURE_TENANT_ID")

	// subscriptionID (ARM)
	Subscription_ID = os.Getenv("AZURE_SUBSCRIPTION_ID")
	// Cloud_Name=""
	return nil
}
