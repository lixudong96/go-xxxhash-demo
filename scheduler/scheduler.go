package scheduler

import (
	"github.com/cespare/xxhash/v2"
	"github.com/fsnotify/fsnotify"
	"os"
	"regexp"
	"strconv"
)

func ModelSchedulerModel(eventPath, eventName string, watcher *fsnotify.Watcher) {
	switch eventName {
	case "CREATE":
		createModel(eventPath, eventName, watcher)
		break
	case "REMOVE":

		break
	case "WRITE":
		writeModel(eventPath, eventName, watcher)
		break
	}
}

func createModel(eventPath, eventName string, watcher *fsnotify.Watcher) {
	ms, _ := regexp.MatchString("yaml", eventPath)
	isHash, _ := regexp.MatchString("\\.[\\d]+\\.", eventPath)
	if isHash {
		return
	}
	if ms {
		file, _ := os.ReadFile(eventPath)
		body := string(file)
		//if len(body) <= 0 {
		//	return
		//}
		fileRenameToXXHash(body, eventPath)
	} else {
		err := watcher.Add(eventPath)
		if err != nil {
			return
		}
	}
}
func writeModel(eventPath, eventName string, watcher *fsnotify.Watcher) {
	ms, _ := regexp.MatchString("yaml", eventPath)
	if ms {
		file, _ := os.ReadFile(eventPath)
		body := string(file)
		//if len(body) <= 0 {
		//	return
		//}
		updateFilenameXXHash(body, eventPath)
	}
}

func fileRenameToXXHash(body string, path string) {
	sum64String := xxhash.Sum64String(body)
	compile := regexp.MustCompile("yaml")
	allString := compile.ReplaceAllString(path, "")
	allString = allString + strconv.FormatUint(sum64String, 10) + ".yaml"
	err := os.Rename(path, allString)
	if err != nil {
		return
	}
}

func updateFilenameXXHash(body string, path string) {
	sum64String := xxhash.Sum64String(body)
	compile := regexp.MustCompile("\\.[\\d]+\\.")
	allString := compile.ReplaceAllString(path, "."+strconv.FormatUint(sum64String, 10)+".")
	err := os.Rename(path, allString)
	if err != nil {
		return
	}
}
