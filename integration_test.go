package mongodb

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/portofolio-mager/goravel-mongodb/contracts"
)

// TestDialectorInitialization tests that our MongoDB dialector initializes correctly
// and provides the necessary connection pool functionality
func TestDialectorInitialization(t *testing.T) {
	tests := []struct {
		name string
		dsn  string
	}{
		{
			name: "valid mongodb dsn",
			dsn:  "mongodb://localhost:27017/testdb",
		},
		{
			name: "empty dsn",
			dsn:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create MongoDB dialector
			dialector := Open(tt.dsn)
			require.NotNil(t, dialector)

			// Test that the dialector has the correct name
			assert.Equal(t, "mongodb", dialector.Name())

			// Test that we can create a connection pool
			connPool, err := newMongoConnPool()
			require.NoError(t, err)
			require.NotNil(t, connPool)

			// Test that the connection pool has a valid embedded sql.DB
			assert.NotNil(t, connPool.DB, "Connection pool should have valid embedded sql.DB")
			assert.IsType(t, (*sql.DB)(nil), connPool.DB)

			// Test that the connection pool implements required methods
			ctx := context.Background()
			assert.NotPanics(t, func() {
				connPool.PrepareContext(ctx, "SELECT 1")
			})
			assert.NotPanics(t, func() {
				connPool.ExecContext(ctx, "SELECT 1")
			})
			assert.NotPanics(t, func() {
				connPool.QueryContext(ctx, "SELECT 1")
			})
			assert.NotPanics(t, func() {
				connPool.QueryRowContext(ctx, "SELECT 1")
			})
		})
	}
}

// TestConnectionPoolInterface verifies our mongoConnPool implements the required interfaces
func TestConnectionPoolInterface(t *testing.T) {
	// Create a connection pool
	pool, err := newMongoConnPool()
	require.NoError(t, err)
	require.NotNil(t, pool)

	// Verify that our pool embeds *sql.DB directly
	assert.NotNil(t, pool.DB, "ConnectionPool should embed valid sql.DB")
	assert.IsType(t, (*sql.DB)(nil), pool.DB)

	// Test that it implements gorm.ConnPool interface methods
	ctx := context.Background()

	// These methods should not panic (they might return errors with our dummy implementation)
	assert.NotPanics(t, func() {
		pool.PrepareContext(ctx, "SELECT 1")
	})

	assert.NotPanics(t, func() {
		pool.ExecContext(ctx, "SELECT 1")
	})

	assert.NotPanics(t, func() {
		pool.QueryContext(ctx, "SELECT 1")
	})

	assert.NotPanics(t, func() {
		pool.QueryRowContext(ctx, "SELECT 1")
	})
}

// TestMongoDBDialectorCompatibility tests the MongoDB dialector's compatibility with GORM
// This test verifies that our implementation works as expected within the GORM framework
func TestMongoDBDialectorCompatibility(t *testing.T) {
	// Create MongoDB dialector
	dialector := Open("mongodb://localhost:27017/testdb")
	require.NotNil(t, dialector)

	// Verify dialector properties
	assert.Equal(t, "mongodb", dialector.Name())

	// Test that our dialector has the ClauseBuilders method (it's on our concrete type)
	if mongoDialector, ok := dialector.(*Dialector); ok {
		assert.NotNil(t, mongoDialector.ClauseBuilders())
	}

	// Test connection pool creation
	connPool, err := newMongoConnPool()
	require.NoError(t, err)
	require.NotNil(t, connPool)

	// Verify the connection pool provides a valid sql.DB
	assert.NotNil(t, connPool.DB, "Connection pool should provide valid sql.DB")
	assert.IsType(t, (*sql.DB)(nil), connPool.DB)

	// Test that the sql.DB instance has the expected properties
	assert.NotPanics(t, func() {
		// These calls should not panic - the actual functionality is handled
		// by our dummy driver, but the interface should be satisfied
		stats := connPool.DB.Stats()
		assert.NotNil(t, stats)
	}, "SQL DB operations should not panic")

	t.Logf("SUCCESS: MongoDB dialector provides valid connection pool with embedded sql.DB")
}

// TestInstanceDBMethod tests the exact framework issue: instance.DB() returning "invalid db"
// This is the critical test that verifies our GetDBConnector implementation works
func TestInstanceDBMethod(t *testing.T) {
	// Create MongoDB dialector - same as framework would do
	dialector := Open("mongodb://localhost:27017/testdb")
	require.NotNil(t, dialector)

	// Initialize GORM instance - same as framework would do
	instance, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err, "GORM initialization should succeed")
	require.NotNil(t, instance)

	// Debug: Check if ConnPool is nil (this is what's causing the issue)
	t.Logf("Debug: instance.ConnPool = %v", instance.ConnPool)
	if instance.ConnPool == nil {
		t.Fatal("PROBLEM FOUND: instance.ConnPool is nil!")
	}

	// This is the exact call that was failing in the framework:
	// db, err := instance.DB()
	// if err != nil { return nil, err } // Was returning "invalid db"
	db, err := instance.DB()

	// These assertions verify our fix works
	assert.NoError(t, err, "instance.DB() should NOT return 'invalid db' error")
	assert.NotNil(t, db, "instance.DB() should return valid *sql.DB, not nil")

	// Verify it's actually a *sql.DB type
	assert.IsType(t, (*sql.DB)(nil), db)

	// Test that we can use the returned sql.DB
	assert.NotPanics(t, func() {
		stats := db.Stats()
		assert.NotNil(t, stats)
	}, "Returned *sql.DB should be usable")

	t.Logf("✅ SUCCESS: instance.DB() returned valid *sql.DB: %p", db)
	t.Logf("✅ Framework issue RESOLVED: No more 'invalid db' error!")
}

// TestFrameworkScenarios tests different ways the framework might initialize GORM
func TestFrameworkScenarios(t *testing.T) {
	scenarios := []struct {
		name   string
		config *gorm.Config
	}{
		{
			name:   "default_config",
			config: &gorm.Config{},
		},
		{
			name:   "nil_config",
			config: nil,
		},
		{
			name: "complex_config",
			config: &gorm.Config{
				DisableForeignKeyConstraintWhenMigrating: true,
				NamingStrategy:                           nil,
			},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			dialector := Open("mongodb://localhost:27017/testdb")

			var instance *gorm.DB
			var err error

			if scenario.config == nil {
				instance, err = gorm.Open(dialector)
			} else {
				instance, err = gorm.Open(dialector, scenario.config)
			}

			require.NoError(t, err, "GORM should initialize successfully")
			require.NotNil(t, instance, "GORM instance should not be nil")

			// Check if ConnPool is set
			t.Logf("Scenario %s: instance.ConnPool = %v", scenario.name, instance.ConnPool)
			if instance.ConnPool == nil {
				t.Errorf("PROBLEM: ConnPool is nil in scenario %s", scenario.name)
				return
			}

			// Test DB() method
			db, err := instance.DB()
			t.Logf("Scenario %s: DB() returned err=%v, db=%p", scenario.name, err, db)

			assert.NoError(t, err, "instance.DB() should work in scenario %s", scenario.name)
			assert.NotNil(t, db, "DB should not be nil in scenario %s", scenario.name)
		})
	}
}

// TestGoravelBuildOrmPattern tests the exact pattern used by Goravel's BuildGorm method
// This reproduces the exact scenario where instance.DB() was failing in the framework
func TestGoravelBuildOrmPattern(t *testing.T) {
	// Simulate Goravel's BuildGorm method pattern

	// Step 1: Create gormConfig similar to Goravel
	gormConfig := &gorm.Config{
		DisableAutomaticPing:                     true,
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		Logger:                                   nil, // Using nil like framework might
	}

	// Step 2: Create dialector (this is pool.Writers[0].Dialector in Goravel)
	dialector := Open("mongodb://localhost:27017/testdb")

	// Step 3: Open GORM instance exactly like Goravel does
	instance, err := gorm.Open(dialector, gormConfig)
	require.NoError(t, err, "gorm.Open should succeed (like in Goravel)")
	require.NotNil(t, instance, "GORM instance should not be nil")

	// Step 4: Test the Ping functionality (Goravel does this)
	if pinger, ok := instance.ConnPool.(interface{ Ping() error }); ok {
		err := pinger.Ping()
		t.Logf("Ping result: %v (this is expected with dummy driver)", err)
	} else {
		t.Logf("ConnPool doesn't support Ping interface")
	}

	// Step 5: This is the CRITICAL part - the exact code from Goravel that was failing
	// if len(pool.Writers) == 1 && len(pool.Readers) == 0 {
	//     db, err := instance.DB()  // THIS WAS FAILING
	//     if err != nil {
	//         return nil, err
	//     }
	// }

	t.Logf("About to call instance.DB() - the critical call that was failing...")

	db, err := instance.DB()

	// These are the assertions that verify our fix for Goravel
	assert.NoError(t, err, "instance.DB() should NOT fail in Goravel BuildGorm pattern")
	assert.NotNil(t, db, "instance.DB() should return valid *sql.DB for Goravel")

	if db != nil {
		// Step 6: Test the connection pool settings that Goravel applies
		assert.NotPanics(t, func() {
			db.SetMaxIdleConns(10)
			db.SetMaxOpenConns(100)
			db.SetConnMaxIdleTime(3600 * time.Second)
			db.SetConnMaxLifetime(3600 * time.Second)
		}, "Should be able to set connection pool settings like Goravel does")

		t.Logf("✅ SUCCESS: Goravel BuildGorm pattern works!")
		t.Logf("✅ instance.DB() returned valid *sql.DB: %p", db)
		t.Logf("✅ Connection pool settings applied successfully")
	}
}

// TestDiagnosticConnPoolNil diagnoses why db.ConnPool might be nil in Goravel
func TestDiagnosticConnPoolNil(t *testing.T) {
	t.Log("=== DIAGNOSTIC: Checking why db.ConnPool might be nil ===")

	// Test 1: Check if dialector Initialize method is being called
	dialector := Open("mongodb://localhost:27017/testdb")

	// Create a custom dialector that tracks if Initialize is called
	type trackingDialector struct {
		*Dialector
	}

	// Custom dialector with tracking
	customDialector := &trackingDialector{Dialector: dialector.(*Dialector)}

	// Test with minimal config (like Goravel might use)
	instance, err := gorm.Open(customDialector, &gorm.Config{
		DisableAutomaticPing: true,
	})

	require.NoError(t, err)
	require.NotNil(t, instance)

	// Test 2: Check current state
	t.Logf("Current db.ConnPool = %v", instance.ConnPool)
	t.Logf("Current db.Config = %v", instance.Config)
	if instance.Config != nil {
		t.Logf("Current db.Config.ConnPool = %v", instance.Config.ConnPool)
	}

	// Test 3: Check Statement
	if instance.Statement != nil {
		t.Logf("Current db.Statement.ConnPool = %v", instance.Statement.ConnPool)
	} else {
		t.Log("db.Statement is nil")
	}

	// Test 4: Try the exact DB() call logic
	t.Log("\n=== Simulating GORM DB() method logic ===")

	connPool := instance.ConnPool
	if instance.Statement != nil && instance.Statement.ConnPool != nil {
		connPool = instance.Statement.ConnPool
		t.Log("Using Statement.ConnPool")
	} else {
		t.Log("Using main ConnPool")
	}

	t.Logf("Final connPool to check = %v", connPool)

	if connPool == nil {
		t.Fatal("❌ PROBLEM CONFIRMED: connPool is nil - this will cause ErrInvalidDB")
	}

	// Test the actual DB() call
	db, err := instance.DB()
	if err != nil {
		t.Fatalf("❌ DB() failed with error: %v", err)
	}

	t.Logf("✅ DB() succeeded, returned: %p", db)
}

// TestPrepareStmtFalseScenario tests the specific issue where PrepareStmt = false
// This reproduces the exact condition that causes connPool to be nil
func TestPrepareStmtFalseScenario(t *testing.T) {
	t.Log("=== TESTING PrepareStmt = false SCENARIO ===")

	// This is the exact scenario that causes the issue
	gormConfig := &gorm.Config{
		PrepareStmt:                              false, // This is the critical setting that causes connPool to be nil
		DisableAutomaticPing:                     true,
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		Logger:                                   nil,
	}

	dialector := Open("mongodb://localhost:27017/testdb")

	t.Log("Opening GORM with PrepareStmt = false...")
	instance, err := gorm.Open(dialector, gormConfig)
	require.NoError(t, err, "GORM should initialize successfully")
	require.NotNil(t, instance, "GORM instance should not be nil")

	// Debug the state
	t.Logf("After GORM initialization:")
	t.Logf("  db.ConnPool = %v", instance.ConnPool)
	t.Logf("  db.Config.ConnPool = %v", instance.Config.ConnPool)

	if instance.Statement != nil {
		t.Logf("  db.Statement.ConnPool = %v", instance.Statement.ConnPool)
	} else {
		t.Log("  db.Statement is nil")
	}

	// This is the critical test - the call that was failing in Goravel
	t.Log("Calling instance.DB()...")
	db, err := instance.DB()

	// Verify our fix works
	assert.NoError(t, err, "instance.DB() should work even with PrepareStmt = false")
	assert.NotNil(t, db, "instance.DB() should return valid *sql.DB")

	if err != nil {
		t.Fatalf("❌ FAILED: instance.DB() returned error: %v", err)
	}

	if db == nil {
		t.Fatal("❌ FAILED: instance.DB() returned nil")
	}

	// Test that it works like Goravel expects
	assert.NotPanics(t, func() {
		db.SetMaxIdleConns(10)
		db.SetMaxOpenConns(100)
		stats := db.Stats()
		assert.NotNil(t, stats)
	}, "Should work with connection pool operations")

	t.Logf("✅ SUCCESS: PrepareStmt = false scenario works!")
	t.Logf("✅ Returned valid *sql.DB: %p", db)
}

// TestPoolDialectorIntegration tests that MongoDB Pool() method now includes Dialector field
// This test verifies the fix that makes MongoDB work like PostgreSQL driver
func TestPoolDialectorIntegration(t *testing.T) {
	t.Log("=== TESTING fullConfigToDialector method directly ===")

	// Create MongoDB instance (need to test method directly due to config complexity)
	mongodb := &MongoDB{}

	// Test fullConfigToDialector method directly
	fullConfig := contracts.FullConfig{
		Config: contracts.Config{
			URI:      "mongodb://localhost:27017/testdb",
			Database: "testdb",
		},
		Connection: "mongodb",
		Driver:     "mongodb",
	}

	// Test the new fullConfigToDialector method
	dialector := mongodb.fullConfigToDialector(fullConfig)
	assert.NotNil(t, dialector, "fullConfigToDialector should return valid dialector")
	assert.Equal(t, "mongodb", dialector.Name(), "Dialector should have correct name")

	// Test that the dialector can be used with GORM (like framework does)
	instance, err := gorm.Open(dialector, &gorm.Config{
		PrepareStmt:                              false, // Test the critical scenario
		DisableAutomaticPing:                     true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	require.NoError(t, err, "GORM should initialize with our dialector")
	require.NotNil(t, instance, "GORM instance should not be nil")

	// The critical test - this should work now
	db, err := instance.DB()
	assert.NoError(t, err, "instance.DB() should work with fullConfigToDialector")
	assert.NotNil(t, db, "instance.DB() should return valid *sql.DB")

	t.Logf("✅ SUCCESS: fullConfigToDialector creates working GORM dialector")
	t.Logf("✅ MongoDB now has PostgreSQL-compatible dialector creation")
	t.Logf("✅ instance.DB() works: %p", db)

	// Test edge case with empty URI (should use default)
	emptyConfig := contracts.FullConfig{
		Config: contracts.Config{
			URI:      "",
			Database: "testdb2",
		},
		Connection: "mongodb",
		Driver:     "mongodb",
	}

	dialector2 := mongodb.fullConfigToDialector(emptyConfig)
	assert.NotNil(t, dialector2, "fullConfigToDialector should work with empty URI")
	assert.Equal(t, "mongodb", dialector2.Name(), "Dialector should still have correct name")

	t.Logf("✅ SUCCESS: fullConfigToDialector handles empty URI correctly")
}
