package main

import (
	"flag"
	"github.com/howeyc/fsnotify"
	"log"
	"os/exec"
	"regexp"
)

var (
	listenDir string
	extReg    *regexp.Regexp
)

func init() {
	flag.StringVar(&listenDir, "d", "", "help message for flagname")
	flag.Parse()
	extReg = regexp.MustCompile("(.*?)\\.scss$")
}
func main() {
	log.Println("Listening for SASS changes in", listenDir)
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				m, err := regexp.MatchString("\\.scss$", ev.Name)
				if err != nil {
					log.Fatal(err)
				}
				if m {
					log.Print("mask", ev)
					destFile := extReg.ReplaceAllString(ev.Name, "$1.css")
					log.Println("Processing", ev.Name, "into", destFile)
					cmd := exec.Command("sass", ev.Name, destFile)
					stdout, err := cmd.StdoutPipe()
					if err != nil {
						log.Fatal(err)
					}
					if err := cmd.Start(); err != nil {
						log.Fatal(err)
					}
					log.Println(stdout)
				}
			case err := <-watcher.Error:
				log.Println("Error:", err)
			}
		}
	}()

	// err = watcher.Watch(listenDir)
	watcher.WatchFlags(listenDir, fsnotify.FSN_MODIFY)
	if err != nil {
		log.Fatal(err)
	}

	<-done

	watcher.Close()
}
