# GitHub Profile Scraper

A Go-based web scraper that extracts public GitHub user profile details and stores them in MongoDB. Ideal for building a local database of developer profiles for analytics, research, or personal use.

## Features

- **Extracts**: Username, bio, location, and repositories from GitHub profiles.
- **Stores**: Data in MongoDB for easy querying and analysis.
- **Crawls**: A single GitHub profile and its public repositories.
- **No Recursion**: Does not follow external links or crawl unrelated pages.
- **Polite Scraping**: Implements rate-limiting to respect GitHub's terms of service.

## Requirements

- Go 1.18 or higher
- MongoDB instance (local or Atlas)
- GitHub profile URL to scrape

## Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/sksumit141/crawler.git
   cd crawler
   ```

2. Install dependencies:

   ```bash
   go get github.com/gocolly/colly/v2
   go get go.mongodb.org/mongo-driver/mongo
   ```

3. Update the MongoDB connection URI in the code:

   ```go
   const mongoURI = "your-mongodb-uri-here"
   ```

4. Replace the GitHub profile URL in the `main` function:

   ```go
   profileURL := "https://github.com/username"
   ```

## Usage

Run the scraper:

```bash
go run main.go
```

The profile data will be stored in the `github_profiles` collection of your MongoDB database.

## Example Output

```bash
🔍 Scraping profile: https://github.com/sksumit141
👤 Username: sksumit141
📍 Location: India
📝 Bio: Full-stack Developer
📦 Repos: 5
✅ Profile stored in MongoDB: sksumit141
```


## Acknowledgments

- [GoColly](https://github.com/gocolly/colly) for the web scraping framework.
- [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver) for MongoDB integration.
- [GitHub Profile ReadMe Maker](https://github-readme-maker.vercel.app/) for inspiration on structuring README files.
