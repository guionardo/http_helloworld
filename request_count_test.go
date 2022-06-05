package main

import (
	"testing"
)

func TestRequestCount_Inc(t *testing.T) {

	t.Run("inc_value", func(t *testing.T) {
		rc := &RequestCount{}
		if rc.Value() != 0 {
			t.Errorf("RequestCount.Value() = %v, want %v", rc.Value(), 0)
		}
		rc.Inc()
		if rc.Value() != 1 {
			t.Errorf("RequestCount.Value() = %v, want %v", rc.Value(), 1)
		}
	})

}
