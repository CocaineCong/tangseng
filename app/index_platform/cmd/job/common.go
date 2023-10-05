package job

import (
	"context"

	"github.com/CocaineCong/tangseng/pkg/clone"
	logs "github.com/CocaineCong/tangseng/pkg/logger"
)

type Command struct {
	Name    string
	Exec    func(ctx context.Context) error
	Context context.Context

	done context.CancelFunc
}

func (cmd *Command) Run() {
	if cmd.Context == nil {
		cmd.Context = context.Background()
	}
	go func() {
		newCtx := clone.NewContextWithoutDeadline()
		err := cmd.Exec(newCtx)
		if err != nil {
			logs.LogrusObj.Error(err)
			return
		}
	}()
}
