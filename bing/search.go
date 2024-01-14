package bing

import (
	"archive/zip"
	"bytes"
	"fmt"
	"log"
	"net/url"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/dmdhrumilmistry/gohtools/http"
	"github.com/dmdhrumilmistry/gohtools/metadata"
)

type BingSearch struct {
	Results []interface{}
}

func NewBingSearch() *BingSearch {
	return &BingSearch{}
}

func (c *BingSearch) IsPresentInResults(target interface{}) bool {
	for _, result := range c.Results {
		if result == target {
			return true
		}
	}
	return false
}

func (c *BingSearch) documentSearchHandler(i int, s *goquery.Selection) {
	url, exists := s.Attr("href")
	if !exists {
		return
	}

	if !c.IsPresentInResults(url) && http.IsValidURL(url) {
		c.Results = append(c.Results, url)
	}
}

func (c *BingSearch) ExtractExcelSheetMetaData(data []byte) (*metadata.OfficeCoreProperties, *metadata.OfficeAppProperties, error) {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		log.Printf("[BING-SEARCH-ERROR] Unable to extract excel sheet document metadata due to error: %s\n", err)
		return nil, nil, err
	}

	return metadata.NewXlsxProperties(reader)
}

func (c *BingSearch) SearchDocument(domain, filetype string) error {
	query := fmt.Sprintf("site:%s && filetype:%s && instreamset:(url title):%s\n", domain, filetype, filetype)
	url := fmt.Sprintf("https://bing.com/search?q=%s", url.QueryEscape(query))

	// scrape first page for doc links
	htmlPage := http.NewApiClient().Get(url)
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlPage))
	if err != nil {
		log.Printf("[BING-SEARCH-ERROR] Unable to create document for %s page due to error: %s\n", url, err)
		return err
	}

	// store doc links in c.Results
	selectorString := fmt.Sprintf("a[href][href$='%s']", filetype)
	doc.Find(selectorString).Each(c.documentSearchHandler)

	return nil
}

func (c *BingSearch) ExtractDocumentMetaData(docLink string) {
	docResp := http.NewApiClient().Get(docLink)
	docType := http.GetDocTypeFromLink(docLink)

	switch docType {
	case "xlsx", "xls", "xlsm", "xltx", "xltm":
		coreProps, appProps, err := c.ExtractExcelSheetMetaData(docResp)
		if err != nil {
			log.Printf("[BING-SEARCH-ERROR] Unable to extract %s document metadata due to error: %s\n", docLink, err)
		} else {
			log.Printf("%s - %s %s - %s %s\n", docLink, coreProps.Creator, coreProps.LastModifiedBy, appProps.Application, appProps.GetMajorVersion())
		}
	default:
		log.Printf("[BING-SEARCH-WARNING] Unable to create document metadata for %s page due to unknown file data type.\n", docLink)
	}
}

func (c *BingSearch) ExtractDocumentsMetaData(docLinks []string) {
	var wg sync.WaitGroup

	if len(docLinks) < 1 {
		log.Println("[BING-SEARCH-WARNING] Slice is empty!")
		return
	}

	for _, docLink := range docLinks {
		wg.Add(1)
		go func(c *BingSearch, docLink *string, wg *sync.WaitGroup) {
			c.ExtractDocumentMetaData(*docLink)
			wg.Done()
		}(c, &docLink, &wg)
	}
	wg.Wait()
}
