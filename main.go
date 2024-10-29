package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func doHello(w http.ResponseWriter, r *http.Request) {
	log.Println("do hello")
}

func doExec(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		return
	}

	cmdStr := r.Form["cmd"][0]
	log.Printf("cmd is: %v", cmdStr)

	if err := os.Setenv("PATH", "/home/user/.goenv/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"); err != nil {
		log.Fatalf("os setenv error: %v", err.Error())
	}

	commands := strings.Fields(cmdStr)
	cmd := exec.Command(commands[0], commands[1:]...)

	// 捕获输出
	var wg sync.WaitGroup

	if stdout, err := cmd.StdoutPipe(); err != nil {
		log.Panicf("StdoutPipe Error: %v", err.Error())
	} else {
		defer stdout.Close()
		wg.Add(1)
		go func(p io.ReadCloser) {
			defer wg.Done()
			reader := bufio.NewReader(stdout)
			output(reader)
		}(stdout)
	}
	if stderr, err := cmd.StderrPipe(); err != nil {
		log.Panicf("StdoutPipe ERROR: %v", err.Error())
	} else {
		defer stderr.Close()
		wg.Add(1)
		go func(p io.ReadCloser) {
			defer wg.Done()
			reader := bufio.NewReader(stderr)
			output(reader)
		}(stderr)
	}

	// run
	if err := cmd.Run(); err != nil {
		log.Printf("exec failed: %v", err.Error())
		_, _ = w.Write([]byte(err.Error()))
	} else {
		log.Println("exec success")
	}
	wg.Wait()
}

func main() {
	http.HandleFunc("/hello", doHello)
	http.HandleFunc("/exec", doExec)

	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Printf("listen and serve failed: %v", err)
		return
	}
}

func output(reader *bufio.Reader) {
	line, err := reader.ReadString('\n')
	for err == nil {
		log.Println(strings.TrimSpace(line))
		line, err = reader.ReadString('\n')
	}
}
