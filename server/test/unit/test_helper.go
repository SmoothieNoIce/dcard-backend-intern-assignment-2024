package unit

import (
	"fmt"
	"github/SmoothieNoIce/dcard-backend-intern-assignment-2024/app/models"
	"github/SmoothieNoIce/dcard-backend-intern-assignment-2024/config"
	"github/SmoothieNoIce/dcard-backend-intern-assignment-2024/database/cache"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"gorm.io/gorm"
)

func PerformRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// 改良 https://jarifibrahim.github.io/blog/test-cleanup-with-gorm-hooks/
// 這個 hook 須在每個測試完成後執行，以清除資料
func DeleteCreatedEntities(db *gorm.DB) func() {
	type entity struct {
		table   string
		keyname string
		key     interface{}
	}
	var entries []entity
	hookName := "cleanupHook"

	db.Callback().Create().After("gorm:create").Register(hookName, func(dba *gorm.DB) {
		_, _ = dba.Statement.Schema.PrioritizedPrimaryField.ValueOf(dba.Statement.ReflectValue)
		if value, isZero := dba.Statement.Schema.PrioritizedPrimaryField.ValueOf(dba.Statement.ReflectValue); !isZero {
			fmt.Printf("Inserted entities of %s with %s=%v\n", dba.Statement.Table, dba.Statement.Schema.PrioritizedPrimaryField.Name, value)
			entries = append(entries, entity{table: dba.Statement.Table, keyname: dba.Statement.Schema.PrioritizedPrimaryField.Name, key: value})
		}
	})
	return func() {
		// Remove the hook once we're done
		defer db.Callback().Create().Remove(hookName)
		tx := db

		// Loop from the end. It is important that we delete the entries in the
		// reverse order of their insertion
		for i := len(entries) - 1; i >= 0; i-- {
			entry := entries[i]
			fmt.Printf("Deleting entities from '%s' table with key %v\n", entry.table, entry.key)
			tx.Table(entry.table).Where(entry.keyname+" = ?", entry.key).Delete("")
		}
		tx.Commit()
	}
}

func SetupIntegrationDB(t *testing.T) func() {
	// setup the connection, use t.Fatal for all errors.
	config.Setup("../../../config.json.test")
	db := models.Setup(true, true)
	cache.SetUpDefaultDB()
	cleaner := DeleteCreatedEntities(db)

	cleanup := func() {
		// here you have access to both t, and *sql.DB
		// you can do all the clean up required,
		// and return this anonymous function to be called later
		cache.Rdb0.FlushAll(cache.Ctx)
		cleaner()
	}

	return cleanup
}
