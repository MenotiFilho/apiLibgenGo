package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// DownloadBook handles the download requests and redirects to the actual download URL
func DownloadBook(w http.ResponseWriter, r *http.Request) {
	md5 := r.URL.Query().Get("md5")
	if md5 == "" {
		http.Error(w, "MD5 parameter is missing", http.StatusBadRequest)
		return
	}

	// Construct the URL to the page containing the download link
	pageURL := fmt.Sprintf("http://libgen.is/get.php?md5=%s", md5)

	// Fetch the page
	resp, err := http.Get(pageURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch download page", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Parse the HTML to find the actual download link
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		http.Error(w, "Failed to parse download page", http.StatusInternalServerError)
		return
	}

	// Find the <a> tag with the text "GET" and extract the href attribute
	var downloadURL string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if s.Text() == "GET" {
			downloadURL, _ = s.Attr("href")
		}
	})

	if downloadURL == "" {
		http.Error(w, "Download link not found", http.StatusInternalServerError)
		return
	}

	// Prepend the base URL if the download link is relative
	if !strings.HasPrefix(downloadURL, "http") {
		downloadURL = fmt.Sprintf("http://libgen.is/%s", downloadURL)
	}

	// Redirect to the actual download URL
	http.Redirect(w, r, downloadURL, http.StatusFound)
}
