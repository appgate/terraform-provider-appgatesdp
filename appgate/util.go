package appgate

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sort"

	"github.com/appgate/terraform-provider-appgate/client/v12/openapi"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	return map[string]*schema.Schema{
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
		"tags": {
			Type:        schema.TypeSet,
			Description: "Array of tags.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
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
	return fmt.Errorf("%+v", err)
}

func schemaExtractTags(d *schema.ResourceData) []string {
	rawtags := d.Get("tags").(*schema.Set).List()
	tags := make([]string, 0)
	for _, raw := range rawtags {
		tags = append(tags, raw.(string))
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

// validateCIDRNetworkAddress ensures that the string value is a valid CIDR that
// represents a network address - it adds an error otherwise
func validateCIDRNetworkAddress(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid CIDR, got error parsing: %s", k, err))
		return
	}

	if ipnet == nil || value != ipnet.String() {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid network CIDR, got %q", k, value))
	}

	return
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

// getResourceFileContent gets content from "file" filepath schema.ResourceData or string payload "content".
func getResourceFileContent(d *schema.ResourceData) ([]byte, error) {
	var content []byte
	if v, ok := d.GetOk("file"); ok {
		path := v.(string)
		file, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("Error opening file (%s): %s", path, err)
		}
		defer func() {
			err := file.Close()
			if err != nil {
				log.Printf("[WARN] Error closing file (%s): %s", path, err)
			}
		}()
		reader := bufio.NewReader(file)
		content, err = ioutil.ReadAll(reader)
		if err != nil {
			return nil, fmt.Errorf("Error reading file (%s): %s", path, err)
		}
	} else if v, ok := d.GetOk("content"); ok {
		content = []byte(v.(string))
	}
	return content, nil
}
