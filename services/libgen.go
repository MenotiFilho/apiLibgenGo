package services

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const LibgenSearchUrl = "https://libgen.is/search.php"
const LIBGEN_DOWNLOAD_URL = "http://libgen.is/get.php?md5="

type Book struct {
	Title    string `json:"title"`
	Author   string `json:"author"`
	Year     string `json:"year"`
	Language string `json:"language"`
	MD5      string `json:"md5"`
	ISBN     string `json:"isbn"`
	FileType string `json:"file_type"`
}

type SearchType string

const (
	SearchByTitle  SearchType = "title"
	SearchByAuthor SearchType = "author"
	SearchByISBN   SearchType = "isbn"
)

// SearchBooks busca livros no Libgen baseado na query e tipo de busca
func SearchBooks(query string, searchType SearchType) ([]Book, error) {
	column := "def"
	switch searchType {
	case SearchByTitle:
		column = "title"
	case SearchByAuthor:
		column = "author"
	case SearchByISBN:
		column = "identifier"
	}

	searchURL := fmt.Sprintf("%s?req=%s&res=25&view=simple&phrase=1&column=%s", LibgenSearchUrl, query, column)
	resp, err := http.Get(searchURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search failed with status code %d", resp.StatusCode)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// log error if needed
		}
	}(resp.Body)

	return parseSearchResults(resp)
}

// parseSearchResults analisa os resultados da pesquisa HTML em um slice de Book
func parseSearchResults(resp *http.Response) ([]Book, error) {
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// Enhanced regex to find ISBN patterns (ISBN-10 and ISBN-13 with optional hyphens or spaces)
	isbnRegex := regexp.MustCompile(`\b(?:ISBN(?:-1[03])?:?\s*)?(?:97[89][-\s]?)?\d{1,5}[-\s]?\d{1,7}[-\s]?\d{1,6}[-\s]?(?:\d|X)\b`)

	var books []Book
	doc.Find("table.c tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			// Pular a linha do cabeçalho
			return
		}
		cols := s.Find("td")
		if cols.Length() < 11 {
			return
		}

		// Extrair título e MD5 da URL do título
		titleLink := cols.Eq(2).Find("a[href]")
		fullTitle := titleLink.Text()
		href, exists := titleLink.Attr("href")
		md5 := ""
		if exists {
			md5Parts := strings.Split(href, "=")
			if len(md5Parts) > 1 {
				md5 = md5Parts[1]
			}
		}

		// Extrair ISBNs
		isbnMatches := isbnRegex.FindAllString(fullTitle, -1)
		isbn := strings.Join(isbnMatches, ", ")
		// Remover ISBNs do título
		title := isbnRegex.ReplaceAllString(fullTitle, "")
		title = strings.TrimSpace(title)
		// Remover vírgulas e espaços extras do título
		title = strings.TrimRight(title, ", ")
		title = strings.TrimSpace(title)

		// Extrair tipo de arquivo
		fileType := strings.TrimSpace(cols.Eq(8).Text())

		book := Book{
			Title:    title,
			Author:   strings.TrimSpace(cols.Eq(1).Text()),
			Year:     strings.TrimSpace(cols.Eq(4).Text()),
			Language: strings.TrimSpace(cols.Eq(6).Text()),
			MD5:      md5,
			ISBN:     isbn,
			FileType: fileType,
		}
		books = append(books, book)
	})

	return books, nil
}

// DownloadBook baixa o livro pelo MD5
func DownloadBook(md5 string) (*http.Response, error) {
	downloadURL := fmt.Sprintf("%s%s", LIBGEN_DOWNLOAD_URL, md5)
	resp, err := http.Get(downloadURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download book")
	}
	return resp, nil
}
