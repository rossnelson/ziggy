// Package workflows imports all workflow packages to trigger their init() registration
package workflow

import (
	"ziggy/internal/workflow/chat"
	"ziggy/internal/workflow/need_updater"
	"ziggy/internal/workflow/pool_regenerator"
	"ziggy/internal/workflow/ziggy"
)

func RegisterWorkflows() {
	chat.Register()
	need_updater.Register()
	pool_regenerator.Register()
	ziggy.Register()
}
