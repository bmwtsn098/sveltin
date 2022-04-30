/**
 * Copyright © 2021-present Sveltin contributors <github@sveltin.io>
 *
 * Use of this source code is governed by Apache 2.0 license
 * that can be found in the LICENSE file.
 */

// Package css ...
package css

import (
	"embed"

	"github.com/spf13/afero"
	"github.com/sveltinio/sveltin/config"
)

// NewTailwindCSS returns a pointer to a CSSLib struct for TailwindCSS.
func NewTailwindCSS(efs *embed.FS, fs afero.Fs, conf *config.SveltinConfig, tplData *config.TemplateData) *CSSLib {
	return &CSSLib{
		Name:    TailwindCSS,
		EFS:     efs,
		FS:      fs,
		Config:  conf,
		TplData: tplData,
	}
}
