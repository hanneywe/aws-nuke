package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func validResourceDef() ResourceDef {
	return ResourceDef{
		Service:      "mailmanager",
		Resource:     "relay",
		ResourceName: "MailManagerRelay",
		Scope:        "account",
		SDKPackage:   "mailmanager",
		List: ListDef{
			Strategy:   "flat",
			Operation:  "ListRelays",
			Pagination: "paginator",
			ItemsField: "Relays",
			ItemType:   "Relay",
		},
		Delete: DeleteDef{
			Operation:   "DeleteRelay",
			InputFields: []string{"RelayId"},
		},
		Fields: []FieldDef{
			{Name: "RelayId", Type: "*string", FromList: "RelayId"},
		},
		StringRep: StringRepDef{
			Field: "RelayName",
		},
	}
}

func TestValidate_ValidDSL_NoErrors(t *testing.T) {
	a := assert.New(t)
	def := validResourceDef()
	errs := Validate(&def)
	a.Empty(errs)
}

func TestValidate_MissingService(t *testing.T) {
	a := assert.New(t)
	def := validResourceDef()
	def.Service = ""
	errs := Validate(&def)
	a.NotEmpty(errs)
	a.Equal("service", errs[0].Field)
}

func TestValidate_MissingListStrategy(t *testing.T) {
	a := assert.New(t)
	def := validResourceDef()
	def.List.Strategy = ""
	errs := Validate(&def)
	found := false
	for _, e := range errs {
		if e.Field == "list.strategy" {
			found = true
			break
		}
	}
	a.True(found, "expected error with field path list.strategy")
}

func TestValidate_InvalidStrategyValue(t *testing.T) {
	a := assert.New(t)
	def := validResourceDef()
	def.List.Strategy = "invalid"
	errs := Validate(&def)
	found := false
	for _, e := range errs {
		if e.Field == "list.strategy" {
			found = true
			a.Contains(e.Message, "must be flat, nested, or singleton")
			break
		}
	}
	a.True(found, "expected error with field path list.strategy")
}

func TestValidate_NestedListOnlyOneLevel(t *testing.T) {
	a := assert.New(t)
	def := validResourceDef()
	def.List.Strategy = StrategyNested
	def.List.Levels = []NestedLevelDef{
		{Operation: "ListChannelGroups", Pagination: "paginator", ItemsField: "Items", IteratorVar: "cg"},
	}
	errs := Validate(&def)
	found := false
	for _, e := range errs {
		if e.Field == "list.levels" {
			found = true
			a.Contains(e.Message, "at least 2 levels")
			break
		}
	}
	a.True(found, "expected error with field path list.levels")
}

func TestValidate_OverrideUnbalancedBraces(t *testing.T) {
	a := assert.New(t)
	def := validResourceDef()
	def.Delete.Override = "if true { doSomething()"
	errs := Validate(&def)
	found := false
	for _, e := range errs {
		if e.Field == "delete.override" {
			found = true
			a.Contains(e.Message, "unbalanced braces")
			break
		}
	}
	a.True(found, "expected error with field path delete.override")
}

func TestValidate_MissingItemTypeForFlatList(t *testing.T) {
	a := assert.New(t)
	def := validResourceDef()
	def.List.ItemType = ""
	errs := Validate(&def)
	found := false
	for _, e := range errs {
		if e.Field == "list.itemType" {
			found = true
			break
		}
	}
	a.True(found, "expected error for missing list.itemType")
}

func TestValidate_MissingItemTypeForNestedLevel(t *testing.T) {
	a := assert.New(t)
	def := validResourceDef()
	def.List.Strategy = StrategyNested
	def.List.ItemType = ""
	def.List.Levels = []NestedLevelDef{
		{Operation: "ListParents", Pagination: "paginator", ItemsField: "Parents", ItemType: "Parent", IteratorVar: "parent"},
		{Operation: "ListChildren", Pagination: "paginator", ItemsField: "Children", IteratorVar: "child"},
	}
	errs := Validate(&def)
	found := false
	for _, e := range errs {
		if e.Field == "list.levels[1].itemType" {
			found = true
			break
		}
	}
	a.True(found, "expected error for missing itemType on level 1")
}

func TestValidate_ItemTypeNotRequiredForConcreteType(t *testing.T) {
	a := assert.New(t)
	def := validResourceDef()
	def.SvcType = SvcTypeConcrete
	def.List.ItemType = ""
	errs := Validate(&def)
	for _, e := range errs {
		a.NotEqual("list.itemType", e.Field, "concrete types should not require itemType")
	}
}

func TestValidate_ItemTypeNotRequiredWithListOverride(t *testing.T) {
	a := assert.New(t)
	def := validResourceDef()
	def.List.ItemType = ""
	def.List.Override = "return nil, nil"
	errs := Validate(&def)
	for _, e := range errs {
		a.NotEqual("list.itemType", e.Field, "list override should skip itemType requirement")
	}
}
