package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var visited = make(map[string]bool)

func isValidUrl(url string) bool {
	return strings.HasPrefix(url, "https://")
}

func storeURL(client *mongo.Client, url string) {
	collection := client.Database("crawler").Collection("urls")
	_, err := collection.InsertOne(context.Background(), bson.M{"url": url})
	if err != nil {
		log.Println("Error storing Url: ", err)
	}
}

type GitHubProfile struct {
	URL          string
	Username     string
	Repositories []string
	Bio          string
	Location     string
}

func storeProfile(client *mongo.Client, profile GitHubProfile) {
	fmt.Printf("\n--- GitHub Profile Data ---\n")
	fmt.Printf("URL: %s\n", profile.URL)
	fmt.Printf("Username: %s\n", profile.Username)
	fmt.Printf("Bio: %s\n", profile.Bio)
	fmt.Printf("Location: %s\n", profile.Location)
	fmt.Printf("Repositories:\n")
	for _, repo := range profile.Repositories {
		fmt.Printf("- %s\n", repo)
	}
	fmt.Printf("------------------------\n\n")

	collection := client.Database("crawler").Collection("github_profiles")
	_, err := collection.InsertOne(context.Background(), profile)
	if err != nil {
		log.Println("Error storing profile: ", err)
	}
}

func crawl(client *mongo.Client, url string, depth int, maxDepth int) {
	if depth > maxDepth {
		return
	}
	visited[url] = true
	fmt.Println("Crawling", url)

	c := colly.NewCollector(
		colly.AllowedDomains("github.com", "www.github.com"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
		colly.MaxDepth(2),
	) // Removed Async mode

	profile := GitHubProfile{URL: url}

	c.OnHTML(".vcard-names", func(e *colly.HTMLElement) {
		profile.Username = e.ChildText(".p-nickname")
		fmt.Println("Found username:", profile.Username)
	})

	c.OnHTML(".user-profile-bio", func(e *colly.HTMLElement) {
		profile.Bio = strings.TrimSpace(e.Text)
		fmt.Println("Found bio:", profile.Bio)
	})

	c.OnHTML(".js-profile-editable-area", func(e *colly.HTMLElement) {
		profile.Location = e.ChildText(".p-label")
		fmt.Println("Found location:", profile.Location)
	})

	c.OnHTML("#user-repositories-list", func(e *colly.HTMLElement) {
		e.ForEach("h3 a", func(_ int, repo *colly.HTMLElement) {
			repoLink := repo.Attr("href")
			if repoLink != "" {
				profile.Repositories = append(profile.Repositories, "https://github.com"+repoLink)
				fmt.Println("Found repository:", repoLink)
			}
		})
	})

	// Add debug logging
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Response received from:", r.Request.URL)
	})

	c.OnScraped(func(r *colly.Response) {
		storeProfile(client, profile)
	})

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

	// Wait for collector to finish
	c.Wait()
}

func main() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://singhksumit2004:XLRI%40581@cluster0.wtzocnh.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	seedURL := "https://github.com/sksumit141"
	storeURL(client, seedURL)
	crawl(client, seedURL, 0, 2)
}
