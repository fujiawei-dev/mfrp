/*
 * @Date: 2022.01.22 18:49
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.01.22 18:49
 */

package main

import "testing"

func TestConfig_Load(t *testing.T) {
	if err := conf.Load("mfrpc.yaml"); err != nil {
		t.Error(err)
	}
}
