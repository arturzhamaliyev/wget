package transporter

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"

	"wget/pkg/utils"
)

var (
	isBackground bool
	name         string
)

func Switcher() error {
	flag.BoolVar(&isBackground, "B", false, "background download")
	flag.StringVar(&name, "O", "tempfile", "give name to file")
	flag.Parse()

	URL := flag.Arg(0)
	arr := strings.Split(URL, "/")
	fileName := arr[len(arr)-1]

	if name != "tempfile" {
		fileName = name
	}

	switch {
	case isBackground:
		log, err := os.Create(utils.DefaultLog)
		if err != nil {
			return err
		}
		fmt.Println(`Output will be written to "wget-log"`)

		syscall.Kill(syscall.Getppid(), syscall.SIGTSTP)

		if err := Download(URL, fileName, log); err != nil {
			return err
		}

		syscall.Kill(syscall.Getppid(), syscall.SIGCONT)

	default:
		if err := Download(URL, fileName, os.Stdout); err != nil {
			return err
		}
	}

	return nil
}
