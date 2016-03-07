package main

import (
	//	"fmt"
	"testing"
)

func TestCamelToKebab(t *testing.T) {
	if v := camelToKebab("describeErrorLogs"); v != "describe-error-logs" {
		//fmt.Printf("%v\n", []rune("describe-error-logs"))
		//fmt.Printf("%v\n", []rune(v))
		t.Errorf("must be describe-error-logs, but %v", v)
	}
	if v := camelToKebab("describeDBInstances"); v != "describe-db-instances" {
		t.Errorf("must be describe-db-instances, but %v", v)
	}
	if v := camelToKebab("describeDBA"); v != "describe-dba" {
		t.Errorf("must be describe-dba, but %v", v)
	}
	if v := camelToKebab("describeDBAMode"); v != "describe-dba-mode" {
		t.Errorf("must be describe-dba-mode, but %v", v)
	}
	if v := camelToKebab("setDRMModeByDBA"); v != "set-drm-mode-by-dba" {
		t.Errorf("must be set-drm-mode-by-dba, but %v", v)
	}
	if v := camelToKebab("DBAModeEnabled"); v != "dba-mode-enabled" {
		t.Errorf("must be dba-mode-enabled, but %v", v)
	}
}

func TestCamelToSnake(t *testing.T) {
	if v := camelToSnake("DBAModeEnabled"); v != "dba_mode_enabled" {
		t.Errorf("must be dba_mode_enabled, but %v", v)
	}
}
