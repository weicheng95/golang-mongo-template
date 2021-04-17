package helper

import (
	"fmt"

	"github.com/spf13/viper"
)

// Load read the env filename and load it into ENV for this process.
func Load(env string, filename string) error {
	// load variables from system environment
	if env != "production" {
		viper.SetConfigFile(filename)
		viper.SetConfigType("yaml")
		err := viper.ReadInConfig()
		if err != nil {
			return fmt.Errorf("Error while reading config file %s", err)
		}
	}

	viper.AutomaticEnv()
	return nil
}

func Get(key string) (string, error) {
	// viper.Get() returns an empty interface{}
	// to get the underlying type of the key,
	// we have to do the type assertion, we know the underlying value is string
	// if we type assert to other type it will throw an error
	value, _ := viper.Get(key).(string)

	// If the type is a string then ok will be true
	// ok will make sure the program not break
	// if !ok {
	// 	return "", fmt.Errorf("Invalid type assertion")
	// }

	if value == "" {
		return "", fmt.Errorf("Couldn't get configuration value for %s")
	}

	return value, nil
}
