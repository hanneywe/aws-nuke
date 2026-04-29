package listsettings

import (
	"context"
	"slices"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"

	"github.com/ekristen/aws-nuke/v3/pkg/commands/global"
	"github.com/ekristen/aws-nuke/v3/pkg/common"

	"github.com/ekristen/libnuke/pkg/registry"

	_ "github.com/ekristen/aws-nuke/v3/resources"
)

func execute(_ context.Context, c *cli.Command) error {
	regs := registry.GetRegistrations()

	type entry struct {
		resource string
		setting  string
	}

	var entries []entry
	for name, reg := range regs {
		for _, s := range reg.Settings {
			entries = append(entries, entry{resource: name, setting: s})
		}
	}

	slices.SortFunc(entries, func(a, b entry) int {
		if a.resource < b.resource {
			return -1
		}
		if a.resource > b.resource {
			return 1
		}
		if a.setting < b.setting {
			return -1
		}
		if a.setting > b.setting {
			return 1
		}
		return 0
	})

	for _, e := range entries {
		color.New(color.Bold).Printf("%-55s", e.resource)
		color.New(color.FgCyan).Printf("%-45s", e.setting)
		// NOTE: all registered settings are currently bool. The registry does not store type
		// metadata, so if int or string settings are registered in the future, this will need
		// to be updated to resolve the type (e.g. by inspecting GetBool/GetInt/GetString usage).
		color.New(color.FgHiBlack).Printf("bool\n")
	}

	return nil
}

func init() {
	cmd := &cli.Command{
		Name:    "resource-settings",
		Aliases: []string{"list-settings"},
		Usage:   "list available settings for resource types",
		Flags:   global.Flags(),
		Before:  global.Before,
		Action:  execute,
	}

	common.RegisterCommand(cmd)
}
