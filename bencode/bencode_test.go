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
	result := EncodeDictionary([]string{EncodeKV("key", "value")})

	if result != expectedResult {
		t.Fatalf("Expected %s, got %s", expectedResult, result)
	}
}

func TestEncodeDictionarySubList(t *testing.T) {
	expectedResult := "d3:keyl5:value4:testee"
	result := EncodeDictionary([]string{EncodeKV("key", EncodeList([]string{"value", "test"}))})

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

func TestEncodePeerList(t *testing.T) {
	expectedResult := "d5:peersld7:peer id20:111111111111111111112:ip9:127.0.0.14:port4:1276eee"
	result := EncodePeerList([]string{"127.0.0.1:1276"})

	if result != expectedResult {
		t.Fatalf("Expected %s, got %s", expectedResult, result)
	}
}

func TestEncodeKV(t *testing.T) {
	expectedResult := "7:peer id1:a"
	result := EncodeKV("peer id", "a")

	if result != expectedResult {
		t.Fatalf("Expected %s, got %s", expectedResult, result)
	}
}

func TestEncodeKVValueIsInt(t *testing.T) {
	expectedResult := "7:peer idi10e"
	result := EncodeKV("peer id", "i10e")

	if result != expectedResult {
		t.Fatalf("Expected %s, got %s", expectedResult, result)
	}
}

func TestWriteStringData(t *testing.T) {
	expectedResult := "test1234"
	result := writeStringData("test", "1234")

	if expectedResult != result {
		t.Fatalf("Expected %s, got %s", expectedResult, result)
	}
}
