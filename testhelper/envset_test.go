package testhelper_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

// makeEnvMap transforms a slice of environment variables into a map
func makeEnvMap(env []string) map[string]string {
	envMap := make(map[string]string)

	for _, e := range env {
		k, v, _ := strings.Cut(e, "=")
		envMap[k] = v
	}

	return envMap
}

// checkEnv checks that the environment matches the suplied list of
// EnvEntry's
func checkEnv(initenv []string, entries []testhelper.EnvEntry) error {
	expEnv := makeEnvMap(initenv)
	for _, ee := range entries {
		expEnv[ee.Key] = ee.Value
	}

	crntEnv := makeEnvMap(os.Environ())

	for k, v := range expEnv {
		if ev := crntEnv[k]; ev != v {
			return fmt.Errorf(
				"Bad envvar %q. Expected %q, got %q",
				k, v, ev)
		}
	}

	if len(expEnv) > len(crntEnv) {
		for k, v := range expEnv {
			if _, ok := crntEnv[k]; !ok {
				return fmt.Errorf(
					"Bad envvar %q. Expected %q, but doesn't exist",
					k, v)
			}
		}
	} else if len(expEnv) < len(crntEnv) {
		for k, v := range crntEnv {
			if _, ok := expEnv[k]; !ok {
				return fmt.Errorf(
					"Bad envvar %q = %q, is not expected",
					k, v)
			}
		}
	}

	return nil
}

func TestEnvCache(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		initEnv      []testhelper.EnvEntry
		newEnv       []testhelper.EnvEntry
		expEnv       []testhelper.EnvEntry
		postResetEnv []testhelper.EnvEntry
	}{
		{
			ID: testhelper.MkID("set and reset"),
			initEnv: []testhelper.EnvEntry{
				{Key: "TestKey_X", Value: "TestVal_X"},
				{Key: "TestKey_Y", Value: "TestVal_Y"},
			},
			newEnv: []testhelper.EnvEntry{
				{Key: "TestKey_X", Value: "TestVal_NewX"},
			},
			expEnv: []testhelper.EnvEntry{
				{Key: "TestKey_X", Value: "TestVal_NewX"},
			},
		},
		{
			ID: testhelper.MkID("set multiple times and reset"),
			initEnv: []testhelper.EnvEntry{
				{Key: "TestKey_X", Value: "TestVal_X"},
				{Key: "TestKey_Y", Value: "TestVal_Y"},
			},
			newEnv: []testhelper.EnvEntry{
				{Key: "TestKey_X", Value: "TestVal_X1"},
				{Key: "TestKey_X", Value: "TestVal_X2"},
				{Key: "TestKey_X", Value: "TestVal_NewX"},
			},
			expEnv: []testhelper.EnvEntry{
				{Key: "TestKey_X", Value: "TestVal_NewX"},
			},
		},
		{
			ID: testhelper.MkID("set and reset with new envvar"),
			initEnv: []testhelper.EnvEntry{
				{Key: "TestKey_X", Value: "TestVal_X"},
				{Key: "TestKey_Y", Value: "TestVal_Y"},
			},
			newEnv: []testhelper.EnvEntry{
				{Key: "TestKey_Z", Value: "TestVal_Z"},
			},
			expEnv: []testhelper.EnvEntry{
				{Key: "TestKey_Z", Value: "TestVal_Z"},
			},
			postResetEnv: []testhelper.EnvEntry{
				{Key: "TestKey_Z"},
			},
		},
	}

	for _, tc := range testCases {
		var ec testhelper.EnvCache

		initenv := os.Environ()

		for _, ee := range tc.initEnv {
			err := os.Setenv(ee.Key, ee.Value)
			if err != nil {
				t.Log(tc.IDStr())
				t.Fatalf("Couldn't set the envvar %q to %q: %s",
					ee.Key, ee.Value, err)
			}
		}

		if err := checkEnv(initenv, tc.initEnv); err != nil {
			t.Log(tc.IDStr())
			t.Fatalf("Couldn't set the initial environment: %s", err)
		}

		initenv = os.Environ()

		err := ec.Setenv(tc.newEnv...)
		testhelper.CheckExpErr(t, err, tc)

		if err != nil {
			continue
		}

		if err := checkEnv(initenv, tc.expEnv); err != nil {
			t.Log(tc.IDStr())
			t.Errorf("EnvCache.Setenv failed: %s", err)

			continue
		}

		ec.ResetEnv()

		if err := checkEnv(initenv, tc.postResetEnv); err != nil {
			t.Log(tc.IDStr())
			t.Errorf("EnvCache.ResetEnv failed: %s", err)
		}
	}
}
