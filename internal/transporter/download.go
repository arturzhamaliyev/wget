package transporter

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/mxk/go-flowrate/flowrate"
)

func Download(credentials *Credentials) error {
	credentials.Mutex.Lock()

	fmt.Fprintf(credentials.OutPut, "start at %v\n", time.Now().Format("2006-01-02 15:04:05"))

	fmt.Fprint(credentials.OutPut, "sending request, awaiting response... status ")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			MaxVersion: tls.VersionTLS13,
		},
	}

	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", credentials.URL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "Chuhuahua.akl/pidor")

	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// fmt.Println(response.Header)

	if response.Status != "200 OK" {
		return errors.New(response.Status)
	}

	fmt.Fprintln(credentials.OutPut, response.Status)

	fmt.Fprintf(credentials.OutPut, "content size: %d [~%.2fMB]\n", response.ContentLength, float64(response.ContentLength)/1000000)

	fmt.Fprintf(credentials.OutPut, "saving file to: %s\n", credentials.Path+credentials.FileName)

	// TODO: duplicate file
	// if _, err := os.Stat(credentials.Path + credentials.FileName); err == nil {

	// } else if errors.Is(err, os.ErrNotExist) {

	// } else {

	// }

	file, err := os.Create(credentials.Path + credentials.FileName)
	if err != nil {
		return err
	}
	defer file.Close()

	if credentials.RateLimit != 0 {
		response.Body = flowrate.NewReader(response.Body, credentials.RateLimit)
	}

	credentials.Mutex.Unlock()

	if credentials.OutPut == os.Stdout && !credentials.IsInDir {
		template := `{{ counters .}} {{ bar . "[" "=" (cycle . "BRUH" ) "." "]"}} {{percent .}} {{speed .}} {{rtime .}}`

		bar := pb.ProgressBarTemplate(template).Start64(response.ContentLength)

		barReader := bar.NewProxyReader(response.Body)

		_, err = io.Copy(file, barReader)
		if err != nil {
			return err
		}

		bar.Finish()

	} else {
		_, err = io.Copy(file, response.Body)
		if err != nil {
			return err
		}
	}

	fmt.Fprintf(credentials.OutPut, "Downloaded [%s]\n", credentials.URL)
	fmt.Fprintf(credentials.OutPut, "finished at %v\n", time.Now().Format("2006-01-02 15:04:05"))

	return nil
}
