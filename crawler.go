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

const mongoURI = "mongodb+srv://singhksumit2004:XLRI%40581@cluster0.wtzocnh.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

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

func crawlProfile(client *mongo.Client, url string) {
	fmt.Println("🔍 Scraping profile:", url)

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
		}
	})

	c.OnScraped(func(_ *colly.Response) {
		fmt.Printf("👤 Username: %s\n📍 Location: %s\n📝 Bio: %s\n📦 Repos: %d\n", profile.Username, profile.Location, profile.Bio, len(profile.Repositories))
		storeProfile(client, profile)
	})

	if err := c.Visit(url); err != nil {
		log.Println("Error visiting URL:", err)
	}
}

func main() {
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	profileURL := "https://github.com/sksumit141"

	crawlProfile(client, profileURL)
}
