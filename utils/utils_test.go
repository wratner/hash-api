package utils

import (
	"encoding/base64"
	"reflect"
	"testing"
)

func TestHashPassword(t *testing.T) {
	actualResult := HashPassword([]byte("angryMonkey"))
	expectedResult, err := base64.StdEncoding.DecodeString("ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actualResult, expectedResult) {
		t.Fatalf("Expected %s but got %s", expectedResult, actualResult)
	}
}

func TestBase64(t *testing.T) {
	actualResult := Base64([]byte("testing"))
	expectedResult := "dGVzdGluZw=="

	if actualResult != expectedResult {
		t.Fatalf("Expected %s but got %s", expectedResult, actualResult)
	}
}
