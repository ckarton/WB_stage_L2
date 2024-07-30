package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: wget <url>")
		return
	}

	startURL := os.Args[1]

	// Создаем папку для хранения скачанных файлов
	err := os.MkdirAll("downloaded_site", os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	// Скачиваем и сохраняем главную страницу
	err = downloadPage(startURL, "downloaded_site/index.html")
	if err != nil {
		fmt.Println("Error downloading page:", err)
		return
	}

	fmt.Println("Site downloaded successfully!")
}

func downloadPage(pageURL, outputPath string) error {
	resp, err := http.Get(pageURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Сохраняем HTML страницы
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	// Парсим HTML для поиска ссылок на ресурсы
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	doc, err := html.Parse(f)
	if err != nil {
		return err
	}

	baseURL, err := url.Parse(pageURL)
	if err != nil {
		return err
	}

	// Рекурсивно скачиваем ресурсы
	var downloadResources func(*html.Node)
	downloadResources = func(n *html.Node) {
		if n.Type == html.ElementNode {
			var resourceURL string
			switch n.Data {
			case "img", "script":
				for _, attr := range n.Attr {
					if attr.Key == "src" {
						resourceURL = attr.Val
					}
				}
			case "link":
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						resourceURL = attr.Val
					}
				}
			}

			if resourceURL != "" {
				absURL := toAbsURL(baseURL, resourceURL)
				localPath := toLocalPath(absURL)
				err := downloadPage(absURL, localPath)
				if err != nil {
					fmt.Println("Error downloading resource:", err)
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			downloadResources(c)
		}
	}
	downloadResources(doc)

	return nil
}

func toAbsURL(base *url.URL, href string) string {
	u, err := url.Parse(href)
	if err != nil {
		return href
	}
	return base.ResolveReference(u).String()
}

func toLocalPath(resourceURL string) string {
	u, err := url.Parse(resourceURL)
	if err != nil {
		return resourceURL
	}
	path := filepath.Join("downloaded_site", filepath.FromSlash(u.Path))
	if strings.HasSuffix(u.Path, "/") {
		path = filepath.Join(path, "index.html")
	}
	return path
}
