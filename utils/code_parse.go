package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Parse_code(htmlStr string) string {
	start := strings.Index(htmlStr, "<code")
	for ; start != -1; start = strings.Index(htmlStr, "<code") {
		fin := strings.Index(htmlStr[start:], "</code>")
		fin = fin + start
		auxStr := htmlStr[start:fin]

		fin += strings.Index(htmlStr[fin:], ">") + 1

		code_start := strings.Index(auxStr, ">")
		code := auxStr[(code_start + 1):]

		code = strings.Replace(code, "\n", "\\n", 1000)
		code = strings.Replace(code, "\"", "\\\"", 1000)
		fmt.Println(code)
		r := highlight("cpp", code)
		fmt.Println(r)
		htmlStr = strings.Replace(htmlStr, htmlStr[start:fin], r, 100)
		fmt.Println(htmlStr)
	}
	return htmlStr
}

func highlight(lang, code string) string {
	scriptPath := filepath.Join("utils", "hilite.me", "run_python.bash")
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		// Fallback to absolute workspace path if called from another working directory
		scriptPath = "/home/matei/Projects/Coalg/utils/hilite.me/run_python.bash"
	}
	cmd := exec.Command("/bin/sh", scriptPath, code)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("cmd.Run() failed with %s\n", err)
		return "<pre><code>" + code + "</code></pre>"
	}
	return string(out)
}
