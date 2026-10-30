// Package signinpage is the page a browser lands on at the end of a sign-in
// magpie runs: a provider account's and an MCP server's alike. It is
// magpie's, in its own quiet black and white, never the vendor's "return
// to the CLI".
package signinpage

import (
	_ "embed"
	"fmt"
	"html"
	"net/http"
	"regexp"
)

// logo is magpie's mark (build/icon/magpie-small.svg), signing the page.
//
//go:embed magpie.svg
var logo string

// hyphened is a short word with a hyphen in it, such as sign-in, which a
// browser would otherwise break after its hyphen. A long one (an id) is
// left free to wrap, so it can't push out of the card.
var hyphened = regexp.MustCompile(`[\p{L}\p{N}]+(?:-[\p{L}\p{N}]+)+`)

// text escapes s and keeps each short hyphened word on one line.
func text(s string) string {
	return hyphened.ReplaceAllStringFunc(html.EscapeString(s), func(w string) string {
		if len(w) > 20 {
			return w
		}
		return `<span class="nb">` + w + `</span>`
	})
}

// Write answers with the page: a check when ok, a cross when not, the
// title over msg. Both are escaped.
func Write(w http.ResponseWriter, status int, ok bool, title, msg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	mark := `<path d="M5 12.5l4.5 4.5L19 7.5"/>`
	if !ok {
		mark = `<path d="M7 7l10 10M17 7L7 17"/>`
	}
	fmt.Fprintf(w, `<!doctype html><html><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>magpie · %[1]s</title><style>
:root{--bg:#fafafa;--fg:#111;--muted:#6b6b6b;--line:#e4e4e4;--card:#fff}
@media (prefers-color-scheme:dark){:root{--bg:#0e0e0e;--fg:#f2f2f2;--muted:#9a9a9a;--line:#262626;--card:#161616}}
*{box-sizing:border-box}html,body{height:100%%;margin:0}
body{background:var(--bg);color:var(--fg);font:15px/1.5 -apple-system,BlinkMacSystemFont,"Segoe UI",system-ui,sans-serif;display:grid;place-items:center;padding:16px}
.card{width:100%%;max-width:380px;background:var(--card);border:1px solid var(--line);border-radius:16px;padding:32px 28px;text-align:center;animation:in .4s cubic-bezier(.2,.8,.2,1)}
.mark{width:44px;height:44px;margin:0 auto 18px;border-radius:50%%;display:grid;place-items:center;background:var(--fg);color:var(--bg)}
.mark svg{width:22px;height:22px;fill:none;stroke:currentColor;stroke-width:2.2;stroke-linecap:round;stroke-linejoin:round}
h1{font-size:18px;font-weight:620;margin:0 0 6px;letter-spacing:-.01em}p{margin:0;color:var(--muted);font-size:14px;overflow-wrap:anywhere}
.by{margin-top:24px;display:flex;align-items:center;justify-content:center;gap:7px;font-size:12.5px;font-weight:560;color:var(--muted);letter-spacing:.01em}
.by svg{width:20px;height:20px;fill:currentColor}
.nb{white-space:nowrap}
@media (prefers-reduced-motion:reduce){.card{animation:none}}
@keyframes in{from{opacity:0;transform:translateY(6px) scale(.98)}}
</style></head><body><div class="card"><div class="mark"><svg viewBox="0 0 24 24">%[4]s</svg></div>
<h1>%[2]s</h1><p>%[3]s</p><div class="by">%[5]smagpie</div></div></body></html>`, html.EscapeString(title), text(title), text(msg), mark, logo)
}
