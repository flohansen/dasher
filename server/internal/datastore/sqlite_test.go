package datastore

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"

	"github.com/flohansen/dasher-server/internal/sqlc"
	"github.com/stretchr/testify/assert"

	_ "github.com/mattn/go-sqlite3"
)

const (
	testDbFileName = "./test.db"
)

func TestSQLiteDatastore_Integration(t *testing.T) {
	t.Run("given an empty database", func(t *testing.T) {
		db := createTestDatabase(t)
		store := NewSQLite(db.db)

		db.init()

		t.Run("when adding a feature toggle", func(t *testing.T) {
			err := store.Upsert(context.TODO(), sqlc.Feature{
				FeatureID: "FEATURE_TOGGLE_TEST",
			})

			t.Run("there should be no error", func(t *testing.T) {
				assert.NoError(t, err)
			})

			t.Run("the database should contain a new item", func(t *testing.T) {
				assert.Equal(t, &sqlc.Feature{
					FeatureID: "FEATURE_TOGGLE_TEST",
				}, db.getFeature("FEATURE_TOGGLE_TEST"))
			})
		})
	})

	t.Run("given a database with a feature toggle", func(t *testing.T) {
		db := createTestDatabase(t)
		store := NewSQLite(db.db)

		db.init()
		db.insertFeature("FEATURE_TOGGLE_TEST")

		t.Run("when getting all feature toggles", func(t *testing.T) {
			features, err := store.GetAll(context.TODO())

			t.Run("there should be no error", func(t *testing.T) {
				assert.NoError(t, err)
			})

			t.Run("the number of features should be 1", func(t *testing.T) {
				assert.Len(t, features, 1)
			})
		})

		t.Run("when adding a feature toggle with the same ID", func(t *testing.T) {
			err := store.Upsert(context.TODO(), sqlc.Feature{
				FeatureID: "FEATURE_TOGGLE_TEST",
				Enabled:   true,
			})

			t.Run("there should be no error", func(t *testing.T) {
				assert.NoError(t, err)
			})

			t.Run("the existing feature toggle should be updated", func(t *testing.T) {
				feature := db.getFeature("FEATURE_TOGGLE_TEST")
				assert.NotNil(t, feature)
				assert.Equal(t, "FEATURE_TOGGLE_TEST", feature.FeatureID)
				assert.True(t, feature.Enabled)
			})
		})

		t.Run("when deleting the feature toggle", func(t *testing.T) {
			err := store.Delete(context.TODO(), "FEATURE_TOGGLE_TEST")

			t.Run("there should be no error", func(t *testing.T) {
				assert.NoError(t, err)
			})

			t.Run("the database should not contain the item anymore", func(t *testing.T) {
				assert.Nil(t, db.getFeature("FEATURE_TOGGLE_TEST"))
			})
		})
	})
}

type testDb struct {
	t  *testing.T
	db *sql.DB
}

func createTestDatabase(t *testing.T) *testDb {
	db, err := sql.Open("sqlite3", testDbFileName)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		db.Close()
		os.Remove(testDbFileName)
	})

	return &testDb{t, db}
}

func (db *testDb) init() {
	if _, err := db.db.Exec(`
		create table if not exists features (
			feature_id text not null primary key,
			description text not null,
			enabled boolean not null
		)
	`); err != nil {
		db.t.Fatal(err)
	}
}

func (db *testDb) insertFeature(id string) {
	if _, err := db.db.Exec(`
		insert into features (feature_id, description, enabled)
		values (?, '', 0)
	`, id); err != nil {
		db.t.Fatal(err)
	}
}

func (db *testDb) getFeature(id string) *sqlc.Feature {
	res := db.db.QueryRow(`
		select feature_id,
			   description,
			   enabled
		from features
		where feature_id = ?
	`, id)
	if err := res.Err(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		db.t.Fatal(err)
	}

	var feature sqlc.Feature
	if err := res.Scan(
		&feature.FeatureID,
		&feature.Description,
		&feature.Enabled,
	); err != nil {
		return nil
	}

	return &feature
}
