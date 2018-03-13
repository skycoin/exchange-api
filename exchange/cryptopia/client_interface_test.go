// +build cryptopia_integration_test

package cryptopia_test

import (
	"os"
)

var key, secret = func() (key string, secret string) {
	var found bool
	if key, found = os.LookupEnv("CRYPTOPIA_TEST_KEY"); found {
		if secret, found = os.LookupEnv("CRYPTOPIA_TEST_SECRET"); found {
			return
		}
		panic("Cryptopia secret not provided")
	}
	panic("Cryptopia key not provided")
}()
