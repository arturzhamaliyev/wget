package transporter

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"wget/pkg/utils"
)

var isBackground bool

func Switcher() error {
	flag.BoolVar(&isBackground, "B", false, "background download")
	flag.Parse()

	URL := flag.Arg(0)

	switch {
	case isBackground:
		log, err := os.Create(utils.DefaultLog)
		if err != nil {
			return err
		}
		fmt.Println(`Output will be written to "wget-log"`)

		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGURG)

		go func() {
			<-sigc
			// fmt.Print("\033[H\033[2J")

			// fmt.Println(1)
		}()

		syscall.Kill(syscall.Getppid(), syscall.SIGTSTP)

		if err := Download(URL, log); err != nil {
			return err
		}

		syscall.Kill(syscall.Getppid(), syscall.SIGCONT)

	default:
		if err := Download(URL, os.Stdout); err != nil {
			return err
		}
	}

	return nil
}
