package main

import(
	"strings"
	"golang.org/x/net/html"
)


func stripHTMLTags(input string) string {
    domDoc := html.NewTokenizer(strings.NewReader(input))
    var text strings.Builder
    
    for {
        tokenType := domDoc.Next()
        if tokenType == html.ErrorToken {
            break
        }
        if tokenType == html.TextToken {
            text.WriteString(strings.TrimSpace(string(domDoc.Text())) + " ")
        }
    }
    return strings.TrimSpace(text.String())
}