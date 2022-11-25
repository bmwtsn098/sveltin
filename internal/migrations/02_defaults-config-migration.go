/**
 * Copyright © 2021-present Sveltin contributors <github@sveltin.io>
 *
 * Use of this source code is governed by Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package migrations

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"github.com/sveltinio/sveltin/common"
)

// UpdateDefaultsConfigMigration is the struct representing the migration update the defaults.js.ts file.
type UpdateDefaultsConfigMigration struct {
	Mediator IMigrationMediator
	Services *MigrationServices
	Data     *MigrationData
}

// MakeMigration implements IMigrationFactory interface,
func (m *UpdateDefaultsConfigMigration) MakeMigration(migrationManager *MigrationManager, services *MigrationServices, data *MigrationData) IMigration {
	return &UpdateDefaultsConfigMigration{
		Mediator: migrationManager,
		Services: services,
		Data:     data,
	}
}

// implements IMigration interface.
func (m *UpdateDefaultsConfigMigration) getServices() *MigrationServices { return m.Services }
func (m *UpdateDefaultsConfigMigration) getData() *MigrationData         { return m.Data }

// Execute return error if migration execution over up and down methods fails (IMigration interface).
func (m UpdateDefaultsConfigMigration) Execute() error {
	if err := m.up(); err != nil {
		return err
	}
	if err := m.down(); err != nil {
		return err
	}
	return nil
}

func (m *UpdateDefaultsConfigMigration) up() error {
	if !m.Mediator.canRun(m) {
		return nil
	}

	exists, err := common.FileExists(m.getServices().fs, m.Data.FileToMigrate)
	if err != nil {
		return err
	}

	if exists {
		fileContent, err := retrieveFileContent(m.getServices().fs, m.getData().FileToMigrate)
		if err != nil {
			return err
		}

		migrationTriggers := []string{patterns[semVersion]}
		if isMigrationRequired(fileContent, migrationTriggers, findStringMatcher) {
			m.getServices().logger.Info(fmt.Sprintf("Migrating %s", filepath.Base(m.Data.FileToMigrate)))
			if err := updateConfigFile(m, fileContent); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *UpdateDefaultsConfigMigration) down() error {
	if err := m.Mediator.notifyAboutCompletion(); err != nil {
		return err
	}
	return nil
}

func (m *UpdateDefaultsConfigMigration) allowUp() error {
	if err := m.up(); err != nil {
		return err
	}
	return nil
}

//=============================================================================

func updateConfigFile(m *UpdateDefaultsConfigMigration, content []byte) error {
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		rules := []*migrationRule{newSveltinVersionRule(line)}
		if res, ok := applyMigrationRules(rules); ok {
			lines[i] = res
		} else {
			lines[i] = line
		}
	}
	output := strings.Join(lines, "\n")
	err := m.getServices().fs.Remove(m.Data.FileToMigrate)
	if err != nil {
		return err
	}

	if err = afero.WriteFile(m.getServices().fs, m.Data.FileToMigrate, []byte(output), 0644); err != nil {
		return err
	}
	return nil
}

//=============================================================================

func newSveltinVersionRule(line string) *migrationRule {
	return &migrationRule{
		value:           line,
		pattern:         patterns[semVersion],
		replaceFullLine: true,
		replacerFunc: func(string) string {
			return `import { sveltin } from '../sveltin.json';

const sveltinVersion = sveltin.version;`
		},
	}
}
