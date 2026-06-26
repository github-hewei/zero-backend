package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"zero-backend/cmd/cli/runner"
	"zero-backend/internal/config"
	"zero-backend/internal/modules/rbac"
	"zero-backend/providers"

	"github.com/241x/zero-kit/gormutil"
	"github.com/241x/zero-kit/locker"
	"github.com/241x/zero-kit/logger"
	"github.com/241x/zero-kit/mongodb"
	"github.com/241x/zero-kit/mysql"
	"github.com/241x/zero-kit/queue"
	"github.com/241x/zero-kit/redis"
	"github.com/241x/zero-web/ctxkeys"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

func main() {
	cfg := config.New()

	mongoCfg := providers.NewMongoDBConfig(cfg)
	conn, err := mongodb.NewConn(mongoCfg)
	if err != nil {
		panic(err)
	}
	l := providers.ProvideLogger(cfg.Logger, conn.DB)

	rdb := redis.New(providers.NewRedisConfig(cfg))

	gormLog := gormutil.NewLogger(l)
	db, err := mysql.NewDB(providers.NewMySQLConfig(cfg), gormLog)
	if err != nil {
		panic(err)
	}

	root := &cobra.Command{
		Use:   "cli",
		Short: "CLI 工具",
		Long:  `命令行管理工具，用于执行数据迁移、用户管理等操作`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			traceID := uuid.New().String()
			cl := l.With("traceId", traceID)
			cl.Info("Start Command")

			ctx := cl.WithContext(cmd.Context())
			ctx = ctxkeys.WithTraceID(ctx, traceID)
			ctx = gormutil.WithTraceID(ctx, traceID)
			ctx = ctxkeys.WithBeginTime(ctx, time.Now())

			id, _ := cmd.Flags().GetInt("instance-id")
			key := fmt.Sprintf("%s:%d", strings.ReplaceAll(cmd.CommandPath(), " ", "-"), id)
			lock, err := locker.NewRedisLocker(rdb).Lock(ctx, key, locker.WithTTL(time.Minute), locker.WithWatchDog())
			if err != nil {
				return err
			}
			cl.Info("Locked Command", "key", key)

			cmd.SetContext(context.WithValue(ctx, lockKey, lock))
			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			if lock, ok := ctx.Value(lockKey).(locker.Lock); ok && lock != nil {
				unlockCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				lock.Unlock(unlockCtx)
				logger.Ctx(ctx).Info("Unlocked Command")
			}
			if t, ok := ctxkeys.BeginTime(ctx); ok {
				logger.Ctx(ctx).Info("End Command", "cost", fmt.Sprintf("%.4f", time.Since(t).Seconds()))
			}
		},
	}
	root.PersistentFlags().IntP("instance-id", "i", 0, "实例编号, 防止并发执行")

	root.AddCommand(migrateCmd(db, l))
	root.AddCommand(queueCmd(queue.NewQueueManager(rdb)))
	root.AddCommand(syncApiCmd(runner.NewSyncApiRunner(l, rbac.NewRbacApiRepository(db))))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	done := make(chan struct{})
	go func() {
		defer close(done)
		if err := root.ExecuteContext(ctx); err != nil {
			var cost time.Duration
			if t, ok := ctxkeys.BeginTime(root.Context()); ok {
				cost = time.Since(t)
			}
			logger.Ctx(root.Context()).Err(err, "Error Command", "cost", cost)
		}
	}()

	select {
	case <-sigChan:
		l.Info("Exiting manually...")
		cancel()
		timeout := time.NewTimer(10 * time.Second)
		defer timeout.Stop()
		select {
		case <-done:
		case <-timeout.C:
			l.Warn("Timeout exit")
		}
	case <-done:
	}
}

type ctxKey struct{}

var lockKey ctxKey
