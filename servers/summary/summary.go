package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

//PreviewImage represents a preview image for a page
type PreviewImage struct {
	URL       string `json:"url,omitempty"`
	SecureURL string `json:"secureURL,omitempty"`
	Type      string `json:"type,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Alt       string `json:"alt,omitempty"`
}

//PageSummary represents summary properties for a web page
type PageSummary struct {
	Type        string          `json:"type,omitempty"`
	URL         string          `json:"url,omitempty"`
	Title       string          `json:"title,omitempty"`
	SiteName    string          `json:"siteName,omitempty"`
	Description string          `json:"description,omitempty"`
	Author      string          `json:"author,omitempty"`
	Keywords    []string        `json:"keywords,omitempty"`
	Icon        *PreviewImage   `json:"icon,omitempty"`
	Images      []*PreviewImage `json:"images,omitempty"`
}

const headerCORS = "Access-Control-Allow-Origin"
const corsAnyOrigin = "*"

//SummaryHandler handles requests for the page summary API.
//This API expects one query string parameter named `url`,
//which should contain a URL to a web page. It responds with
//a JSON-encoded PageSummary struct containing the page summary
//meta-data.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	/*TODO: add code and additional functions to do the following:
	- Add an HTTP header to the response with the name
	 `Access-Control-Allow-Origin` and a value of `*`. This will
	  allow cross-origin AJAX requests to your server.
	- Get the `url` query string parameter value from the request.
	  If not supplied, respond with an http.StatusBadRequest error.
	- Call fetchHTML() to fetch the requested URL. See comments in that
	  function for more details.
	- Call extractSummary() to extract the page summary meta-data,
	  as directed in the assignment. See comments in that function
	  for more details
	- Close the response HTML stream so that you don't leak resources.
	- Finally, respond with a JSON-encoded version of the PageSummary
	  struct. That way the client can easily parse the JSON back into
	  an object. Remember to tell the client that the response content
	  type is JSON.
	Helpful Links:
	https://golang.org/pkg/net/http/#Request.FormValue
	https://golang.org/pkg/net/http/#Error
	https://golang.org/pkg/encoding/json/#NewEncoder
	*/
	w.Header().Set(headerCORS, corsAnyOrigin)

	myURL := r.URL.Query().Get("url")
	if len(myURL) == 0 {
		log.Printf("Bad request")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	resp, err := fetchHTML(myURL)
	if err != nil {
		log.Printf("Error fetching URL")
	}
	p, err := extractSummary(myURL, resp)
	if err != nil {
		http.Error(w, "Error extracting summary", http.StatusBadRequest)
		return
	}
	defer resp.Close()

	w.Header().Add("Content-Type", "application/json") // add json content type before encoding
	json.NewEncoder(w).Encode(p)

}

//fetchHTML fetches `pageURL` and returns the body stream or an error.
//Errors are returned if the response status code is an error (>=400),
//or if the content type indicates the URL is not an HTML page.
func fetchHTML(pageURL string) (io.ReadCloser, error) {
	/*TODO: Do an HTTP GET for the page URL. If the response status
	code is >= 400, return a nil stream and an error. If the response
	content type does not indicate that the content is a web page, return
	a nil stream and an error. Otherwise return the response body and
	no (nil) error.
	To test your implementation of this function, run the TestFetchHTML
	test in summary_test.go. You can do that directly in Visual Studio Code,
	or at the command line by running:
		go test -run TestFetchHTML
	Helpful Links:
	https://golang.org/pkg/net/http/#Get
	*/
	if !strings.HasPrefix(pageURL, "https://") && !strings.HasPrefix(pageURL, "http://") {
		pageURL = "https://" + pageURL
	}
	log.Printf("fetching URL %s", pageURL)
	resp, err := http.Get(pageURL)
	if err != nil {
		log.Printf("Get error found")
		return nil, errors.New("Get Error Found")
	}
	if resp.StatusCode >= 400 {
		//log.Fatalf("response status code %d was >= 400\n", resp.StatusCode)  don't use fatalf() it will stop the program
		log.Printf("response status code was >= 400")
		return nil, errors.New("response status code was >= 400") // create a new error
	}
	ctype := resp.Header.Get("Content-type")
	if !strings.HasPrefix(ctype, "text/html") {
		//log.Fatalf("response content type was %s not text/html\n", ctype)
		log.Printf("response content type was not html")
		return nil, errors.New("response content type was not html")
	}

	return resp.Body, nil
}

//extractSummary tokenizes the `htmlStream` and populates a PageSummary
//struct with the page's summary meta-data.
func extractSummary(pageURL string, htmlStream io.ReadCloser) (*PageSummary, error) {
	/*TODO: tokenize the `htmlStream` and extract the page summary meta-data
	according to the assignment description.
	To test your implementation of this function, run the TestExtractSummary
	test in summary_test.go. You can do that directly in Visual Studio Code,
	or at the command line by running:
		go test -run TestExtractSummary
	Helpful Links:
	https://drstearns.github.io/tutorials/tokenizing/
	http://ogp.me/
	https://developers.facebook.com/docs/reference/opengraph/
	https://golang.org/pkg/net/url/#URL.ResolveReference
	*/

	// create a tokenizer to extract html based on htmlstream
	tokenizer := html.NewTokenizer(htmlStream)
	p := PageSummary{}
	// property array used to store the tags information
	// it also stores image information
	propArray := []string{}
	propArray = append(propArray, "og:type")
	propArray = append(propArray, "og:url")
	propArray = append(propArray, "og:title")
	propArray = append(propArray, "og:site_name")
	propArray = append(propArray, "og:description")
	propArray = append(propArray, "description")
	propArray = append(propArray, "author")
	propArray = append(propArray, "keywords")
	propArray = append(propArray, "og:image")
	propArray = append(propArray, "og:image:secure_url")
	propArray = append(propArray, "og:image:type")
	propArray = append(propArray, "og:image:width")
	propArray = append(propArray, "og:image:height")
	propArray = append(propArray, "og:image:alt")
	// the image struct is used to store information
	image := []*PreviewImage{}
	for {
		// read the token one by one and store the token info in tokentype
		tokenType := tokenizer.Next()
		// sanity check for the tokentype, if error, quit loop
		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				break
			}
			log.Printf("error tokenizing HTML")
		}
		// if the token begins with <head> ... </head>
		if tokenType == html.StartTagToken || tokenType == html.EndTagToken || tokenType == html.SelfClosingTagToken {
			// read after the head tags
			token := tokenizer.Token()
			// first case: self closed meta information <meta> ... </meta>
			if "meta" == token.Data {
				// used stored info in propArray to iterate and find the property to compare
				// with the token's tag. e.x <meta property="og:url" content="..."> compares
				// with the string "og:url", if matches, then assign the value
				for _, prop := range propArray {
					// the helper function to extract the meta property from the token
					content, state := extractMetaProperty(token, prop)
					// if the token is valid and has the content we need
					if state == true {
						//fmt.Println("true: " + prop + " " + content)
						if prop == "og:type" {
							p.Type = content
						} else if prop == "og:url" {
							p.URL = content
						} else if prop == "og:title" {
							p.Title = content
						} else if prop == "og:site_name" {
							p.SiteName = content
						} else if prop == "og:description" {
							p.Description = content
						} else if prop == "description" && len(p.Description) == 0 {
							p.Description = content
						} else if prop == "author" {
							p.Author = content
						} else if prop == "keywords" {
							if len(strings.Split(content, ", ")) == 1 {
								p.Keywords = strings.Split(content, ",")
							} else {
								p.Keywords = strings.Split(content, ", ")
							}
						} else if prop == "og:image" {
							imageStruct := PreviewImage{}
							if !strings.Contains(content, "http://") || strings.HasPrefix(content, "/") {
								content = parseURL(pageURL, content)
							}
							imageStruct.URL = content
							image = append(image, &imageStruct)
						} else if strings.Contains(prop, "og:image:") {
							if prop == "og:image:secure_url" {
								image[len(image)-1].SecureURL = content
							} else if prop == "og:image:type" {
								image[len(image)-1].Type = content
							} else if prop == "og:image:width" {
								width, _ := strconv.Atoi(content)
								image[len(image)-1].Width = width
							} else if prop == "og:image:height" {
								height, _ := strconv.Atoi(content)
								image[len(image)-1].Height = height
							} else if prop == "og:image:alt" {
								image[len(image)-1].Alt = content
							}
						}
					}
				}
				// second case: self closed title tag (remaining finish...)
			} else if "title" == token.Data {
				tokenType = tokenizer.Next()
				if len(p.Title) == 0 && tokenType == html.TextToken {
					p.Title = tokenizer.Token().Data
				}
				// third case: the link self closed tag. need a new previewimage
				// struct to store info
			} else if "link" == token.Data {
				tokenType = tokenizer.Next()
				imageStruct := PreviewImage{}
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						if !strings.Contains(attr.Val, "http://") || strings.HasPrefix(attr.Val, "/") {
							imageStruct.URL = parseURL(pageURL, attr.Val)
						} else {
							imageStruct.URL = attr.Val
						}
					} else if attr.Key == "sizes" {
						size := strings.Split(attr.Val, "x")
						if len(size) > 1 {
							height, _ := strconv.Atoi(size[0])
							width, _ := strconv.Atoi(size[1])
							imageStruct.Height = height
							imageStruct.Width = width
						}
					} else if attr.Key == "type" {
						imageStruct.Type = attr.Val
					}
				}
				p.Icon = &imageStruct
			}
		}
	}
	// if the image is not empty, assign the value
	if len(image) > 0 {
		p.Images = image
	}
	return &p, nil
}

// the helper function is used to find the attribute value from the token
// we iterate through the attributes to make sure the key of attribute would
// be property (most of) and name (1 case)
func extractMetaProperty(t html.Token, property string) (content string, isValid bool) {
	for _, attr := range t.Attr {
		if (attr.Key == "property" && attr.Val == property) || (attr.Key == "name" && property == attr.Val) {
			isValid = true
		}
		if attr.Key == "content" {
			content = attr.Val
		}
	}
	return content, isValid
}

// helper function to parse the base URL and the relative URL
// return a string of a URL with the complete address to the img
func parseURL(pageURL string, relativeURL string) string {
	base, err := url.Parse(pageURL)
	checkErr(err)
	relative, err := url.Parse(relativeURL)
	checkErr(err)
	return base.ResolveReference(relative).String()
}

// helper function to log the error
func checkErr(err error) {
	if err != nil {
		log.Printf("error log")
	}
}
