package appgate

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/appgate/sdp-api-client-go/api/v22/openapi"
	"github.com/appgate/terraform-provider-appgatesdp/appgate/hashcode"
	"github.com/cenkalti/backoff/v4"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func BaseAuthContext(token string) context.Context {
	return context.WithValue(context.Background(), openapi.ContextAccessToken, token)
}

func AppendErrorf(diags diag.Diagnostics, format string, a ...any) diag.Diagnostics {
	return append(diags, diag.Errorf(format, a...)...)
}

func AppendFromErr(diags diag.Diagnostics, err error) diag.Diagnostics {
	if err == nil {
		return diags
	}
	return append(diags, diag.FromErr(err)...)
}

func mergeSchemaMaps(maps ...map[string]*schema.Schema) map[string]*schema.Schema {
	result := make(map[string]*schema.Schema)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

func baseEntitySchema() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Name of the object.",
			Required:    true,
		},
		"notes": {
			Type:        schema.TypeString,
			Description: "Notes for the object. Used for documentation purposes.",
			Default:     DefaultDescription,
			Optional:    true,
		},
	}
	return mergeSchemaMaps(s, baseTagsSchema())
}

func baseTagsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"tags": tagsSchema(),
	}
}

func tagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "Array of tags.",
		Optional:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		StateFunc: func(val interface{}) string {
			return strings.ToLower(val.(string))
		},
		Set: func(v interface{}) int {
			var buf bytes.Buffer
			str := v.(string)
			buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(str)))
			return hashcode.String(buf.String())
		},
	}
}

func ipPoolRange() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{

				"first": {
					Type:     schema.TypeString,
					Required: true,
				},

				"last": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
}

func readBaseEntityFromConfig(d *schema.ResourceData) (*openapi.BaseEntity, error) {
	base := &openapi.BaseEntity{}
	base.SetId(uuid.New().String())
	if v, ok := d.GetOk("name"); ok {
		base.SetName(v.(string))
	}
	if v, ok := d.GetOk("notes"); ok {
		base.SetNotes(v.(string))
	}
	if _, ok := d.GetOk("tags"); ok {
		base.SetTags(schemaExtractTags(d))
	}
	return base, nil
}

// prettyPrintAPIError is used to show a formatted error message
// from a HTTP 400-503 response from the http client.
func prettyPrintAPIError(err error) error {
	if err, ok := err.(*openapi.GenericOpenAPIError); ok {
		model := err.Model()
		if err, ok := model.(openapi.Error); ok {
			return fmt.Errorf("%s - %s", err.GetId(), err.GetMessage())
		}
		if err, ok := model.(openapi.ValidationError); ok {
			var ValidationErrors string
			errorMessage := "Validation error"
			for _, ve := range err.GetErrors() {
				ValidationErrors = ValidationErrors + ve.GetField() + " " + ve.GetMessage() + "\n"
			}
			if msg, o := err.GetMessageOk(); o {
				errorMessage = fmt.Sprintf("%s %s", errorMessage, *msg)
			}
			return fmt.Errorf("%s \n %s", errorMessage, ValidationErrors)
		}
		return fmt.Errorf("%s", err.Error())
	}
	return fmt.Errorf("%w", err)
}

func schemaExtractTags(d *schema.ResourceData) []string {
	rawtags := d.Get("tags").(*schema.Set).List()
	tags := make([]string, 0)
	for _, raw := range rawtags {
		tags = append(tags, strings.ToLower(raw.(string)))
	}
	return tags
}

func listToMapList(in []interface{}) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)
	for _, a := range in {
		source := a.(map[string]interface{})
		result = append(result, source)
	}
	return result, nil
}

func readAllowSources(in []interface{}) ([]openapi.AllowSourcesInner, error) {
	r := make([]openapi.AllowSourcesInner, 0)
	as, err := listToMapList(in)
	if err != nil {
		return r, err
	}

	for _, source := range as {
		row := openapi.NewAllowSourcesInnerWithDefaults()
		if v, ok := source["address"].(string); ok {
			row.SetAddress(v)
		}
		if v, ok := source["netmask"].(int); ok {
			row.SetNetmask(int32(v))
		}
		if v, ok := source["nic"].(string); ok && len(v) > 0 {
			row.SetNic(v)
		}
		r = append(r, *row)
	}

	return r, nil
}

func readAllowedUsers(in []interface{}) ([]openapi.PrometheusExporterAllowedUsersInner, error) {
	r := make([]openapi.PrometheusExporterAllowedUsersInner, 0)
	as, err := listToMapList(in)
	if err != nil {
		return r, err
	}
	for _, source := range as {
		row := openapi.NewPrometheusExporterAllowedUsersInner()
		if v, ok := source["username"].(string); ok {
			row.SetUsername(v)
		}
		if v, ok := source["password"].(string); ok {
			row.SetPassword(v)
		}
		r = append(r, *row)
	}
	return r, nil
}

func readLabelsDisabled(v []interface{}) ([]string, error) {
	l := []string{}
	for _, i := range v {
		s, ok := i.(string)
		if !ok {
			return nil, fmt.Errorf("invalid type: 'i', expected string")
		}
		l = append(l, s)
	}
	return l, nil
}
func readP12(in interface{}) (openapi.P12, error) {
	p12 := openapi.P12{}
	raw := in.(map[string]interface{})
	p12.SetId(uuid.New().String())
	if v, ok := raw["content"]; ok {
		content, err := appliancePortalReadp12Content(v.(string))
		if err != nil {
			return p12, fmt.Errorf("unable to read https_p12 file content %w", err)
		}
		p12.SetContent(content)
	}
	if v, ok := raw["password"]; ok {
		p12.SetPassword(v.(string))
	}
	return p12, nil
}

func readArrayOfFunctionsFromConfig(list []interface{}) ([]openapi.ApplianceFunction, error) {
	result := make([]openapi.ApplianceFunction, 0)
	for _, item := range list {
		if item == nil {
			continue
		}
		function := openapi.ApplianceFunction(item.(string))
		result = append(result, function)
	}
	return result, nil
}

func readArrayOfStringsFromConfig(list []interface{}) ([]string, error) {
	result := make([]string, 0)
	for _, item := range list {
		if item == nil {
			continue
		}
		result = append(result, item.(string))
	}
	return result, nil
}

func sliceToLowercase(l []openapi.ApplianceFunction) []string {
	result := make([]string, 0, len(l))
	for _, s := range l {
		result = append(result, strings.ToLower(string(s)))
	}
	return result
}

// validateIPaddress validate both IPv4 and IPv6 addresses.
func validateIPaddress(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if net.ParseIP(value) == nil {
		errors = append(errors, fmt.Errorf("Invalid ip address, got %s", value))
		return
	}
	return
}

func inArray(needle string, haystack []string) bool {
	sort.Strings(haystack)
	i := sort.Search(len(haystack),
		func(i int) bool { return haystack[i] >= needle })
	if i < len(haystack) && haystack[i] == needle {
		return true
	}
	return false
}

func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// getResourceFileContent gets content from key "file" filepath schema.ResourceData or string payload "content".
func getResourceFileContent(d *schema.ResourceData, key string) ([]byte, error) {
	var content []byte
	if v, ok := d.GetOk(key); ok {
		path := v.(string)
		file, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("Error opening file %q (%s): %w", key, path, err)
		}
		defer func() {
			err := file.Close()
			if err != nil {
				log.Printf("[WARN] Error closing file %q (%s): %s", key, path, err)
			}
		}()
		reader := bufio.NewReader(file)
		content, err = ioutil.ReadAll(reader)
		if err != nil {
			return nil, fmt.Errorf("Error reading file %q (%s): %w", key, path, err)
		}
	} else if v, ok := d.GetOk("content"); ok {
		content = []byte(v.(string))
	}
	return content, nil
}

// suppressMissingOptionalConfigurationBlock handles configuration block attributes in the following scenario:
//   - The resource schema includes an optional configuration block with defaults
//   - The API response includes those defaults to refresh into the Terraform state
//   - The operator's configuration omits the optional configuration block
func suppressMissingOptionalConfigurationBlock(k, old, new string, d *schema.ResourceData) bool {
	return old == "1" && new == "0"
}

// supressComputedResourceID handles a resource with a computed resource id.
// for example, most resource support that you set a UUID as its ID, if
// the ID is omitted, the API will computed one.
// if you set a custom id, for example
// appgatesdp_policy.example.policy_id = UUID
// the resource will get 2 values from this
// appgatesdp_policy.example.policy_id and appgatesdp_policy.example.id
// appgatesdp_policy.example.id is always Compute only
// and appgatesdp_policy.example.policy_id is ForceNew
// if the computed ID and the optional resource ID is not the same, we will show
// it as a diff.
func supressComputedResourceID(k, old, new string, d *schema.ResourceData) bool {
	// only detect diff if new is not null
	if d.Id() == old && len(new) == 0 {
		return true
	}
	return false
}

// resourceUUID provides a generic schema for all resource that support UUID as ID
func resourceUUID() *schema.Schema {
	return &schema.Schema{
		Type:             schema.TypeString,
		Description:      "ID of the object.",
		Optional:         true,
		ForceNew:         true,
		ValidateFunc:     validation.IsUUID,
		DiffSuppressFunc: supressComputedResourceID,
	}
}

func convertStringArrToInterface(strs []string) []interface{} {
	arr := make([]interface{}, len(strs))
	for i, str := range strs {
		arr[i] = str
	}
	return arr
}

// Nprintf is a Printf sibling (Nprintf; Named Printf), which handles strings like
// Nprintf("Hello %{target}!", map[string]interface{}{"target":"world"}) == "Hello world!".
// This is particularly useful for generated tests, where we don't want to use Printf,
// since that would require us to generate a very particular ordering of arguments.
func Nprintf(format string, params map[string]interface{}) string {
	for key, val := range params {
		format = strings.Replace(format, "%{"+key+"}", fmt.Sprintf("%v", val), -1)
	}
	return format
}

// ApplianceStatsRetryableError is used when /stats/appliance should be retried.
type ApplianceStatsRetryableError struct {
	err error
}

// Error returns non-empty string if there was an error.
func (e ApplianceStatsRetryableError) Error() string {
	return e.err.Error()
}

const (
	ApplianceStateInit                 = "init"
	ApplianceStateWaitingConfig        = "waiting_config"
	ApplianceStateMigratingdata        = "data_migration"
	ApplianceStateUpgrading            = "upgrading"
	ApplianceStateCloudinitializing    = "cloud_initializing"
	ApplianceStateApplianceActivating  = "appliance_activating"
	ApplianceStateApplianceRegistering = "appliance_registering"
	ApplianceStateApplianceReady       = "appliance_ready"
	ApplianceStateControllerReady      = "controller_ready"
)

// waitForApplianceState is a blocking function that does exponential backOff on appliance stats
// and make sure a certain appliance has reached state.
func waitForApplianceState(ctx context.Context, meta interface{}, applianceID, state string, b *backoff.ExponentialBackOff) error {
	return backoff.Retry(func() error {
		appliancesAPI := meta.(*Client).API.AppliancesApi
		token, err := meta.(*Client).GetToken()
		if err != nil {
			return ApplianceStatsRetryableError{err: err}
		}
		ctx = context.WithValue(ctx, openapi.ContextAccessToken, token)
		stats, _, err := appliancesAPI.AppliancesStatusGet(ctx).Execute()
		if err != nil {
			log.Printf("[ERROR] Failed to get appliance status for %s: %s", applianceID, err)
			return ApplianceStatsRetryableError{err: err}
		}
		var appliance openapi.ApplianceWithStatus
		for _, data := range stats.GetData() {
			if data.GetId() == applianceID {
				appliance = data
			}
		}
		if appliance.GetId() != applianceID {
			return fmt.Errorf("could not find appliance %q in stats list", applianceID)
		}
		got := appliance.GetState()
		log.Printf("[DEBUG] Appliance %s state is %s want state %s", applianceID, got, state)
		if got == state {
			return nil
		}
		return fmt.Errorf("appliance %q is in state %s expected %s", applianceID, appliance.GetState(), state)
	}, b)
}

func FileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}
