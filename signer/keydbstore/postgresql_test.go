// +build postgresqldb

// Initializes a PostgreSQL DB for testing purposes

package keydbstore

import (
	"os"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func init() {
	// Get the PostgreSQL connection string from an environment variable
	dburl := os.Getenv("DBURL")
	if dburl == "" {
		logrus.Fatal("PostgreSQL environment variable not set")
	}

	for i := 0; i < 30; i++ {
		gormDB, err := gorm.Open("postgres", dburl)
		if err == nil {
			err := gormDB.DB().Ping()
			if err == nil {
				break
			}
		}
		if i == 29 {
			logrus.Fatalf("Unable to connect to %s after 60 seconds", dburl)
		}
		time.Sleep(2 * time.Second)
	}

	sqldbSetup = func(t *testing.T) (*SQLKeyDBStore, func()) {
		var cleanup1 = func() {
			gormDB, err := gorm.Open("postgres", dburl)
			require.NoError(t, err)

			// drop all tables, if they exist
			gormDB.DropTable(&GormPrivateKey{})
		}
		cleanup1()
		dbStore := SetupSQLDB(t, "postgres", dburl)

		require.Equal(t, "postgres", dbStore.Name())

		return dbStore, func() {
			dbStore.db.Close()
			cleanup1()
		}
	}
}
