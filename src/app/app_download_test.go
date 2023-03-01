package app

import (
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/inhies/go-bytesize"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

// asyncRetryDownload
// This is called only when file not-exist or exist but partial.
func asyncRetryDownload(readFileFD *os.File, installerCfg InstallerConfig, fullFileLength int64, wg *sync.WaitGroup) {
	defer wg.Done() // decrement dynamically

	// Check file existence
	retryFileInfo, err := readFileFD.Stat()
	if err != nil {
		log.Fatalln("Failed to open test file : ", installerCfg.Name)
	}

	// Check file download complete
	// This is for protecting code for recursive function
	if fullFileLength == retryFileInfo.Size() {
		readFileFD.Close()
		// only this printed when download complete without checking file size
		log.Println("Success to download complete: ", installerCfg.Name, "after retying")
		return
	}

	log.Println("Retrying file size : ", bytesize.New(float64(retryFileInfo.Size())))

	streamFile, err := os.OpenFile(readFileFD.Name(), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("Failed to open readFileFD : ", err)
	}
	defer streamFile.Close()

	res, err := req.R().
		SetHeader("Range", fmt.Sprintf("bytes=%d-", retryFileInfo.Size())).
		Get(DockerInstaller.Url)

	defer res.Body.Close()

	written, copyErr := io.Copy(streamFile, res.Body)

	log.Println("Written : ", bytesize.New(float64(written)))

	if err != nil {
		log.Println("Http request fail : ", err)
		time.Sleep(2 * time.Second)
		wg.Add(1)
		go asyncRetryDownload(readFileFD, installerCfg, fullFileLength, wg)
		return
	}

	if copyErr != nil {
		log.Println("Error occurs appending stream from response .. ", err)
		time.Sleep(2 * time.Second)
		wg.Add(1)
		go asyncRetryDownload(readFileFD, installerCfg, fullFileLength, wg)
		return
	}

	// only this printed when download complete without checking file size
	log.Println("Success to download complete: ", installerCfg.Name, "after retying")
}

func downloadInstaller(installerCfg InstallerConfig, wg *sync.WaitGroup) {
	defer wg.Done()

	// 헤더로 총 파일 크기 검사
	headerRes, err := req.R().Head(installerCfg.Url)
	if err != nil {
		log.Println("헤더 응답 오류 ", err)
	}

	// Convert string to int64
	contentLength := headerRes.GetHeader("Content-Length")
	fullFileLength, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		log.Fatal("Failed to parsing content Length : ", err)
	}

	// 이미 파일이 있는 경우
	var targetFile = installerCfg.Name + installerCfg.Ext
	if _, err := os.Stat(targetFile); err == nil {
		readFileFD, err := os.Open(targetFile)
		// 여기서 바로 닫아버리면 다음 재 다운로드 시도에서 읽다가 도중에 끊김..
		// defer readFileFD.Close()
		if err != nil {
			log.Println("Failed to open readFileFD : ", readFileFD)
			log.Fatalln("error : ", err)
		}

		fileInfo, err := readFileFD.Stat()
		if err != nil {
			log.Fatalln("Failed to get stat from File Descriptor ", err)
		}

		// 이미 완전한 파일이 존재.
		if fullFileLength == fileInfo.Size() {
			readFileFD.Close()
			log.Println("Already installed tool : ", installerCfg.Name)
			return
		}

		wg.Add(1)
		go asyncRetryDownload(readFileFD, installerCfg, fullFileLength, wg)
		return
	} else {
		// 최초 요청 케이스
		res, err := req.R().
			SetOutputFile(installerCfg.Name + installerCfg.Ext).
			Get(installerCfg.Url)

		if err != nil {
			log.Println("Failed to download : ", err)
			readFileFD, err := os.Open(targetFile)
			defer readFileFD.Close()
			if err != nil {
				log.Println("Failed to open readFileFD : ", readFileFD)
				log.Println("error : ", err)
			}

			wg.Add(1)
			go asyncRetryDownload(readFileFD, installerCfg, fullFileLength, wg)
			return
		}

		if res.GetStatusCode() != http.StatusOK {
			log.Printf("Status code : %d\n", res.GetStatusCode())
			log.Fatalln("Error", installerCfg.Name+" download failed")
			// 요청 응답 코드가 다른 경우 그냥 중단.
		}

		// 재시도 없이 한번에 성공한경우
		log.Println("Success at once ", installerCfg.Name+" download complete.")
	}
}

func Test_All_App(t *testing.T) {
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(4)

	go downloadInstaller(DockerInstaller, waitGroup)
	go downloadInstaller(NotionInstaller, waitGroup)
	go downloadInstaller(NodeInstaller, waitGroup)
	go downloadInstaller(GoInstaller, waitGroup)

	waitGroup.Wait()
}
