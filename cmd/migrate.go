package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

var migrateCount *int

func dsn() string {
	dbhost := "127.0.0.1"
	dbport := os.Getenv("PG_PORT")
	dbusername := os.Getenv("PG_USER")
	dbpassword := os.Getenv("PG_PASS")
	dbname := os.Getenv("PG_DB")
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbusername, dbpassword, dbhost, dbport, dbname)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate the db schema",
	Long:  `Migrate the db schema required by country-incomes`,
	Run: func(cmd *cobra.Command, args []string) {
		m, err := migrate.New(
			"file://db/migrations",
			dsn())
		if err != nil {
			log.WithFields(log.Fields{"err": err.Error()}).Error("failed to create migration")
			return
		}

		if *migrateCount == 0 {

			err := m.Up()
			if err == migrate.ErrNoChange {
				log.Info("no change")
				return
			}

			if err != nil {
				log.WithFields(log.Fields{"err": err.Error()}).Error("migration failed")
				return
			}
			log.Info("all migrations were run successfully")

		} else {
			err := m.Steps(*migrateCount)

			if err == migrate.ErrNoChange {
				log.Printf("no Change")
				return
			}

			if err != nil {
				log.WithFields(log.Fields{"err": err.Error()}).Error("migration failed")
				return
			}
			log.Infof("%d migrations were run successfully", *migrateCount)
		}
	},
}

func setupMigrate() {
	migrateCount = migrateCmd.Flags().IntP("count", "n", 0, "migrate --count [+-]N")
}
