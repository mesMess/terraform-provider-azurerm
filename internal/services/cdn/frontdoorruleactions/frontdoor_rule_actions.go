package frontdoorruleactions

import (
	"fmt"
	"strings"

	track1 "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func expandStringSliceToCsvFormat(input []interface{}) *string {
	if len(input) == 0 {
		return nil
	}

	v := utils.ExpandStringSlice(input)
	csv := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(*v)), ","), "[]")

	return &csv
}

func flattenCsvToStringSlice(input *string) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	v := strings.Split(*input, ",")

	for _, s := range v {
		results = append(results, s)
	}

	return results
}

func ExpandFrontdoorRequestHeaderAction(input []interface{}) (*[]track1.BasicDeliveryRuleAction, error) {
	output := make([]track1.BasicDeliveryRuleAction, 0)

	for _, v := range input {
		item := v.(map[string]interface{})

		requestHeaderAction := track1.DeliveryRuleRequestHeaderAction{
			Name: track1.NameBasicDeliveryRuleActionNameModifyRequestHeader,
			Parameters: &track1.HeaderActionParameters{
				TypeName:     utils.String("DeliveryRuleHeaderActionParameters"),
				HeaderAction: track1.HeaderAction(item["header_action"].(string)),
				HeaderName:   utils.String(item["header_name"].(string)),
				Value:        utils.String(item["value"].(string)),
			},
		}

		output = append(output, requestHeaderAction)
	}

	return &output, nil
}

func ExpandFrontdoorResponseHeaderAction(input []interface{}) (*[]track1.BasicDeliveryRuleAction, error) {
	output := make([]track1.BasicDeliveryRuleAction, 0)

	for _, v := range input {
		item := v.(map[string]interface{})

		responseHeaderAction := track1.DeliveryRuleResponseHeaderAction{
			Name: track1.NameBasicDeliveryRuleActionNameModifyResponseHeader,
			Parameters: &track1.HeaderActionParameters{
				TypeName:     utils.String("DeliveryRuleHeaderActionParameters"),
				HeaderAction: track1.HeaderAction(item["header_action"].(string)),
				HeaderName:   utils.String(item["header_name"].(string)),
				Value:        utils.String(item["value"].(string)),
			},
		}

		output = append(output, responseHeaderAction)
	}

	return &output, nil
}

func ExpandFrontdoorUrlRedirectAction(input []interface{}) (*[]track1.BasicDeliveryRuleAction, error) {
	output := make([]track1.BasicDeliveryRuleAction, 0)

	for _, v := range input {
		item := v.(map[string]interface{})

		urlRedirectAction := track1.URLRedirectAction{
			Name: track1.NameBasicDeliveryRuleActionNameURLRedirect,
			Parameters: &track1.URLRedirectActionParameters{
				TypeName:            utils.String("DeliveryRuleUrlRedirectActionParameters"),
				RedirectType:        track1.RedirectType(item["redirect_type"].(string)),
				DestinationProtocol: track1.DestinationProtocol(item["redirect_protocol"].(string)),
				CustomPath:          utils.String(item["destination_path"].(string)),
				CustomHostname:      utils.String(item["destination_hostname"].(string)),
				CustomQueryString:   utils.String(item["query_string"].(string)),
				CustomFragment:      utils.String(item["destination_fragment"].(string)),
			},
		}

		output = append(output, urlRedirectAction)
	}

	return &output, nil
}

func ExpandFrontdoorUrlRewriteAction(input []interface{}) (*[]track1.BasicDeliveryRuleAction, error) {
	output := make([]track1.BasicDeliveryRuleAction, 0)

	for _, v := range input {
		item := v.(map[string]interface{})

		urlRedirectAction := track1.URLRewriteAction{
			Name: track1.NameBasicDeliveryRuleActionNameURLRedirect,
			Parameters: &track1.URLRewriteActionParameters{
				TypeName:              utils.String("DeliveryRuleUrlRewriteActionParameters"),
				Destination:           utils.String(item["destination"].(string)),
				PreserveUnmatchedPath: utils.Bool(item["preserve_unmatched_path"].(bool)),
				SourcePattern:         utils.String(item["source_pattern"].(string)),
			},
		}

		output = append(output, urlRedirectAction)
	}

	return &output, nil
}

func ExpandFrontdoorRouteConfigurationOverrideAction(input []interface{}) (*[]track1.BasicDeliveryRuleAction, error) {
	output := make([]track1.BasicDeliveryRuleAction, 0)

	for _, v := range input {
		item := v.(map[string]interface{})

		originGroupOverride := &track1.OriginGroupOverride{
			OriginGroup: &track1.ResourceReference{
				ID: utils.String(item["origin_group_id"].(string)),
			},
			ForwardingProtocol: track1.ForwardingProtocol(item["forwarding_protocol"].(string)),
		}

		compressionEnabled := track1.RuleIsCompressionEnabledEnabled
		if !item["compression_enabled"].(bool) {
			compressionEnabled = track1.RuleIsCompressionEnabledDisabled
		}

		// RuleQueryStringCachingBehavior
		cacheConfiguration := &track1.CacheConfiguration{
			QueryStringCachingBehavior: track1.RuleQueryStringCachingBehavior(item["query_string_caching_behavior"].(string)),
			QueryParameters:            expandStringSliceToCsvFormat(item["query_string_parameters"].([]interface{})),
			IsCompressionEnabled:       compressionEnabled,
			CacheBehavior:              track1.RuleCacheBehavior(item["cache_behavior"].(string)),
			CacheDuration:              utils.String(item["cache_duration"].(string)),
		}

		routeConfigurationOverrideAction := track1.DeliveryRuleRouteConfigurationOverrideAction{
			Name: track1.NameBasicDeliveryRuleActionNameRouteConfigurationOverride,
			Parameters: &track1.RouteConfigurationOverrideActionParameters{
				TypeName:            utils.String("DeliveryRuleRouteConfigurationOverrideAction"),
				OriginGroupOverride: originGroupOverride,
				CacheConfiguration:  cacheConfiguration,
			},
		}

		if queryStringCachingBehavior := *cacheConfiguration.QueryParameters; queryStringCachingBehavior == "" {
			if cacheConfiguration.QueryStringCachingBehavior == track1.RuleQueryStringCachingBehaviorIncludeSpecifiedQueryStrings || cacheConfiguration.QueryStringCachingBehavior == track1.RuleQueryStringCachingBehaviorIgnoreSpecifiedQueryStrings {
				return nil, fmt.Errorf("%q can not be empty if the %q is either %q or %q", "query_string_parameters", "query_string_caching_behavior", "IncludeSpecifiedQueryStrings", "IgnoreSpecifiedQueryStrings")
			}
		}

		output = append(output, routeConfigurationOverrideAction)
	}

	return &output, nil
}

func FlattenFrontdoorRequestHeaderAction(input track1.BasicDeliveryRuleAction) (*map[string]interface{}, error) {
	action, ok := input.AsDeliveryRuleRequestHeaderAction()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule request header action")
	}

	actionType := ""
	name := ""
	value := ""

	if params := action.Parameters; params != nil {
		actionType = string(params.HeaderAction)
		name = string(*params.HeaderName)
		value = string(*params.Value)
	}

	return &map[string]interface{}{
		"header_action": actionType,
		"header_name":   name,
		"value":         value,
	}, nil
}

func FlattenFrontdoorResponseHeaderAction(input track1.BasicDeliveryRuleAction) (*map[string]interface{}, error) {
	action, ok := input.AsDeliveryRuleResponseHeaderAction()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule reesponse header action")
	}

	actionType := ""
	name := ""
	value := ""

	if params := action.Parameters; params != nil {
		actionType = string(params.HeaderAction)
		name = string(*params.HeaderName)
		value = string(*params.Value)
	}

	return &map[string]interface{}{
		"header_action": actionType,
		"header_name":   name,
		"value":         value,
	}, nil
}

func FlattenFrontdoorUrlRedirectAction(input track1.BasicDeliveryRuleAction) (*map[string]interface{}, error) {
	action, ok := input.AsURLRedirectAction()
	if !ok {
		return nil, fmt.Errorf("expected a URL redirect action")
	}

	destinationHost := ""
	destinationPath := ""
	queryString := ""
	destinationProtocol := ""
	redirectType := ""
	fragment := ""

	if params := action.Parameters; params != nil {
		destinationHost = string(*params.CustomHostname)
		destinationPath = string(*params.CustomPath)
		queryString = string(*params.CustomQueryString)
		destinationProtocol = string(params.DestinationProtocol)
		redirectType = string(params.RedirectType)
		fragment = string(*params.CustomFragment)
	}

	return &map[string]interface{}{
		"destination_hostname": destinationHost,
		"destination_path":     destinationPath,
		"query_string":         queryString,
		"redirect_protocol":    destinationProtocol,
		"redirect_type":        redirectType,
		"destination_fragment": fragment,
	}, nil
}

func FlattenFrontdoorUrlRewriteAction(input track1.BasicDeliveryRuleAction) (*map[string]interface{}, error) {
	action, ok := input.AsURLRewriteAction()
	if !ok {
		return nil, fmt.Errorf("expected a URL redirect action")
	}

	destination := ""
	preservePath := false
	sourcePattern := ""

	if params := action.Parameters; params != nil {
		destination = string(*params.Destination)
		preservePath = bool(*params.PreserveUnmatchedPath)
		sourcePattern = string(*params.SourcePattern)
	}

	return &map[string]interface{}{
		"destination":             destination,
		"preserve_unmatched_path": preservePath,
		"source_pattern":          sourcePattern,
	}, nil
}

func FlattenFrontdoorRouteConfigurationOverrideAction(input track1.BasicDeliveryRuleAction) (*map[string]interface{}, error) {
	action, ok := input.AsDeliveryRuleRouteConfigurationOverrideAction()
	if !ok {
		return nil, fmt.Errorf("expected a route configuration override action")
	}

	queryStringCachingBehavior := ""
	cacheBehavior := ""
	compressionEnabled := false
	cacheDuration := ""
	queryParameters := make([]interface{}, 0)
	forwardingProtocol := ""
	originGroupId := ""

	if params := action.Parameters; params != nil {
		queryStringCachingBehavior = string(params.CacheConfiguration.QueryStringCachingBehavior)
		cacheBehavior = string(params.CacheConfiguration.CacheBehavior)
		compressionEnabled = (params.CacheConfiguration.IsCompressionEnabled == track1.RuleIsCompressionEnabledEnabled)
		cacheDuration = string(*params.CacheConfiguration.CacheDuration)
		queryParameters = flattenCsvToStringSlice(params.CacheConfiguration.QueryParameters)
		forwardingProtocol = string(params.OriginGroupOverride.ForwardingProtocol)
		originGroupId = string(*params.OriginGroupOverride.OriginGroup.ID)
	}

	return &map[string]interface{}{
		"query_string_caching_behavior": queryStringCachingBehavior,
		"cache_behavior":                cacheBehavior,
		"compression_enabled":           compressionEnabled,
		"cache_duration":                cacheDuration,
		"query_string_parameters":       queryParameters,
		"forwarding_protocol":           forwardingProtocol,
		"origin_group_id":               originGroupId,
	}, nil
}
