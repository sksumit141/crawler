# 🕷️ GitHub Profile Crawler

A Go-based **crawler** and scraper that extracts public GitHub profile information and stores it in MongoDB. Designed to crawl and parse user profile data for building a local database useful for research, analytics, or developer directories.

---

## 🚀 Features

- 🔍 **Scrapes**: Username, bio, location, and repository count.
- 🕸️ **Crawls**: Starts from a given GitHub profile and parses its public data.
- 🧠 **Extensible**: Structure supports future recursion to discover followers/following.
- 🗃️ **Stores**: Profile data in MongoDB for easy querying.
- 🕒 **Rate Limited**: Polite crawling to comply with GitHub’s scraping policies.

---

## 📦 Requirements

- Go 1.18 or higher
- MongoDB (local or Atlas)
- GitHub profile URL to crawl

---

## ⚙️ Setup

1. **Clone the repository**:

   ```bash
   git clone https://github.com/sksumit141/crawler.git
   cd crawler
