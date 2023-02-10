package transporter

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"

	"wget/pkg/utils"
)

type Credentials struct {
	URL          string
	FileName     string
	Path         string
	IsBackground bool
	OutPut       *os.File
}

func NewCredentialsConstructor(URL string) *Credentials {
	return &Credentials{
		URL:          URL,
		FileName:     getFileName(URL),
		Path:         "./",
		IsBackground: false,
		OutPut:       os.Stdout,
	}
}

func getFileName(URL string) string {
	arr := strings.Split(URL, "/")
	return arr[len(arr)-1]
}

var (
	backgroundFlagVal bool
	nameFlagVal       string
	pathFlagVal       string
)

func Switcher() error {
	flag.BoolVar(&backgroundFlagVal, "B", false, "background download")
	flag.StringVar(&nameFlagVal, "O", "tempfile", "give name to saved file")
	flag.StringVar(&pathFlagVal, "P", "./", "path to where you want to save the file")
	flag.Parse()

	credentials := NewCredentialsConstructor(flag.Arg(0))

	if nameFlagVal != "tempfile" {
		credentials.FileName = nameFlagVal
	}

	if pathFlagVal != "./" {
		credentials.Path = pathFlagVal
	}

	if backgroundFlagVal {
		logFile, err := os.Create(utils.DefaultLog)
		if err != nil {
			return err
		}
		fmt.Println(`Output will be written to "wget-log"`)
		credentials.OutPut = logFile

		parentProc := syscall.Getppid()
		syscall.Kill(parentProc, syscall.SIGTSTP)
		defer syscall.Kill(parentProc, syscall.SIGCONT)
	}

	if err := Download(credentials); err != nil {
		return err
	}

	return nil
}
