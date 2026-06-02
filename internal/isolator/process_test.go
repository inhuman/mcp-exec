package isolator

import (
	"strings"
	"testing"
)

func TestOutputCapper_TruncatesCombined(t *testing.T) {
	c := &outputCapper{remaining: 10}
	out := c.writer(&c.stdout)
	errw := c.writer(&c.stderr)

	if _, err := out.Write([]byte("aaaaaa")); err != nil { // 6 bytes
		t.Fatal(err)
	}
	if _, err := errw.Write([]byte("bbbbbb")); err != nil { // 6 bytes, only 4 fit
		t.Fatal(err)
	}

	if !c.truncated {
		t.Error("expected truncated=true")
	}
	if total := c.stdout.Len() + c.stderr.Len(); total != 10 {
		t.Errorf("combined kept %d bytes, want 10", total)
	}
	if c.stderr.String() != "bbbb" {
		t.Errorf("stderr = %q, want %q", c.stderr.String(), "bbbb")
	}
}

func TestOutputCapper_NoTruncationUnderLimit(t *testing.T) {
	c := &outputCapper{remaining: 100}
	if _, err := c.writer(&c.stdout).Write([]byte("hello")); err != nil {
		t.Fatal(err)
	}
	if c.truncated {
		t.Error("did not expect truncation")
	}
	if c.stdout.String() != "hello" {
		t.Errorf("stdout = %q", c.stdout.String())
	}
}

func TestSafeUTF8(t *testing.T) {
	in := "ok\xff\xfebad"
	got := safeUTF8(in)
	if !strings.HasPrefix(got, "ok") || strings.ContainsRune(got, 0xff) {
		t.Errorf("safeUTF8 did not sanitize: %q", got)
	}
}
