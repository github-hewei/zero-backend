package providers

import (
	"zero-backend/modules/cli/command"
	"zero-backend/pkg/queue"

	"github.com/google/wire"
)

// CliCommandProviderSet 提供 CLI 命令依赖集合
var CliCommandProviderSet = wire.NewSet(
	command.NewRootCommand,
	command.NewUserCommand,
	command.NewUserListCommand,
	command.NewMigrateCommand,
	command.NewQueueCommand,
	queue.NewQueueManager,
)
