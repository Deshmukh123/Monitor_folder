package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

func WatchFolders(folders []string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %v", err)
	}
	defer watcher.Close()

	for _, folder := range folders {
		err = watcher.Add(folder)
		if err != nil {
			log.Printf("Error watching folder %s: %v", folder, err)
			continue
		}
		log.Printf("Started watching folder: %s", folder)
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Printf("New file detected: %s", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("Error: %v", err)
			}
		}
	}()
	<-done
	return nil
}

func main() {
	folders := []string{"C:/Users/vaish/go/src/service/new-service/clients", "C:/Users/vaish/go/src/service/new-service/server"} // Specify folders here

	for _, folder := range folders {
		if _, err := os.Stat(folder); os.IsNotExist(err) {
			err := os.Mkdir(folder, os.ModePerm)
			if err != nil {
				log.Fatalf("Failed to create folder %s: %v", folder, err)
			}
		}
	}

	if err := WatchFolders(folders); err != nil {
		log.Fatalf("Error starting folder watch: %v", err)
	}
}
