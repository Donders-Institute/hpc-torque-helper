package grpc

import (
	"testing"
)

// TestSecretFromBuild checks if the secretToken variable is assigned at the build time of the code,
// and the value is "my-test-secret".
//
// The test should be run with the build flag:
//
// `-ldflags "-X github.com/Donders-Institute/hpc-torque-helper/internal/grpc.secretToken=my-test-secret"`
func TestSecretFromBuild(t *testing.T) {
	if GetSecret() != "my-test-secret" {
		t.Errorf("secret not passed from the build, expect %s got %s\n", "my-test-secret", GetSecret())
	}
}
