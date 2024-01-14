package metadata

import (
	"archive/zip"
	"encoding/xml"
	"log"
	"strings"
)

type OfficeCoreProperties struct {
	XmlName        string `xml:"coreProperties"`
	Creator        string `xml:"creator"`
	Title          string `xml:"title"`
	LastModifiedBy string `xml:"lastModifiedBy"`
}

type OfficeAppProperties struct {
	XmlName     string `xml:"Properties"`
	Application string `xml:"Application"`
	Company     string `xml:"Company"`
	AppVersion  string `xml:"AppVersion"`
}

func (c *OfficeAppProperties) GetMajorVersion() string {
	var OfficeVersions = map[string]string{
		"16": "2016",
		"15": "2013",
		"14": "2010",
		"12": "2007",
		"11": "2003",
	}
	tokens := strings.Split(c.AppVersion, ".")

	if len(tokens) < 2 {
		return "unknown"
	}

	version, ok := OfficeVersions[tokens[0]]
	if !ok {
		return "unknown"
	}

	return version
}

func ProcessZipFile(file *zip.File, prop interface{}) error {
	readCloser, err := file.Open()
	if err != nil {
		log.Printf("[METADATA-XLSX-ERROR] Unable to read file %s due to error: %s\n", file.Name, err)
		return err
	}
	defer readCloser.Close()

	if err := xml.NewDecoder(readCloser).Decode(&prop); err != nil {
		log.Printf("[METADATA-XLSX-ERROR] Unable to decode file %s due to error: %s\n", file.Name, err)
		return err
	}
	return nil
}

func NewXlsxProperties(z *zip.Reader) (*OfficeCoreProperties, *OfficeAppProperties, error) {
	var coreProps OfficeCoreProperties
	var appProps OfficeAppProperties

	for _, file := range z.File {
		switch file.Name {
		case "docProps/core.xml":
			if err := ProcessZipFile(file, &coreProps); err != nil {
				log.Printf("[METADATA-XLSX-ERROR] Unable to process %s file due to error: %s\n", file.Name, err)
				return nil, nil, err
			}

		case "docProps/app.xml":
			if err := ProcessZipFile(file, &appProps); err != nil {
				log.Printf("[METADATA-XLSX-ERROR] Unable to process %s file due to error: %s\n", file.Name, err)
				return nil, nil, err
			}
		}
	}

	return &coreProps, &appProps, nil
}
