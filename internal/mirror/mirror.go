package mirror

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"sync"

	"github.com/kurochkinivan/web-local-mirror/internal/crawl"
	"github.com/kurochkinivan/web-local-mirror/internal/fsm"
)

type Mirror struct {
	baseURL     *url.URL
	workerCount int
}

func NewMirror(baseURL string, workerCount int) (*Mirror, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL, err: %w", err)
	}

	return &Mirror{
		baseURL:     parsedURL,
		workerCount: workerCount,
	}, nil
}

func (m *Mirror) Start() error {
	hostname := m.baseURL.Hostname()
	err := fsm.InitializeDirectory(hostname)
	if err != nil {
		return fmt.Errorf("failed to create directory, err: %w", err)
	}

	urlChan := make(chan string, m.workerCount)
	errChan := make(chan error, m.workerCount)
	doneChan := make(chan struct{})

	var wg sync.WaitGroup
	for i := 0; i < m.workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.worker(urlChan, errChan)
		}()
	}

	urlChan <- m.baseURL.String()
	go crawl.CrawlURL(urlChan, errChan, m.baseURL.String())

	go func() {
		wg.Wait()
		doneChan <- struct{}{}
	}()

	for {
		select {
		case err := <-errChan:
			if err != nil {
				return err
			}
		case <-doneChan:
			return nil
		}
	}
}

func (m *Mirror) worker(urlChan <-chan string, errChan chan<- error) {
	for url := range urlChan {
		err := m.processSinglePage(url)
		if err != nil {
			errChan <- err
			continue
		}
	}
}

func (m *Mirror) processSinglePage(link string) error {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return fmt.Errorf("failed to parse URL %s, err: %w", link, err)
	}

	resp, err := http.Get(link)
	if err != nil {
		return fmt.Errorf("failed to get html page for url %s, err: %w", link, err)
	}

	hostname := m.baseURL.Hostname()
	path := filepath.Join(hostname, parsedURL.Path)

	err = fsm.SaveHTMLFile(path, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save HTML file, err: %w", err)
	}

	return nil
}
