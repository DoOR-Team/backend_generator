package main

import (
	"fmt"
	"log"
	"testing"
)

func getType1(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

func TestStr(t *testing.T) {
	initForbiddenChar()

	log.Println(checkName("?aaabb-b"))

}
