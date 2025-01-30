package utils

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func Parse_code(htmlStr string) string {
	start := strings.Index(htmlStr, "<code")
	for ; start != -1; start = strings.Index(htmlStr, "<code") {
		fin := strings.Index(htmlStr[start:], "</code>")
		fin = fin + start
		auxStr := htmlStr[start:fin]

		fin += strings.Index(htmlStr[fin:], ">") + 1

		/*
			lang_start := strings.Index(auxStr, "'")
			lang_fin := strings.Index(auxStr[(start+2):], "'")
			lang_fin = lang_start + lang_fin
			lang := auxStr[(lang_start + 1):(lang_fin + 1)]
		*/
		code_start := strings.Index(auxStr, ">")
		code := auxStr[(code_start + 1):]

		code = strings.Replace(code, "\n", "\\n", 1000)
		r := highlight("cpp", code)
		fmt.Println(r)
		htmlStr = strings.Replace(htmlStr, htmlStr[start:fin], r, 100)
		fmt.Println(htmlStr)
	}
	return htmlStr
}

func highlight(lang, code string) string {
	cmd := exec.Command("/bin/sh", "/home/matei/coalg/utils/hilite.me/run_python.bash", code)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	return string(out)
}
