package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
)

const undefined = "undefined, build process will set it"

var (
	// Version of Go toolchain that produces this build
	goVersion = runtime.Version()
	// Platform for which this build was produced
	goPlatform = fmt.Sprintf("%s:%s", runtime.GOOS, runtime.GOARCH)

	// // Version of this build
	// srcVersion = undefined
	// SHA string of the scm commit of this build
	commitHash = undefined
	// SCM branch of this build
	scmBranch = undefined
	// Date this build was produced
	buildDate = undefined
)

func buildInfo() {
	fmt.Printf("\nStarting Http Server,"+
		"Go version: %s\n"+
		"Go platform: %s\n"+
		"Commit hash: %s\n"+
		"Scm branch: %s\n"+
		"Build date: %s\n\n",
		goVersion, goPlatform, commitHash, scmBranch, buildDate)
}

func main() {
	// For simplity, use console output for now instead of using glog lib dependencies.
	// - Elton Hu (Oct.8, 2022)
	// flag.Set("v", 1)
	// glog.V(2).Info("Starting http server")
	buildInfo()

	// In future, we can split this into differnt service layers from the mux.
	// - Elton Hu (Oct.8, 2022)
	mux := http.NewServeMux()
	mux.HandleFunc("/tasks/v1/copyHeaders", copyHeadersHandler)
	mux.HandleFunc("/tasks/v1/copyEnvVersion", copyEnvVersionHanlder)
	mux.HandleFunc("/tasks/v1/log", logHandler)
	mux.HandleFunc("/healthz", healthProbHandler)
	err := http.ListenAndServe(":8081", mux)

	if err != nil {
		log.Fatal(err)
	}

}

// Service handler section
func copyHeadersHandler(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		for _, value := range v {
			w.Header().Add(k, value)
		}
	}
}

func copyEnvVersionHanlder(w http.ResponseWriter, r *http.Request) {
	if versionNum, exist := os.LookupEnv("VERSION"); exist {
		fmt.Println(versionNum)
		w.Header().Add("Env-Version", versionNum)
	}
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Println(r.RemoteAddr + " return code: " + fmt.Sprint(http.StatusOK))
}

func healthProbHandler(w http.ResponseWriter, r *http.Request) {
	return
}
