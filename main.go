package main

import (
	"fmt"
	// "html/template"
	// "image"
	"log"
	"net/http"
	// "os"
	"sync"
	"time"

	"gocv.io/x/gocv"
)

var (
	err      error
	webcam   *gocv.VideoCapture
	frame []byte
	mutex = &sync.Mutex{}
)

func main() {
	host := "0.0.0.0:3000"

	// open webcam
	webcam, err = gocv.VideoCaptureDevice(1)
	if err != nil {
		fmt.Printf("Error opening capture device: \n")
		return
	}
	// start capturing
	go getframes()

	fmt.Println("Capturing. Open http://" + host)

	// start http server
	http.HandleFunc("/video", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")
		data := ""
		for {
			mutex.Lock()
			data = "--frame\r\n  Content-Type: image/jpeg\r\n\r\n" + string(frame) + "\r\n\r\n"
			mutex.Unlock()
			time.Sleep(33 * time.Millisecond)
			w.Write([]byte(data))
		}
	})

	log.Fatal(http.ListenAndServe(host, nil))
}

func getframes() {
	time.Sleep(1000 * time.Millisecond)
	img := gocv.NewMat()
	defer img.Close()
	defer webcam.Close()
	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("Device closed\n")
			return
		}
		if img.Empty() {
			continue
		}
		// gocv.Resize(img, &img, image.Point{}, float64(0.5), float64(0.5), 0)
		frame, _ = gocv.IMEncode(".jpg", img)
	}
}
