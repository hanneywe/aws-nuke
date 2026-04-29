package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func fullResourceDef() ResourceDef {
	exported := true
	return ResourceDef{
		Service:      "mailmanager",
		Resource:     "relay",
		ResourceName: "MailManagerRelay",
		Scope:        "account",
		SDKPackage:   "mailmanager",
		SvcType:      "concrete",
		List: ListDef{
			Strategy:      "nested",
			Operation:     "ListRelays",
			Pagination:    "paginator",
			ItemsField:    "Relays",
			ItemType:      "Relay",
			ResponseField: "Relay",
			NilCheck:      true,
			Override:      "custom list body",
			Describe: &DescribeDef{
				Operation:     "DescribeRelay",
				InputMapping:  map[string]string{"RelayId": "RelayId"},
				ResponseField: "Relay",
			},
			Tags: &TagFetchDef{
				Operation: "ListTagsForResource",
				ArnField:  "RelayArn",
			},
			Levels: []NestedLevelDef{
				{
					Operation:   "ListChannelGroups",
					Pagination:  "paginator",
					ItemsField:  "Items",
					ItemType:    "ChannelGroupListConfiguration",
					IteratorVar: "cg",
				},
				{
					Operation:          "ListChannels",
					Pagination:         "paginator",
					ItemsField:         "Items",
					ItemType:           "ChannelListConfiguration",
					IteratorVar:        "ch",
					ParentInputMapping: map[string]string{"ChannelGroupName": "cg.ChannelGroupName"},
				},
			},
		},
		Delete: DeleteDef{
			Operation:   "DeleteRelay",
			InputFields: []string{"RelayId"},
			Override:    "custom delete code",
		},
		Fields: []FieldDef{
			{
				Name:         "RelayId",
				Type:         "*string",
				FromList:     "RelayId",
				FromDescribe: "Id",
				FromTags:     true,
				Exported:     &exported,
				PropertyTag:  "name=RelayId",
			},
		},
		Filters: []FilterDef{
			{
				Field:    "Status",
				Operator: "equals",
				Values:   []string{"DeleteInProgress", "Canceled"},
				Message:  "already being deleted",
			},
		},
		Settings: []SettingDef{
			{
				Name:             "DisableDeletionProtection",
				ProtectionField:  "protection",
				DisableOperation: "UpdateRelayConfig",
				DisableInput:     map[string]interface{}{"Name": "Name", "Protection": false},
			},
		},
		Dependencies: []string{"MailManagerRuleSet", "MailManagerEndpoint"},
		PreDeletion: []PreDeletionDef{
			{
				Type:             "listAndBatchDelete",
				ListOperation:    "ListTargets",
				ListInput:        map[string]string{"TargetGroupId": "ARN"},
				ListItemsField:   "Items",
				DeleteOperation:  "DeregisterTargets",
				DeleteInput:      map[string]string{"TargetGroupId": "ARN"},
				DeleteItemsField: "Targets",
				ItemMapping:      map[string]string{"Id": "Id", "Port": "Port"},
				Condition:        "r.Status == Active",
				Operation:        "DisablePortal",
				Input:            map[string]string{"PortalId": "PortalID"},
			},
		},
		StringRep: StringRepDef{
			Field:  "RelayName",
			Format: "%s (%s)",
			Fields: []string{"RelayId", "RelayName"},
			Conditional: &ConditionalStr{
				Field: "Name",
				IfNil: "RelayId",
				IfNotNil: &StringRepDef{
					Format: "%s (%s)",
					Fields: []string{"RelayId", "Name"},
				},
			},
			Override: "custom string body",
		},
		ExtraImports:   []string{"fmt", "strings"},
		FilterOverride: "custom filter body",
		IntegrationTest: &IntegTestDef{
			Create: &IntegCreateDef{
				Operation: "CreateRelay",
				Input:     map[string]interface{}{"RelayName": "test-relay"},
			},
		},
	}
}

func TestResourceDef_RoundTrip(t *testing.T) {
	a := assert.New(t)

	original := fullResourceDef()

	data, err := yaml.Marshal(&original)
	a.NoError(err)
	a.NotEmpty(data)

	var decoded ResourceDef
	err = yaml.Unmarshal(data, &decoded)
	a.NoError(err)

	a.Equal(original, decoded)
}

func TestResourceDef_OptionalFieldsOmitted(t *testing.T) {
	a := assert.New(t)

	minimal := ResourceDef{
		Service:      "eks",
		Resource:     "clusters",
		ResourceName: "EKSCluster",
		Scope:        "account",
		SDKPackage:   "eks",
		List: ListDef{
			Strategy: "flat",
		},
		Delete: DeleteDef{
			Operation:   "DeleteCluster",
			InputFields: []string{"Name"},
		},
		Fields: []FieldDef{
			{
				Name: "Name",
				Type: "*string",
			},
		},
		StringRep: StringRepDef{
			Field: "Name",
		},
	}

	data, err := yaml.Marshal(&minimal)
	a.NoError(err)

	expectedYAML := `service: eks
resource: clusters
resourceName: EKSCluster
scope: account
sdkPackage: eks
list:
    strategy: flat
delete:
    operation: DeleteCluster
    inputFields:
        - Name
fields:
    - name: Name
      type: '*string'
stringRepresentation:
    field: Name
`
	a.Equal(expectedYAML, string(data))

	var decoded ResourceDef
	err = yaml.Unmarshal(data, &decoded)
	a.NoError(err)
	a.Equal(minimal, decoded)
}
