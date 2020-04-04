package handler

import (
	"github.com/url_shortner/helpers"
	"github.com/url_shortner/i13n"
	"github.com/url_shortner/lib"
)

var MsgHandlers = map[string]interface{}{
	SortClicked:  HandleSortClicked,
	ProdPageCall: HandlePageCall,
}

func HandleEvent(message lib.UserProfileEvent) {

	anonymousID := message.AnonymousID
	eventName := message.EventName
	props := message.Properties
	propMap := helpers.ConvertMapInterfaceToString(props)
	i13n.Record(i13n.EventStatus,
		map[string]string{
			"name":    propMap["event_name"],
			"system":  lib.OriginDynamo,
			"bound":   "in",
			"success": "true",
		}, map[string]interface{}{
			"count": int(1),
		})
	if MsgHandlers[eventName] != nil {
		MsgHandlers[eventName].(func(string, map[string]string))(anonymousID, propMap)
	}
}

func EventShouldCompute(event_name string) (should_compute bool) {

	whitelisted_events := []string{
		"ProdPageCall",
		"PlpPageView",
		"PdpPageView",
		"LeadCreatedBackend",
		"FilterClicked",
		"PackageCompareClicked",
		"SortClicked",
		"PackageInteracted",
		"DestinationClicked",
		"LeadFunnel",
		"ExploreFormInteracted",
	}

	for _, allowed_event := range whitelisted_events {
		if allowed_event == event_name {
			should_compute = true
			return
		}
	}

	return
}
