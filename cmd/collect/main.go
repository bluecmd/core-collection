package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <pid of dying process>", os.Args[0])
	}
	pid, err := strconv.ParseInt(os.Args[1], 10, 32)
	if err != nil {
		log.Fatalf("Unable to parse dying process PID: %v", err)
	}

	type CorePayload struct {
		Hostname string `json:"name"`
		PID      int    `json:"pid"`
		Cmdline  string `json:"cmdline"`
		CoreDump []byte `json:"core"`
	}

	cmdline, err := ioutil.ReadFile(path.Join("/proc", fmt.Sprintf("%d", pid), "cmdline"))
	if err != nil {
		log.Fatalf("Unable to read cmdline from dying process: %v", err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "<unknown>"
	}

	body := &CorePayload{
		Hostname: hostname,
		PID:      int(pid),
		Cmdline:  string(cmdline),
	}

	coreBuf := new(bytes.Buffer)
	gzw := gzip.NewWriter(coreBuf)
	if _, err := io.Copy(gzw, os.Stdin); err != nil {
		log.Fatalf("Failed to compress core file: %v", err)
	}
	if err := gzw.Close(); err != nil {
		log.Fatalf("Failed to close compressed core: %v", err)
	}

	body.CoreDump = coreBuf.Bytes()
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)
	if _, err := http.Post("http://speedtest.kamel.network:8888", "application/json", payloadBuf); err != nil {
		log.Fatalf("Error: %+v", err)
	}
}
