package main

import (
	"bytes"
	"flag"
	"gopkg.in/ini.v1"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	nTarget = flag.String("t", "/etc/supervisord.d", "需要拆分的ini目录")
	nDelete = flag.Bool("d", false, "是否删除源文件")
)

func main() {
	flag.Parse()
	err := filepath.Walk(*nTarget, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".ini") {
			return nil
		}
		err2 := splitSingleIni(path, *nTarget)
		if err2 == nil {
			if *nDelete {
				return os.Remove(path)
			}
		}
		return err2
	})
	if err != nil {
		log.Fatal(err)
	}
}

func splitSingleIni(filename string, outputDir string) error {
	cfg, err := ini.Load(filename)
	if err != nil {
		log.Fatal(err)
	}
	for _, section := range cfg.Sections() {
		buf := bytes.Buffer{}
		buf.WriteString("[" + section.Name() + "]\n")
		for k, v := range section.KeysHash() {
			buf.WriteString(k + "=" + v + "\n")
		}
		filename := strings.ReplaceAll(section.Name(), ":", "_")
		filename = path.Join(outputDir, filename+".ini")
		err := ioutil.WriteFile(filename, buf.Bytes(), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
