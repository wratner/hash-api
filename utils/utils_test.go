package utils

import (
	"encoding/base64"
	"reflect"
	"testing"
)

//TestHashPassword tests the given test case as part of the assignment.
func TestHashPassword(t *testing.T) {
	actualResult, err := HashPassword([]byte("angryMonkey"))
	if err != nil {
		t.Fatal(err)
	}
	expectedResult, err := base64.StdEncoding.DecodeString("ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actualResult, expectedResult) {
		t.Fatalf("Expected %s but got %s", expectedResult, actualResult)
	}
}

//TestBase64 tests that the base64 function gives the correct output.
func TestBase64(t *testing.T) {
	actualResult := Base64([]byte("testing"))
	expectedResult := "dGVzdGluZw=="

	if actualResult != expectedResult {
		t.Fatalf("Expected %s but got %s", expectedResult, actualResult)
	}
}
