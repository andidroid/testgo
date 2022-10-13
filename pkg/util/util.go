package util

import (
	"fmt"
	"os"
)

func CheckErr(err error) {
	if err != nil {
		fmt.Println("errrrror occured", err)

	}
}

func LookupEnv(envKey string, defaultVal string) string {
	env, ok := os.LookupEnv(envKey)
	if !ok {
		env = defaultVal
		os.Setenv(envKey, defaultVal)
	}
	fmt.Printf("read env val %s: %s\n", envKey, env)
	return env
}
