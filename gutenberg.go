package twittersounds

import (
	"errors"
	"golang.org/x/net/html"
	"math/rand"
	"net/http"
	"time"
)

const gutenburgRandUrl = "https://www.gutenberg.org/ebooks/search/?sort_order=random"

type Book struct {
	Title    string
	Subtitle string
	Href     string
}

func FindBook() (Book, error) {
	resp, err := http.Get(gutenburgRandUrl)
	if err != nil {
		return Book{}, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return Book{}, errors.New(resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return Book{}, err
	}

	books := parseBookHtml(doc)

	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	index := rng.Intn(len(books))
	book := books[index]

	return book, nil
}

func parseBookHtml(n *html.Node) []Book {
	var books []Book
	bookLinks := getNodesWithTagAndClass(n, "li", "booklink", nil)

	for _, book := range bookLinks {
		// TODO these three could be concurrently
		href := getAttrVal(book, "href")
		title := getTextContentByClass(book, "title")
		subtitle := getTextContentByClass(book, "subtitle")
		nextBook := Book{Href: href, Title: title, Subtitle: subtitle}
		books = append(books, nextBook)
	}

	return books
}

func getNodesWithTagAndClass(n *html.Node, tag, class string, nodes []*html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == tag {
		for _, a := range n.Attr {
			if a.Key == "class" && a.Val == class {
				nodes = append(nodes, n)
			}
		}
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		nodes = getNodesWithTagAndClass(child, tag, class, nodes)
	}

	return nodes
}

func getAttrVal(n *html.Node, attr string) (val string) {
	if n.Type == html.ElementNode {
		for _, a := range n.Attr {
			if a.Key == attr {
				val = a.Val
				return val
			}
		}
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		val = getAttrVal(child, attr)
		if val != "" {
			return val
		}
	}

	return val
}

func getTextContentByClass(n *html.Node, class string) string {
	if n.Type == html.ElementNode {
		for _, a := range n.Attr {
			if a.Key == "class" && a.Val == class {
				return getTextContent(n)
			}
		}
	}

	var result string
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		result = getTextContentByClass(child, class)
		if result != "" {
			return result
		}
	}

	return result
}

func getTextContent(n *html.Node) (textContent string) {
	if n.Type == html.TextNode {
		textContent = n.Data
		return textContent
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		textContent = getTextContent(child)
		if textContent != "" {
			return textContent
		}
	}

	return textContent
}

///
func parseBookLinks(n *html.Node, books []Book) []Book {
	if n.Type == html.ElementNode && n.Data == "li" {
		for _, a := range n.Attr {
			if a.Key == "class" && a.Val == "booklink" {
				// book := extractBookInfo(a)
				// books = append(books, book)
			}
		}
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		books = parseBookLinks(child, books)
	}

	return books
}

func extractBookInfo(n *html.Node) *Book {
	var book Book
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					book.Href = a.Val
					break
				}
			}

			break
		}
	}

	return &book
}
