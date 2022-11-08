package testhelper

import "os"

// EnvEntry records the name and value of an environment variable
type EnvEntry struct {
	Key   string
	Value string
}

// EnvCache maintains a stack of EnvEntry's. It records the previous value of
// each environment variable which has been set by the Setenv method. This
// allows the values to be reset to their original values by the ResetEnv
// method.
type EnvCache struct {
	Stack []EnvEntry
}

// Setenv sets the environment values given by the EnvEntry parameters. It
// records the prior value so that it can be reset later using the ResetEnv
// method. The first failure to set a value returns the error and subsequent
// values are not set.
func (ec *EnvCache) Setenv(entries ...EnvEntry) error {
	for _, ee := range entries {
		val := os.Getenv(ee.Key)
		err := os.Setenv(ee.Key, ee.Value)
		if err != nil {
			return err
		}
		ec.Stack = append(ec.Stack, EnvEntry{Key: ee.Key, Value: val})
	}
	return nil
}

// ResetEnv resets the environment to its state prior to the modifications
// made through the use of the Setenv method. It clears the stack after the
// environment has been restored. Note that the environment is not restored
// exactly as it was; variables which didn't previously exist at all will
// afterwards exist but with an empty value.
func (ec *EnvCache) ResetEnv() {
	for i := len(ec.Stack) - 1; i >= 0; i-- {
		os.Setenv(ec.Stack[i].Key, ec.Stack[i].Value)
	}
	ec.Stack = ec.Stack[:0]
}
