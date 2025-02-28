# Blog AggreGator 🐊

A command-line application for managing and reading RSS/Atom feeds.

## 🛠️ Prerequisites

To run this application you need:
 * **Go** (1.16 or newer) - [Download Go](https://go.dev/dl/)
 * **PostgreSQL** (12 or newer) - [Download PostgreSQL](https://www.postgresql.org/download/)

## ⚙️ Installation

Install the CLI tool using Go:
```bash
go install github.com/peeta98/blog-aggregator@latest
```

## 🔧 Configuration

Create a configuration file:
```bash
touch ~/.gatorconfig.json
```

Add your database connection details to the file:
```bash
{
  "db_url": "postgresql://username:password@localhost:5432/blog_aggregator?sslmode=disable",
  "current_user": ""
}
```

Make sure to:

1. Replace `username` and `password` with your PostgreSQL credentials
2. Create the `gator` database in PostgreSQL before using the application
3. The `current_user` field will be populated when you register or login

## 📋 Usage Examples

### 👤 User Management

```bash
# Register a new user
blog-aggregator register <username>

# Login as an existing user
blog-aggregator login <username>

# List all users (highlights the current user logged in)
blog-aggregator users
```

### 📚 Feed Management

```bash
# Add a new feed (automatically follows it)
blog-aggregator addfeed <name> <url>

# List all available feeds
blog-aggregator feeds
```

### ✅ Following Feeds

```bash
# Follow a feed
blog-aggregator follow <feed_url>

# List feeds you're following
blog-aggregator following

# Unfollow a feed
blog-aggregator unfollow <feed_url>
```

### 📖 Reading Posts

```bash
# Browse posts from feeds you follow (default: 2 posts)
blog-aggregator browse [limit]
```

### ⏱️ Feed Collection

```bash
# Start the aggregator to collect feeds at regular intervals
blog-aggregator agg <time_between_requests>
# Example: blog-aggregator agg 1m (for every minute)
```

### 🔄 Reset Database

```bash
# Reset the database (removes all users)
blog-aggregator reset
```

## 📝 Notes

 * You must register and login before using most commands
 * Feed URLs must use HTTP or HTTPS protocols
 * Time durations should be specified using Go's duration format (e.g., "30s", "1m", "1h")
 * RSS feeds will be automatically collected at intervals you specify with the `agg` command