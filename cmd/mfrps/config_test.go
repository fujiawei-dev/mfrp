/*
 * @Date: 2022.01.22 18:49
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.01.22 18:49
 */

package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfig_Load(t *testing.T) {
	if err := conf.Load(filepath.Join(os.TempDir(), "mfrps.yaml")); err != nil {
		t.Error(err)
	}
}
