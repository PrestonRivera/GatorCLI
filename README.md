# Gator CLI

- A multi-user CLI tool for aggregating RSS feeds and veiwing the posts.

## Motivation

I had a problem with keeping up with my favorite dev blogs. Having to switch between 
multiple browser tabs and manage notifications from different blog subscriptions was eating up 20-30 minutes of my day.

As a developer who practically lives in the terminal, GatorCLI brings all my favorite blogs to me. 
I'm able to quickly check blogs while coding without context-switching. I made it multi-user capable 
for people who work on a shared development server. A whole development team can use it independently with their own feed configurations.

## Prerequisites

- Go version: go1.23.5
- PostgreSQL version: 15.10

# Quick Start

## Installation

1. Install `gator` with:
  ```
  go install github.com/PrestonRivera/GatorCLI@latest
  ```

2. Install PostgreSQL and set up a database:
  - For Linux (Debian/Ubuntu):
    ```
    sudo apt update
    sudo apt install postgresql postgresql-contrib
    ```
  - For other systems, download from [PostgreSQL official website](https://www.postgresql.org/download/)

## Configuration

- Create a `.gatorconfig.json` file in your home directory with the following structure:
  ```json
  {
    "db_url": "postgres://username:@localhost:5432/database?sslmode=disable"
  }
  ```

- You can test your connection string by trying to connect to psql:
  ```
  psql "your_connection_string"

## Quick Start commands

- `gator register <name>`
- `gator addfeed <url>`
- `gator agg 30s`
- `gator browse 5`

## Contributing

### Clone the repo 

- `git clone https://github.com/PrestonRivera/GatorCLI.git`

### Build the project

- `go build -o gator`

### Run the project

- `./gator help`

### Submit a pull request

If you'd like to contribute, please fork the repository and open a pull request to the `main` branch.

## Available Commands

- Setup & Database:
  - `gator help` - (List available commands and the use cases)
  - `gator register <name>` - (Create a new user)
  - `gator reset` - (Reset/clear the database)
  - `gator login <name>` - (Log in as a user that already exists)

- Feed management:
  - `gator addfeed <url>` - (Add a feed to the database)
  - `gator feeds` - (List all feeds)
  - `gator follow <feed_id>` - (Follow a feed that already exists in the database)
  - `gator following` - (Lists feeds the user is following)
  - `gator unfollow <feed_id>` - (Unfollow a feed that already exists in the database)

- Content & Updates:
  - `gator browse [limit]` - (View the posts, defaults to 2)
  - `gator users` - (List all users)
  - `gator agg <30s>` - (Start the aggregator. Can use preferred time intervals such as 1m 3m 6m etc...) ctrl + c to end