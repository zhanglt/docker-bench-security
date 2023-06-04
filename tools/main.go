/**
*翻译tests目录中的条目
 */
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	translator "github.com/Conight/go-googletrans"
)

func t(text string, t *translator.Translator) string {

	result, err := t.Translate(text, "en", "zh")
	if err != nil {
		panic(err)
	}

	return result.Text

}

func main() {
	//需要翻译的文件列表
	var tmpl []string
	tmpl = append(tmpl, "../tests/1_host_configuration")
	tmpl = append(tmpl, "../tests/2_docker_daemon_configuration")
	tmpl = append(tmpl, "../tests/3_docker_daemon_configuration_files")
	tmpl = append(tmpl, "../tests/4_container_images")
	tmpl = append(tmpl, "../tests/5_container_runtime")
	tmpl = append(tmpl, "../tests/6_docker_security_operations")
	tmpl = append(tmpl, "../tests/7_docker_swarm_configuration")
	tmpl = append(tmpl, "../tests/8_docker_enterprise_configuration")
	tmpl = append(tmpl, "../tests/99_community_checks")
	c := translator.Config{
		Proxy: "http://127.0.0.1:10809",
	}
	ts := translator.New(c)

	for _, file := range tmpl {

		err := ReadLines(file+".sh", file+"_zh.sh", ts)
		if err != nil {
			fmt.Printf("%s文件处理错误:%s", file, err)
		}
	}

}

func ReadLines(inFile, outFile string, ts *translator.Translator) error {
	in, err := os.Open(inFile)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer out.Close()
	write := bufio.NewWriter(out)
	read := bufio.NewReader(in)
	for {
		bytes, _, err := read.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		s := string(bytes)
		tmp := getStr(s, "local desc=", "local remediation=", "local remediationImpact=", ts)
		if tmp != "" {
			s = tmp
		}
		write.WriteString(fmt.Sprintf("%s\n", s))
	}

	err = write.Flush()
	return err
}
func getStr(str, substr1, substr2, substr3 string, ts *translator.Translator) string {

	if strings.Contains(str, substr1) {
		return substring(str, substr1, ts)
	}
	if strings.Contains(str, substr2) {
		return substring(str, substr2, ts)
	}
	if strings.Contains(str, substr3) {
		return substring(str, substr3, ts)
	}
	return ""
}

func substring(str, substr string, ts *translator.Translator) string {
	var s string
	l := len(str)
	i := strings.Index(str, substr)
	s = str[i+len(substr)+1 : l-1]
	if s != "" && s != "None." {
		s = t(s, ts)
	}
	s = str[0:i+len(substr)+1] + s + str[l-1:l]
	return s
}
