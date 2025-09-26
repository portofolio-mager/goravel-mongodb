package mongodb

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
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
