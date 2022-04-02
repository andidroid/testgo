package util

import (
	"fmt"
)

func CheckErr(err error) {
	if err != nil {
		fmt.Println("errrrror occured", err)

	}
}
