package testhelper

import (
	"os"
	"slices"
	"strings"
)

// EnvEntry records the name and value of an environment variable
type EnvEntry struct {
	Key   string
	Value string
}

// EnvCache maintains a stack of [EnvEntry]'s. It records the original value
// of each environment variable before setting the new value using the
// [EnvCache.Setenv] method. This allows the values to be restored to their
// original values by the [EnvCache.ResetEnv] method.
type EnvCache struct {
	Stack []EnvEntry
}

// addEnvEntry adds a new entry to the [EnvCache] Stack. These entries will be
// set in the environment by the [EnvCache.ResetEnv] method.
func (ec *EnvCache) addEnvEntry(key, val string) {
	ec.Stack = append(ec.Stack, EnvEntry{Key: key, Value: val})
}

// Setenv sets the environment values given by the EnvEntry parameters. It
// records the prior value so that it can be restored later using the
// [EnvCache.ResetEnv] method. The first failure to set a value returns the
// error and subsequent values are not set.
func (ec *EnvCache) Setenv(entries ...EnvEntry) error {
	for _, ee := range entries {
		val := os.Getenv(ee.Key)

		err := os.Setenv(ee.Key, ee.Value)
		if err != nil {
			return err
		}

		ec.addEnvEntry(ee.Key, val)
	}

	return nil
}

// Clearenv removes all the environment entries and caches them ready to be
// restored by [EnvCache.ResetEnv]. It behaves like [os.Clearenv] but caches
// the original values so they can be restored.
func (ec *EnvCache) Clearenv() {
	env := os.Environ()
	for _, ekv := range env {
		key, val, _ := strings.Cut(ekv, "=")
		ec.addEnvEntry(key, val)
	}

	os.Clearenv()
}

// ResetEnv resets the environment to its state prior to the modifications
// made through the use of the Setenv method. It clears the stack after the
// environment has been restored. Note that the environment is not restored
// exactly as it was; variables which didn't previously exist at all will
// afterwards exist but with an empty value.
func (ec *EnvCache) ResetEnv() {
	for _, v := range slices.Backward(ec.Stack) {
		_ = os.Setenv(v.Key, v.Value)
	}

	ec.Stack = ec.Stack[:0]
}
