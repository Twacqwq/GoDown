package download

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/Twacqwq/godown/utils"
	"github.com/schollz/progressbar/v3"
)

var progress *progressbar.ProgressBar

type GoDown struct {
	concurrency int
	path        string
}

func NewGodown(n int, path string) *GoDown {
	if !strings.HasSuffix(path, "/") && path != "" {
		path += "/"
	}
	return &GoDown{
		concurrency: n,
		path:        path,
	}
}

func (g *GoDown) Download(url string) error {
	filename := path.Base(url)
	res, err := http.Head(url)
	if err != nil {
		return err
	}
	progress = progressbar.DefaultBytes(
		res.ContentLength,
		"downloading",
	)
	if res.StatusCode == http.StatusOK && res.Header.Get("Accept-Ranges") == "bytes" {
		return g.multiDownload(url, filename, int(res.ContentLength))
	}

	return g.singleDownload(url, filename)
}

func (g *GoDown) multiDownload(url, filename string, contentLength int) error {
	partSize := contentLength / g.concurrency
	partDir := g.path + utils.GetPartDir(filename)
	os.Mkdir(partDir, 0777)
	defer os.RemoveAll(partDir)

	var wg sync.WaitGroup
	wg.Add(g.concurrency)

	sliceStart := 0
	for i := 0; i < g.concurrency; i++ {
		go func(partNum, sliceStart int) {
			defer wg.Done()

			sliceEnd := sliceStart + partSize
			if partNum == g.concurrency-1 {
				sliceEnd = contentLength
			}
			g.multiple(url, filename, sliceStart, sliceEnd, partNum)
		}(i, sliceStart)

		sliceStart += partSize + 1
	}
	wg.Wait()
	if err := g.merge(filename); err != nil {
		log.Fatal(err)
	}

	return nil
}

func (g *GoDown) singleDownload(url, filename string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	rawFile, err := os.OpenFile(g.path+filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(io.MultiWriter(rawFile, progress), res.Body)
	if err != nil {
		if err == io.EOF {
			return nil
		}
		log.Fatal(err)
	}

	return nil
}

func (g *GoDown) multiple(url, filename string, sliceStart, sliceEnd, part int) {
	if sliceStart > sliceEnd {
		return
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", sliceStart, sliceEnd))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	partFile, err := os.OpenFile(g.path+utils.GetPartFilename(filename, part), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, 32*1024)
	_, err = io.CopyBuffer(io.MultiWriter(partFile, progress), res.Body, buf)
	if err != nil {
		if err == io.EOF {
			return
		}
		log.Fatal(err)
	}
}

func (g *GoDown) merge(filename string) error {
	rawFile, err := os.OpenFile(g.path+filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer rawFile.Close()

	for i := 0; i < g.concurrency; i++ {
		partFilename := g.path + utils.GetPartFilename(filename, i)
		partFile, err := os.Open(partFilename)
		if err != nil {
			return err
		}
		io.Copy(rawFile, partFile)
		partFile.Close()
		os.Remove(partFilename)
	}

	return nil
}
