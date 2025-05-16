import (
	"fmt"
	"context"
	"strings"
	"log"

	"github.com/gocolly/colly/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var visited = make(map[string]bool)

func isValidUrl(url string) bool {
	return strings.HashPrefix(url, "https://")
}

func storeURL(client *mongo.Client, url String) {
	collection := clent.Database("crawler").Collection("urls")
	_.err := collection.InsertOne(context.Background(), bson.M{"url": url})
	if err != nil {
		log.Println("Error storing Url: ", err)
	}
}

func crawl(client *mongo.Client, url string, depth int, maxDepth int) {
	if depth > maxDepth {
		return
	}
	visited[url] = true
	fmt.Println("Crawling", url)

	c := colly.NewCollector(
		colly.AllowedDomains("example.com"),
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if isValidUrl(link) && !visited[link] {
			storeURL(client, link)
			crawl(client, link, depth+1, maxDepth)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	err := c.Visit(url)
	if err != nil {
		log.Println("Error visiting url: ", err)
	}
}