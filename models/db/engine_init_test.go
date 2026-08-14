// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package db

import (
	"context"
	"path/filepath"
	"testing"

	"gitea.dev/modules/setting"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type autoMigrationTestModel struct {
	ID int64 `xorm:"pk autoincr"`
}

func TestInitEngineWithMigrationSkipsSchemaSyncWhenAutoMigrationDisabled(t *testing.T) {
	originalDatabase := setting.Database
	originalEngine := xormEngine
	originalModels := registeredModels
	originalInitFuncs := registeredInitFuncs
	t.Cleanup(func() {
		UnsetDefaultEngine()
		setting.Database = originalDatabase
		xormEngine = originalEngine
		registeredModels = originalModels
		registeredInitFuncs = originalInitFuncs
	})

	ResetModels()
	RegisterModel(new(autoMigrationTestModel))
	setting.Database.Type = setting.DatabaseTypeSQLite3
	setting.Database.Path = filepath.Join(t.TempDir(), "auto-migration-disabled.db")
	setting.Database.AutoMigration = false

	migrationChecked := false
	require.NoError(t, InitEngineWithMigration(t.Context(), func(context.Context, EngineMigration) error {
		migrationChecked = true
		return nil
	}))
	exists, err := xormEngine.IsTableExist(new(autoMigrationTestModel))
	require.NoError(t, err)
	assert.True(t, migrationChecked)
	assert.False(t, exists)
}
