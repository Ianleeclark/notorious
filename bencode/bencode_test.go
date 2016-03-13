package bencode

import (
	"testing"
)

func TestEncodeInt(t *testing.T) {
	expectedResult := "i5e"
	result := EncodeInt(5)

	if result != expectedResult {
		t.Fatalf("Expected %s, got %s", expectedResult, result)
	}
}

func TestEncodeList(t *testing.T) {
	expectedResult := "l4:TESTe"
	a := []string{"TEST"}
	result := EncodeList(a)

	if result != expectedResult {
		t.Fatalf("Expected %s, got %s", expectedResult, result)
	}
}

func TestEncodeDictionary(t *testing.T) {
	expectedResult := "d3:key5:valuee"
	result := EncodeDictionary("3:key", "5:value")

	if result != expectedResult {
		t.Fatalf("Expected %s, got %s", expectedResult, result)
	}
}

func TestEncodeDictionarySubList(t *testing.T) {
	expectedResult := "d3:keyl5:value4:testee"
	result := EncodeDictionary("3:key", "l5:value4:teste")

	if result != expectedResult {
		t.Fatalf("Expected %s, got %s", expectedResult, result)
	}
}

func TestEncodeByteString(t *testing.T) {
	expectedResult := "4:TEST"
	result := EncodeByteString("TEST")

	if result != expectedResult {
		t.Fatalf("Expected %s, got %s", expectedResult, result)
	}
}
