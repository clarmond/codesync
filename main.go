package main

import (
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
)

const (
	defaultConcurrency = 16
	eventBufferSize    = 4096
)

func main() {
	concurrency := flag.Int("concurrency", defaultConcurrency, "number of concurrent copy/delete workers")
	flag.Usage = func() { syntax(false) }
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		syntax(true)
	}

	watchDir, err := filepath.Abs(args[0])
	if err != nil {
		fatal("%s: %v", args[0], err)
	}
	destDir, err := filepath.Abs(args[1])
	if err != nil {
		fatal("%s: %v", args[1], err)
	}

	if !isDirectory(watchDir) {
		syntax(false)
		fatal("%s not found or is not a directory", watchDir)
	}
	if !isDirectory(destDir) {
		syntax(false)
		fatal("%s not found or is not a directory", destDir)
	}

	ig := loadIgnoreMatcher(watchDir)
	s := &syncer{watchDir: watchDir, destDir: destDir, ignore: ig}

	events := make(chan syncEvent, eventBufferSize)
	wg := startWorkers(s, events, *concurrency)

	w, err := newFileWatcher(watchDir, ig, events)
	if err != nil {
		fatal("Failed to start watcher: %v", err)
	}

	logInfo("Monitoring %s for changes\n", watchDir)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		logInfo("\nStopping.")
		os.Exit(0)
	}()

	w.run()
	close(events)
	wg.Wait()
}

func isDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func syntax(exit bool) {
	greenColor.Println("Syntax: codesync [-concurrency N] source-directory destination-directory")
	if runtime.GOOS == "windows" {
		whiteColor.Println(`Example: codesync C:\code\my-app \\server-name\code-share`)
	} else {
		whiteColor.Println("Example: codesync ~/code/my-app /Volumes/mount/network-share")
	}
	if exit {
		os.Exit(1)
	}
}
