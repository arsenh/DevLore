# DevLore

A personal programming knowledge base — a web application for storing and organizing programming knowledge, recipes, tutorials, and code snippets. Users can create accounts, write articles, and search through their personal knowledge base.

---

## Table of Contents

1. [Project Overview](#project-overview)
2. [Features](#features)
3. [Tech Stack](#tech-stack)
4. [Getting Started](#getting-started)

   * [Prerequisites](#prerequisites)
   * [Build & Run Standalone](#build--run-standalone)
   * [Run with Docker Compose](#run-with-docker-compose)
5. [Database Setup](#database-setup)
6. [Environment Variables](#environment-variables)
7. [License](#license)

---

## Project Overview

DevLore is a web application that lets you maintain your own **programming knowledge base**. It supports user accounts, article creation and editing, and an organized way to keep technical notes and code snippets searchable.

---

## Features

* User authentication (register/login)
* Create, read, update, and delete articles
* Auto‑updated timestamps
* PostgreSQL database
* Optional Docker deployment

---

## Tech Stack

* **Backend:** Go
* **Database:** PostgreSQL
* **Dev Tools:** Docker + Docker Compose

---

## Getting Started

### Prerequisites

Before running DevLore, make sure you have:

* [Go 1.21+](https://go.dev/doc/install)
* [PostgreSQL 15+](https://www.postgresql.org/download/) (if running without Docker)
* [Docker & Docker Compose](https://docs.docker.com/compose/install/) (optional)

---

## Build & Run Standalone

1. **Clone the repository**

```bash
git clone https://github.com/arsenh/DevLore.git
cd DevLore
```

2. **Edit `.env` if needed**

Example values:

```env
APP_ADDR=0.0.0.0:8080
POSTGRES_USER=postgres_user
POSTGRES_PASSWORD=qwerty
POSTGRES_DB=devlore_db
POSTGRES_HOST=localhost
POSTGRES_PORT=5432```

4. **Build the app**

```bash
go build -o devlore ./cmd/app
```

5. **Run the app**

```bash
./devlore
```

6. **Open in the browser**

Visit: `http://localhost:8080`

---

## Run with Docker Compose

DevLore includes a `docker-compose.yml` that runs both the app and PostgreSQL together.

1. **Make sure Docker & Docker Compose are installed**
3. **Start services**

```bash
docker-compose up --build
```

4. **Stop services**

```bash
docker-compose down
```

By default:

* App runs at `http://localhost:8080`
* PostgreSQL runs on port `5432`

---

## Database Setup

When running without Docker Compose, create the database manually:

```sql
CREATE DATABASE devlore_db;
CREATE USER postgres_user WITH PASSWORD 'qwerty';
GRANT ALL PRIVILEGES ON DATABASE devlore_db TO postgres_user;
```

Run migrations (if provided or using a migration tool), then start the app.

---

## Environment Variables

DevLore uses environment variables for configuration.

Here’s a typical `.env`:

```env
# App
APP_ADDR=0.0.0.0:8080

# PostgreSQL
POSTGRES_USER=postgres_user
POSTGRES_PASSWORD=qwerty
POSTGRES_DB=devlore_db
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
```

> **Note:**
> `.env` in this repo is intended for **example/testing only**.
> Do **not** commit real secrets.

---

## License

MIT License © 2026 Arsen H.
