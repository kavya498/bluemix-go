package endpoints

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/bluemix-go/helpers"
)

//EndpointLocator ...
type EndpointLocator interface {
	AccountManagementEndpoint() (string, error)
	CertificateManagerEndpoint() (string, error)
	CFAPIEndpoint() (string, error)
	ContainerEndpoint() (string, error)
	ContainerRegistryEndpoint() (string, error)
	CisEndpoint() (string, error)
	GlobalSearchEndpoint() (string, error)
	GlobalTaggingEndpoint() (string, error)
	IAMEndpoint() (string, error)
	IAMPAPEndpoint() (string, error)
	ICDEndpoint() (string, error)
	MCCPAPIEndpoint() (string, error)
	ResourceManagementEndpoint() (string, error)
	ResourceControllerEndpoint() (string, error)
	ResourceCatalogEndpoint() (string, error)
	UAAEndpoint() (string, error)
	CseEndpoint() (string, error)
	SchematicsEndpoint() (string, error)
	UserManagementEndpoint() (string, error)
	HpcsEndpoint() (string, error)
	FunctionsEndpoint() (string, error)
}

const (
	//ErrCodeServiceEndpoint ...
	ErrCodeServiceEndpoint = "ServiceEndpointDoesnotExist"
)

var regionToEndpoint = map[string]map[string]string{
	"cf": {
		"us-south": "https://api.ng.bluemix.net",
		"us-east":  "https://api.us-east.bluemix.net",
		"eu-gb":    "https://api.eu-gb.bluemix.net",
		"au-syd":   "https://api.au-syd.bluemix.net",
		"eu-de":    "https://api.eu-de.bluemix.net",
		"jp-tok":   "https://api.jp-tok.bluemix.net",
	},
	"cr": {
		"us-south": "us.icr.io",
		"us-east":  "us.icr.io",
		"eu-de":    "de.icr.io",
		"au-syd":   "au.icr.io",
		"eu-gb":    "uk.icr.io",
		"jp-tok":   "jp.icr.io",
		"jp-osa":   "jp2.icr.io",
	},
	"uaa": {
		"us-south": "https://iam.cloud.ibm.com/cloudfoundry/login/us-south",
		"us-east":  "https://iam.cloud.ibm.com/cloudfoundry/login/us-east",
		"eu-gb":    "https://iam.cloud.ibm.com/cloudfoundry/login/uk-south",
		"au-syd":   "https://iam.cloud.ibm.com/cloudfoundry/login/ap-south",
		"eu-de":    "https://iam.cloud.ibm.com/cloudfoundry/login/eu-central",
	},
}
var privateRegions = []string{"us-south", "us-east"}

func init() {
	//TODO populate the endpoints which can be retrieved from given endpoints dynamically
	//Example - UAA can be found from the CF endpoint
}

type endpointLocator struct {
	region     string
	visibility string
}

//NewEndpointLocator ...
func NewEndpointLocator(region, visibility string) EndpointLocator {
	return &endpointLocator{region: region, visibility: visibility}
}

func (e *endpointLocator) AccountManagementEndpoint() (string, error) {
	//As the current list of regionToEndpoint above is not exhaustive we allow to read endpoints from the env
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_ACCOUNT_MANAGEMENT_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.visibility == "private" && (e.region == "us-south" || e.region == "us-east") {
		return fmt.Sprintf("https://private.%s.accounts.cloud.ibm.com", e.region), nil
	}
	return "https://accounts.cloud.ibm.com", nil
}

func (e *endpointLocator) CertificateManagerEndpoint() (string, error) {
	//As the current list of regionToEndpoint above is not exhaustive we allow to read endpoints from the env
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_CERTIFICATE_MANAGER_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.visibility == "private" {
		return fmt.Sprintf("https://private.%s.certificate-manager.cloud.ibm.com", e.region), nil
	}
	return fmt.Sprintf("https://%s.certificate-manager.cloud.ibm.com", e.region), nil
}

func (e *endpointLocator) CFAPIEndpoint() (string, error) {
	//As the current list of regionToEndpoint above is not exhaustive we allow to read endpoints from the env
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_CF_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	} else if ep, ok := regionToEndpoint["cf"][e.region]; ok {
		return ep, nil
	}
	if e.visibility == "private" {
		return "", bmxerror.New(ErrCodeServiceEndpoint, fmt.Sprintf("Private Endpoints is not supported by this service"))
	}
	return "", bmxerror.New(ErrCodeServiceEndpoint, fmt.Sprintf("Cloud Foundry endpoint doesn't exist for region: %q", e.region))
}

func (e *endpointLocator) ContainerEndpoint() (string, error) {
	//As the current list of regionToEndpoint above is not exhaustive we allow to read endpoints from the env
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_CS_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.visibility == "private" {
		return "", bmxerror.New(ErrCodeServiceEndpoint, fmt.Sprintf("Private Endpoints is not supported by this service"))
	}
	return "https://containers.cloud.ibm.com/global", nil
}

func (e *endpointLocator) SchematicsEndpoint() (string, error) {
	//As the current list of regionToEndpoint above is not exhaustive we allow to read endpoints from the env
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_SCHEMATICS_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.visibility == "private" {
		if e.region == "us-south" || e.region == "us-east" {
			ep := "https://private-us.schematics.cloud.ibm.com"
			return ep, nil
		}
		if e.region == "eu-gb" || e.region == "eu-de" {
			ep := "https://private-eu.schematics.cloud.ibm.com"
			return ep, nil
		}
		return "", bmxerror.New(ErrCodeServiceEndpoint, fmt.Sprintf("Private Endpoints is not supported by this service for the region %s", e.region))

	}
	return fmt.Sprintf("https://%s.schematics.cloud.ibm.com", e.region), nil
}

func (e *endpointLocator) ContainerRegistryEndpoint() (string, error) {
	//As the current list of regionToEndpoint above is not exhaustive we allow to read endpoints from the env
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_CR_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	} else if ep, ok := regionToEndpoint["cr"][e.region]; ok {
		if e.visibility == "private" {
			return fmt.Sprintf("https://private.%s", ep), nil
		}
		return fmt.Sprintf("https://%s", ep), nil
	}
	return "", bmxerror.New(ErrCodeServiceEndpoint, fmt.Sprintf("Container Registry endpoint doesn't exist for region: %q", e.region))
}

// Not used in Provider as we have migrated to go-sdk
func (e *endpointLocator) CisEndpoint() (string, error) {
	//As the current list of regionToEndpoint above is not exhaustive we allow to read endpoints from the env
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_CIS_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.visibility == "private" {
		return "https://api.private.cis.cloud.ibm.com", nil
	}
	return "https://api.cis.cloud.ibm.com", nil
}

func (e *endpointLocator) GlobalSearchEndpoint() (string, error) {
	//As the current list of regionToEndpoint above is not exhaustive we allow to read endpoints from the env
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_GS_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.visibility == "private" && (e.region == "us-south" || e.region == "us-east") {
		return fmt.Sprintf("https://api.private.%s.global-search-tagging.cloud.ibm.com", e.region), nil
	}
	return "https://api.global-search-tagging.cloud.ibm.com", nil
}

func (e *endpointLocator) GlobalTaggingEndpoint() (string, error) {
	//As the current list of regionToEndpoint above is not exhaustive we allow to read endpoints from the env
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_GT_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.visibility == "private" && (e.region == "us-south" || e.region == "us-east") {
		return fmt.Sprintf("https://tags.private.%s.global-search-tagging.cloud.ibm.com", e.region), nil
	}
	return "https://tags.global-search-tagging.cloud.ibm.com", nil
}

func (e *endpointLocator) IAMEndpoint() (string, error) {
	//As the current list of regionToEndpoint above is not exhaustive we allow to read endpoints from the env
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_IAM_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.visibility == "private" {
		if e.region == "us-south" || e.region == "us-east" {
			return fmt.Sprintf("https://private.%s.iam.cloud.ibm.com", e.region), nil
		}
		return "https://private.iam.cloud.ibm.com", nil
	}
	return "https://iam.cloud.ibm.com", nil
}

func (e *endpointLocator) IAMPAPEndpoint() (string, error) {
	//As the current list of regionToEndpoint above is not exhaustive we allow to read endpoints from the env
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_IAMPAP_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.visibility == "private" {
		if e.region == "us-south" || e.region == "us-east" {
			return fmt.Sprintf("https://private.%s.iam.cloud.ibm.com", e.region), nil
		}
		return "https://private.iam.cloud.ibm.com", nil
	}
	return "https://iam.cloud.ibm.com", nil
}

func (e *endpointLocator) ICDEndpoint() (string, error) {
	//As the current list of regionToEndpoint above is not exhaustive we allow to read endpoints from the env
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_ICD_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.visibility == "private" {
		return fmt.Sprintf("https://api.%s.private.databases.cloud.ibm.com", e.region), nil
	}
	return fmt.Sprintf("https://api.%s.databases.cloud.ibm.com", e.region), nil
}

func (e *endpointLocator) MCCPAPIEndpoint() (string, error) {
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_MCCP_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.visibility == "private" {
		return "", bmxerror.New(ErrCodeServiceEndpoint, fmt.Sprintf("Private Endpoints is not supported by this service for the region %s", e.region))
	}
	return fmt.Sprintf("https://mccp.%s.cf.cloud.ibm.com", e.region), nil
}

func (e *endpointLocator) ResourceManagementEndpoint() (string, error) {
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_RESOURCE_MANAGEMENT_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.visibility == "private" && (e.region == "us-south" || e.region == "us-east") {
		return fmt.Sprintf("https://private.%s.resource-controller.cloud.ibm.com", e.region), nil
	}
	return "https://resource-controller.cloud.ibm.com", nil
}

func (e *endpointLocator) ResourceControllerEndpoint() (string, error) {
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_RESOURCE_CONTROLLER_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.visibility == "private" && (e.region == "us-south" || e.region == "us-east") {
		return fmt.Sprintf("https://private.%s.resource-controller.cloud.ibm.com", e.region), nil
	}
	return "https://resource-controller.cloud.ibm.com", nil
}

func (e *endpointLocator) ResourceCatalogEndpoint() (string, error) {
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_RESOURCE_CATALOG_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.visibility == "private" && (e.region == "us-south" || e.region == "us-east") {
		return fmt.Sprintf("https://private.%s.globalcatalog.cloud.ibm.com", e.region), nil
	}
	return "https://globalcatalog.cloud.ibm.com", nil
}

func (e *endpointLocator) UAAEndpoint() (string, error) {
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_UAA_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.visibility == "private" {
		return "", bmxerror.New(ErrCodeServiceEndpoint, fmt.Sprintf("Private Endpoints is not supported by this service for the region %s", e.region))
	}
	if ep, ok := regionToEndpoint["uaa"][e.region]; ok {
		//As the current list of regionToEndpoint above is not exhaustive we allow to read endpoints from the env
		return ep, nil
	}
	return "", bmxerror.New(ErrCodeServiceEndpoint, fmt.Sprintf("UAA endpoint doesn't exist for region: %q", e.region))
}

func (e *endpointLocator) CseEndpoint() (string, error) {
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_CSE_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.visibility == "private" {
		return "", bmxerror.New(ErrCodeServiceEndpoint, fmt.Sprintf("Private Endpoints is not supported by this service"))
	}
	return "https://api.serviceendpoint.cloud.ibm.com", nil
}

func (e *endpointLocator) UserManagementEndpoint() (string, error) {
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_USER_MANAGEMENT_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.visibility == "private" && (e.region == "us-south" || e.region == "us-east") {
		return fmt.Sprintf("https://private.%s.user-management.cloud.ibm.com", e.region), nil
	}
	return "https://user-management.cloud.ibm.com", nil
}

func (e *endpointLocator) HpcsEndpoint() (string, error) {
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_HPCS_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	return fmt.Sprintf("https://%s.broker.hs-crypto.cloud.ibm.com/crypto_v2/", e.region), nil
}

func (e *endpointLocator) FunctionsEndpoint() (string, error) {
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_FUNCTIONS_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.visibility == "private" {
		return "", bmxerror.New(ErrCodeServiceEndpoint, fmt.Sprintf("Private Endpoints is not supported by this service for the region %s", e.region))
	}
	return fmt.Sprintf("https://%s.functions.cloud.ibm.com", e.region), nil
}
