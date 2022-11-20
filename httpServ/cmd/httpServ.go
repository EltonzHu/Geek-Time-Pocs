package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/golang/glog"
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

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", 8081),
		Handler: mux,
	}

	// Adding proper termination support
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			glog.Fatalln("http server error, ", err.Error())
		}
	}()

	<-quit

	glog.Infoln("shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		glog.Errorln("Server shutdown: ", err)
	} else {
		glog.Infoln("Server has been shutdown succesfully")
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
