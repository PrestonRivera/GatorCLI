# Gator CLI

- A multi-player CLI tool for aggregating RSS feeds and veiwing the posts.

## Prerequisites

- Go version: go1.23.5
- PostgreSQL version: 15.10

## Installation

1. Install gator with:
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

  - Ensure the installation worked and you have version 15+:
    ```
    psql version
    ```

  - (Linux only) Update system users postgres password (Do not forget it):
    ```
    sudo passwd postgres
    ```

  - Start PostgreSQL service
    ```
    sudo service postgresql start
    ```

  - Create a new database:
    ```
    CREATE DATABASE gator;
    ```

  - Connect to the new database:
    ```
    \c gator
    ```

  - (Linux only) Set the database users password:
    ```
    ALTER USER postgres PASSWORD 'postgres';
    ```

  - Exit the database:
    ```
    exit
    ```

## Configuration

- Create a `.gatorconfig.json` file in your home directory with the following structure:
  ```json
  {
    "db_url": "postgres://username:@localhost:5432/database?sslmode=disable"
  }
  ```

- Replace the values with your database connection string.
  - macOS will look something like this: 'postgres://yourusername:@localhost:5432/gator'
  - Linux will look something like this: 'postgres://postgres:postgres@localhost:5432/gator'

- You can test your connection string by trying to connect to psql:
  ```
  psql "your_connection_string"

## Available Commands

- Setup & Database:
```
gator help                  (List available commands and the use cases)
gator register <name> 	    (Create a new user)
gator reset 			          (Reset/clear the database)
gator login <name> 		      (Log in as a user that already exists)

- Feed management:
gator addfeed <url> 	      (Add a feed to the database)
gator feeds 			          (List all feeds)
gator follow <feed_id> 	    (Follow a feed that already exists in the database)
gator following 		        (Lists feeds the user is following)
gator unfollow <feed_id> 	  (Unfollow a feed that already exists in the database)

- Content & Updates:
gator browse [limit] 	      (View the posts, defaults to 2)
gator users 			          (List all users)
gator agg <30s> 	          (Start the aggregator. Can use preferred time intervals such as 1m 3m 6m etc...) ctrl + c to end
```