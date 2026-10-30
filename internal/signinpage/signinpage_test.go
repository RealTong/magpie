package signinpage

import (
	"net/http/httptest"
	"strings"
	"testing"
)

// The page is magpie's card with its mark, the outcome's sign and the
// words escaped: a vendor's error text never becomes markup.
func TestWrite(t *testing.T) {
	for _, c := range []struct {
		ok   bool
		sign string
	}{{true, `M5 12.5l4.5 4.5L19 7.5`}, {false, `M7 7l10 10M17 7L7 17`}} {
		rec := httptest.NewRecorder()
		Write(rec, 418, c.ok, "Signed in to <b>x</b>", `refused: <script>alert(1)</script>`)
		body := rec.Body.String()
		if rec.Code != 418 || rec.Header().Get("Content-Type") != "text/html; charset=utf-8" || rec.Header().Get("Cache-Control") != "no-store" {
			t.Fatalf("ok=%v: status %d, headers %v", c.ok, rec.Code, rec.Header())
		}
		if !strings.Contains(body, c.sign) || !strings.Contains(body, `class="card"`) || !strings.Contains(body, "<svg xmlns") {
			t.Fatalf("ok=%v: the card, its sign or magpie's mark is missing:\n%s", c.ok, body)
		}
		if strings.Contains(body, "<script>") || strings.Contains(body, "<b>x") || !strings.Contains(body, "&lt;script&gt;") {
			t.Fatalf("ok=%v: title or message not escaped:\n%s", c.ok, body)
		}
	}
}

// A short hyphened word stays on one line, so "sign-in" never breaks after
// its hyphen; a long id is left to wrap, and the hyphen itself is kept.
func TestHyphenedWordsKeepTogether(t *testing.T) {
	rec := httptest.NewRecorder()
	Write(rec, 200, true, "Signed in", "Every agent uses this sign-in. Client client-0123456789abcdef0123 too.")
	body := rec.Body.String()
	if !strings.Contains(body, `<span class="nb">sign-in</span>`) {
		t.Fatalf("sign-in can break after its hyphen:\n%s", body)
	}
	if strings.Contains(body, `<span class="nb">client-0123`) || !strings.Contains(body, "client-0123456789abcdef0123") {
		t.Fatalf("a long id is held on one line or lost its hyphen:\n%s", body)
	}
}
