package tdx

import "testing"

func TestRealTime(t *testing.T) {
	RealTime("sh000001")
}

func TestBatchRealtime(t *testing.T) {
	BatchRealtime([]string{"sh000001", "sh000905", "sz399001", "sz399006", "sh600600", "sz002528"})
}
