package download

type Downloader interface {
	Download(url, filename string) error
}
