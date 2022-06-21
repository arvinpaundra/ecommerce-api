# Ecommerce-api

Arvin Paundra Ardana - Backend Developer Internship.

### About

This was experimental RESTful API for ecommerce with MVP (Minimum Viable Product). I built this API using these:

- [Go](https://go.dev) for the main recipe to build the API
- [Postgres](https://www.postgresql.org) for store the whole data
- [Gorm](https://gorm.io) to simplifier the query database
- [Mux](https://github.com/gorilla/mux) to handle the routes

In this project, I also use [JWT](https://jwt.io) for authentication and authorization with package from [dgrijalva](https://github.com/dgrijalva/jwt-go). After the API was created, I'm dockerizing this API and upload it into [Docker](https://hub.docker.com) registry.

Since this was an experimental and my first project with golang, there are might be tons of errors will appear.

### Purpose

The purpose of this project is to complete the internship test given by Synapsis.id which build an Ecommerce with RESTful API architecture. Maybe this isn't the best golang API, but at least this is one of my best attempts right now dealing with golang.

### Usage

#### 1. Using Git

Clone this repo if you want to try this API with `git clone`. Try this:

```
$ git clone https://github.com/arvinpaundra/ecommerce-api
```

This was the project structures of ecommerce-api repo after you clone it.

```
├── api // all logic goes here
│   ├── auth
│   ├── controllers
|   ├── middlewares
│   ├── models
│   ├── responses
|   ├── utils
|   |   └── formaterror
│   └── server.go
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
└── main.go // the entry point of this golang app
```

Before you run the whole app, you should install the go modules first. Try this:

```
$ go get github.com/badoux/checkmail
$ go get github.com/dgrijalva/jwt-go
$ go get github.com/gorilla/mux
$ go get github.com/jinzhu/gorm
$ go get github.com/joho/godotenv
$ go get golang.org/x/crypto
```

After modules installed, you can simply run with execute the entry point `main.go`. Try this:

```
$ go run main.go
```

Lastly, you can hit this url endpoint to see response of the app. Try this:

```
http://localhost:5000/
```

#### 2. Using Docker

As mention before I'm dockerize this API, so you can use docker to perform this API without installing any modules or something except the docker itself.

First moment you can pull my image from docker registry. Try this:

```
$ docker pull arvinpaundra/ecommerce-api:latest
```

Since this app run multiple services, you can run services together with command `docker-compose`. Try this:

```
$ docker-compose up --build -d
```

You can check the which service is running. Try this:

```
$ docker ps
```

Lastly, hit this url endpoint to check the service is working. Try this:

- Home endpoint

```
http://localhost:5000/
```

- PgAdmin

```
http://localhost:5050/
```

### Note :

Don't forget to create an .env file to store and hide your credentials.

Add this inside your .env file:

```
# Database Live
API_SECRET=your_secret
DB_HOST=your_db_host    # if you're going to use docker uncomment this
# DB_HOST=your_host     # if you're not use docker uncomment this
DB_DRIVER=postgres      # driver name of postgresql
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
DB_PORT=5432            # default port of postgres

# Used by pgadmin service
PGADMIN_DEFAULT_EMAIL=your_email            # your email to login into pgadmin
PGADMIN_DEFAULT_PASSWORD=your_password      # your password to login into pgadmin
```
