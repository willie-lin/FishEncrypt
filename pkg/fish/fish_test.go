package fish

import (
	"fmt"
	"testing"
)

func TestBase64EncodeMod8(t *testing.T) {
	src := "egg spam"
	expected := "H34qN/uqQnz/"

	encr := Base64Encode([]byte(src))
	fmt.Println(encr)
	fmt.Println(111)

	if encr != expected {
		t.Fatalf("%s != %s", expected, encr)
	}
}

func TestBase64Encode(t *testing.T) {

	src := "The quick brown fox jumps over the lazy dog"
	expected := "xzkrL/ui4oi/uSQrJ/M746F/KPkrE/uuRpA/O/wqz/QX46N/uyBsv/G/" +
		"gnC/.......qQpy/"

	enc := Base64Encode([]byte(src))
	fmt.Println(enc)
	if enc != expected {
		t.Fatalf("%s != %s", expected, enc)
	}

	if len(Base64Encode([]byte(""))) != 0 {
		t.Fatalf("Expected empty string")
	}
}



