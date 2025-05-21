package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoURI = //your mongodb uri

type GitHubProfile struct { 
	URL          string   `bson:"url"`
	Username     string   `bson:"username"`
	Bio          string   `bson:"bio"`
	Location     string   `bson:"location"`
	Repositories []string `bson:"repositories"`
}

func storeProfile(client *mongo.Client, profile GitHubProfile) {
	collection := client.Database("crawler").Collection("github_profiles")
	_, err := collection.InsertOne(context.Background(), profile)
	if err != nil {
		log.Println("Error storing profile:", err)
	} else {
		fmt.Println("✅ Profile stored in MongoDB:", profile.Username)
	}
}

var visitedProfiles = make(map[string]bool)

func crawlProfile(client *mongo.Client, url string, depth int, maxDepth int) {
	if depth > maxDepth || visitedProfiles[url] {
		return
	}
	visitedProfiles[url] = true
	fmt.Printf("🔍 Scraping profile [%d/%d]: %s\n", depth, maxDepth, url)

	profile := GitHubProfile{URL: url}

	c := colly.NewCollector(
		colly.AllowedDomains("github.com", "www.github.com"),
		colly.UserAgent("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*github.*",
		Parallelism: 1,
		Delay:       1 * time.Second,
	})

	c.OnHTML(".vcard-names", func(e *colly.HTMLElement) {
		profile.Username = e.ChildText(".p-nickname")
	})

	c.OnHTML(".user-profile-bio", func(e *colly.HTMLElement) {
		profile.Bio = strings.TrimSpace(e.Text)
	})

	c.OnHTML(".js-profile-editable-area", func(e *colly.HTMLElement) {
		profile.Location = e.ChildText(".p-label")
	})

	c.OnHTML("#user-repositories-list h3 a", func(e *colly.HTMLElement) {
		repoPath := e.Attr("href")
		if repoPath != "" {
			repoURL := "https://github.com" + repoPath
			profile.Repositories = append(profile.Repositories, repoURL)
			fmt.Printf("Found repository: %s\n", repoURL)
		}
	})

	
	c.OnHTML("a[data-hovercard-type='user']", func(e *colly.HTMLElement) {
		userLink := e.Attr("href")
		if userLink != "" && !visitedProfiles["https://github.com"+userLink] {
			go crawlProfile(client, "https://github.com"+userLink, depth+1, maxDepth)
		}
	})

	
	c.OnResponse(func(r *colly.Response) {
		fmt.Printf("Status Code: %d for %s\n", r.StatusCode, r.Request.URL)
	})

	c.OnScraped(func(_ *colly.Response) {
		fmt.Printf("👤 Username: %s\n📍 Location: %s\n📝 Bio: %s\n📦 Repos: %d\n", profile.Username, profile.Location, profile.Bio, len(profile.Repositories))
		storeProfile(client, profile)
	})

	if err := c.Visit(url); err != nil {
		log.Printf("Error visiting %s: %v\n", url, err)
	}
}

func main() {
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	maxDepth := 3
	seedProfiles := []string{
		"https://github.com/sksumit141", //add more profiles
		
	}

	
	for {
		visitedProfiles = make(map[string]bool) 

		fmt.Println("\n🔄 Starting new crawl cycle...")
		for _, profile := range seedProfiles {
			crawlProfile(client, profile, 0, maxDepth)
		}

		fmt.Println("\n😴 Waiting before next cycle...")
		time.Sleep(5 * time.Minute) // Wait 5 minutes between cycles
	}
}
