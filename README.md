# Gator - RSS Feed Reader CLI

Gator is a command-line RSS feed reader that helps you manage and read content from your favorite websites.

## Prerequisites

Before you can use Gator, you'll need:

- **Go** (version 1.20 or later)
- **PostgreSQL** database

## Installation

Install Gator using Go:

```bash
go install github.com/sianwa11/gator@latest
```

## Configuration

1. Create a configuration directory and file:

```bash
mkdir -p ~/.config/gator
touch ~/.config/gator/config.json
```

2. Edit the config file with your database connection details:

```json
{
  "database_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

Replace `username` and `password` with your PostgreSQL credentials.

## Usage

Once installed, run Gator by typing:

```bash
gator
```

### Available Commands

#### User Management
- `register <username>` - Create a new user
- `login <username>` - Log in as a user

#### Feed Management
- `add-feed <name> <url>` - Add a new RSS feed
- `list-feeds` - Show all your feeds
- `fetch-feeds` - Update all your feeds

#### Reading Content
- `browse [limit]` - Browse recent posts (default limit is 10)

#### Help
- `help` - Display available commands

## Example Workflow

```
$ gator
> register john_doe
User created successfully!

> login john_doe
Logged in as john_doe

> add-feed "Tech News" "https://news.ycombinator.com/rss"
Feed added successfully!

> browse 5
Showing 5 most recent posts...
```

## License