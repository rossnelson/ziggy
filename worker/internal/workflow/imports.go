// Package workflows imports all workflow packages to trigger their init() registration
package workflow

import (
	_ "ziggy/internal/workflow/chat"
	_ "ziggy/internal/workflow/need_updater"
	_ "ziggy/internal/workflow/pool_regenerator"
	_ "ziggy/internal/workflow/ziggy"
)
