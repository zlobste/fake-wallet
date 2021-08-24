package cli

import (
	"github.com/alecthomas/kingpin"
	"github.com/zlobste/fake-wallet/internal/app"
	"github.com/zlobste/fake-wallet/internal/config"
	"os"
)

func Run(args []string) bool {
	cfg := config.New(os.Getenv("CONFIG"))
	log := cfg.Logging()

	defer func() {
		if rvr := recover(); rvr != nil {
			log.Error("internal panicked", rvr)
		}
	}()

	mainCmd := kingpin.New("fake-wallet", "")
	runCmd := mainCmd.Command("run", "run command")
	appCmd := runCmd.Command("app", "run appCmd")

	migrateCmd := mainCmd.Command("migrate", "migrate command")
	migrateUpCmd := migrateCmd.Command("up", "migrate db up")
	migrateDownCmd := migrateCmd.Command("down", "migrate db down")

	cmd, err := mainCmd.Parse(args[1:])
	if err != nil {
		log.WithError(err).Error("failed to parse arguments")
		return false
	}

	switch cmd {
	case appCmd.FullCommand():
		if err := app.New(cfg).Run(); err != nil {
			log.Error("failed to start appCmd", err)
			return false
		}
	case migrateUpCmd.FullCommand():
		err = MigrateUp(cfg)
	case migrateDownCmd.FullCommand():
		err = MigrateDown(cfg)
	default:
		log.Errorf("unknown command %s", cmd)
		return false
	}

	if err != nil {
		log.WithError(err).Error("failed to exec cmd")
		return false
	}

	return true
}
