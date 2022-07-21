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
	"time"

	"github.com/appgate/sdp-api-client-go/api/v17/openapi"
	"github.com/appgate/terraform-provider-appgatesdp/appgate/hashcode"
	"github.com/cenkalti/backoff/v4"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

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

func readBaseEntityFromConfig(d *schema.ResourceData) (*openapi.BaseEntity, error) {
	base := &openapi.BaseEntity{}
	base.Id = uuid.New().String()
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
	if err, ok := err.(openapi.GenericOpenAPIError); ok {
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

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
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

func applianceStatsRetryable(ctx context.Context, meta interface{}) *resource.RetryError {
	if err := checkApplianceStatus(ctx, meta)(); err != nil {
		if err, ok := err.(ApplianceStatsRetryableError); ok {
			return resource.RetryableError(err)
		}
		return resource.NonRetryableError(err)
	}
	return nil
}

// ApplianceStatsRetryableError is used when /stats/appliance should be retried.
type ApplianceStatsRetryableError struct {
	err error
}

// Error returns non-empty string if there was an error.
func (e ApplianceStatsRetryableError) Error() string {
	return e.err.Error()
}

func checkApplianceStatus(ctx context.Context, meta interface{}) func() error {
	return func() error {
		statsAPI := meta.(*Client).API.ApplianceStatsApi
		token, err := meta.(*Client).GetToken()
		if err != nil {
			return err
		}
		stats, _, err := statsAPI.StatsAppliancesGet(ctx).Authorization(token).Execute()
		if err != nil {
			return ApplianceStatsRetryableError{err: err}
		}
		numberOfControllers := int(stats.GetControllerCount())
		controllers := make([]openapi.StatsAppliancesListAllOfData, 0, numberOfControllers)
		for _, data := range stats.GetData() {
			c := data.GetController()
			// all none controller appliances will return n/a as status
			if c.GetStatus() != "n/a" {
				controllers = append(controllers, data)
			}
		}
		if len(controllers) != numberOfControllers {
			log.Printf("[DEBUG] Found %d controller expected %d", len(controllers), numberOfControllers)
		}
		for _, controller := range controllers {
			log.Printf("[DEBUG] Wait for controllers %s %s %s", controller.GetName(), controller.GetState(), controller.GetStatus())
			if controller.GetStatus() == "busy" {
				return ApplianceStatsRetryableError{err: fmt.Errorf("%s is busy, got %s", controller.GetName(), controller.GetStatus())}
			}
		}
		return nil
	}
}

// waitForControllers is a blocking function that does exponential backOff on appliance stats
// and make sure all the controllers are healthy before returning nil
func waitForControllers(ctx context.Context, meta interface{}) error {
	return backoff.Retry(checkApplianceStatus(ctx, meta), &backoff.ExponentialBackOff{
		InitialInterval:     2 * time.Second,
		RandomizationFactor: 0.7,
		Multiplier:          2,
		MaxInterval:         5 * time.Minute,
		Stop:                backoff.Stop,
		Clock:               backoff.SystemClock,
	})
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
		statsAPI := meta.(*Client).API.ApplianceStatsApi
		token, err := meta.(*Client).GetToken()
		if err != nil {
			return ApplianceStatsRetryableError{err: err}
		}
		stats, _, err := statsAPI.StatsAppliancesGet(ctx).Authorization(token).Execute()
		if err != nil {
			return ApplianceStatsRetryableError{err: err}
		}
		var appliance openapi.StatsAppliancesListAllOfData
		for _, data := range stats.GetData() {
			if data.GetId() == applianceID {
				appliance = data
			}
		}
		if appliance.GetId() != applianceID {
			return fmt.Errorf("could not find appliance %q in stats list", applianceID)
		}
		if appliance.GetState() == state {
			log.Printf("[DEBUG] Appliance %q reached expected state %s", applianceID, state)
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
