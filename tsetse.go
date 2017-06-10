package main

import (
	"os/exec"
	"runtime"
	"net/http"
	"log"
	"time"
	"fmt"
	"github.com/tealeg/xlsx"
	"sync"
)

func main() {

	start := time.Now()

	excelFileName := "test.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)

	if err != nil {
		fmt.Println(err)
	}

	var wg sync.WaitGroup

	fmt.Println("Making calls")

	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {

			wg.Add(1)

			name, _ := row.Cells[0].String()
			url, _ := row.Cells[2].String()

			fmt.Print(".")
			go makeHttpCall(name, url, &wg)
		}
	}

	fmt.Println("Waiting..")
	wg.Wait()

	elapsed := time.Since(start)
	log.Printf("Binomial took %s", elapsed)

	time.Sleep(60 * time.Second)
	fmt.Println("This is working")
}


func makeHttpCall(name string, urll string, wg *sync.WaitGroup){

	baseUrl := "http://10.252.169.12:7788" + urll
	resp, err := http.Get(baseUrl)

	if resp != nil {

		if resp.StatusCode == 200 {

			fmt.Print("+")
		} else {

			fmt.Println(baseUrl, resp.StatusCode, err)
		}
	} else {

		fmt.Println("Didn't make call...", baseUrl)
	}

	wg.Done()
}


func startWS(){
	fs := http.FileServer(http.Dir("static-pages"))
	http.Handle("/", fs)

	log.Println("Listening...")

	go openBrowser("http://127.0.0.1:3000/index.html")

	http.ListenAndServe(":3000", nil)
}




// openBrowser tries to open the URL in a browser,
// and returns whether it succeed in doing so.
func openBrowser(url string) bool {

	time.Sleep(3 * time.Second)

	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}