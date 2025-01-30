package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func Parse_code(htmlStr string) string {
	//for start := strings.Index(htmlStr, "<code"); start != -1; {
	start := strings.Index(htmlStr, "<code")
	fin := strings.Index(htmlStr[start:], "</code>")
	fin = fin + start
	auxStr := htmlStr[start:fin]

	fin += strings.Index(htmlStr[fin:], ">") + 1

	//fmt.Println(htmlStr[start:fin])

	lang_start := strings.Index(auxStr, "\"")
	lang_fin := strings.Index(auxStr[(start+2):], "\"")
	lang_fin = lang_start + lang_fin
	lang := auxStr[(lang_start + 1):(lang_fin + 1)]

	code_start := strings.Index(auxStr, ">")
	code := auxStr[(code_start + 1):]

	r := highlight(lang, code)
	htmlStr = strings.Replace(htmlStr, htmlStr[start:fin], r, 100)
	fmt.Println(htmlStr)
	//}
	return htmlStr
}

func highlight(lang, code string) string {
	my_url, _ := url.Parse("http://hilite.me/api?")
	param, _ := url.ParseQuery(my_url.RawQuery)

	param.Add("code", code)
	param.Add("lexer", lang)
	param.Add("style", "monokai")

	r, err := http.Get("http://hilite.me/api?" + param.Encode())
	if err != nil {
		panic(err)
	}
	body, _ := io.ReadAll(r.Body)
	return string(body)
}
