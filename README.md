# Nebula

<p align="center">
  <img src="static/images/nebula-logo.png" alt="Nebula Logo" width="200">
</p>

A secure, full-stack web application designed for cybersecurity enthusiasts, hackers, and technology professionals.

## Overview

Nebula is a robust platform that connects innovators through cutting-edge cybersecurity features, technology news, and community interaction. Built with performance and security at its core, Nebula demonstrates advanced backend capabilities using Go, PostgreSQL, Redis, and the Echo framework.

## Key Features

### Authentication & Security
- Secure user authentication flow with session management
- CSRF protection and Redis-backed session store
- Industry-standard password hashing with bcrypt
- Comprehensive security headers (HSTS, CSP, X-Frame-Options)
- Rate limiting to mitigate brute-force attacks

### Content & Functionality
- Dedicated sections for News, TTPS, Hacks of Fame, and community profiles
- S3 integration for storing and serving articles and images
- Markdown rendering for article content with security sanitization
- Responsive, Matrix-inspired UI with modern interactions

### Technical Architecture
- High-performance Go backend with Echo framework
- PostgreSQL database with efficient connection pooling
- Redis for fast cache and session management
- Docker-based development environment
- Structured project layout optimized for maintainability

## Technology Stack

### Backend
- **Language**: Go (Golang)
- **Framework**: Echo
- **Database**: PostgreSQL
- **Cache/Sessions**: Redis
- **Storage**: AWS S3

### Frontend
- **HTML/CSS Framework**: Custom CSS with Matrix theme
- **JavaScript**: Vanilla JS with Matrix rain effect
- **Templating**: Go's html/template

### DevOps & Tools
- Docker & Docker Compose
- Air (live reload)
- Goose (database migrations)
- Sqlx (database access)

## Development Setup

### Prerequisites
- Go 1.24 or higher
- Docker and Docker Compose
- PostgreSQL client (optional)
- Redis client (optional)
- AWS account with S3 access (for article storage)

### Quick Start

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/nebula.git
   cd nebula
   ```

2. **Start the database services**
   ```bash
   make db-up
   ```

3. **Configure environment variables**
   Create a `.env` file in the project root with the following variables (customize as needed):
   ```
   # Database Configuration
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASS=your_password
   DB_NAME=nebula
   DB_SSL_MODE=disable

   # Redis Configuration
   REDIS_HOST=localhost
   REDIS_PORT=6379
   REDIS_PASSWORD=
   REDIS_DB=3

   # AWS S3 Configuration
   AWS_REGION=your_region
   AWS_ACCESS_KEY=your_access_key
   AWS_SECRET_KEY=your_secret_key
   AWS_BUCKET=your_bucket_name

   # Server Configuration
   PORT=7777

   # Security Settings
   SESSION_SECRET=change_this_in_production
   CSRF_SECRET=change_this_in_production
   ```

4. **Run database migrations**
   ```bash
   goose -dir migrations postgres "host=localhost user=postgres password=your_password dbname=nebula sslmode=disable" up
   ```

5. **Run the application**
   ```bash
   make dev
   ```

6. **Access the application**
   Open `http://localhost:7777` in your browser

## API Endpoints

| Method | Endpoint         | Description                            | Auth Required |
|--------|------------------|----------------------------------------|--------------|
| GET    | /                | Home page                              | No           |
| GET    | /news            | List articles from S3                  | No           |
| GET    | /news/:id        | View specific article                  | No           |
| GET    | /ttps            | Tactics, Techniques, and Procedures    | No           |
| GET    | /hof             | Hacks of Fame                          | No           |
| GET    | /who             | Innovator profiles                     | No           |
| GET    | /login           | Login page                             | No           |
| POST   | /login           | Authenticate user                      | No           |
| GET    | /register        | Registration page                      | No           |
| POST   | /register        | Create new user                        | No           |
| GET    | /logout          | Log out current user                   | Yes          |
| GET    | /images/*        | Serve images from S3                   | No           |
| POST   | /admin/images/upload | Upload image to S3                 | Yes          |

## Project Structure

```
nebula/
├── cmd/                # Application entry points
│   └── main.go         # Main server
├── controllers/        # Request handlers with business logic
├── handlers/           # Simple request handlers
├── migrations/         # Database migration files
├── pkg/                # Reusable packages
│   └── s3client/       # S3 integration
├── static/             # Static assets
│   ├── css/            # Stylesheets
│   ├── js/             # JavaScript files
│   └── images/         # Static images
├── templates/          # HTML templates
├── .env                # Environment variables
├── docker-compose.yml  # Container configuration
├── go.mod              # Go dependencies
└── Makefile            # Development commands
```

## Available Make Commands

| Command          | Description                               |
|------------------|-------------------------------------------|
| `make help`      | Display available commands                |
| `make dev`       | Run with hot-reloading                    |
| `make db-up`     | Start database containers                 |
| `make db-down`   | Stop database containers                  |
| `make db-reset`  | Reset the database                        |
| `make install`   | Install Go dependencies                   |
| `make run`       | Run without hot-reloading                 |
| `make test`      | Run tests                                 |
| `make lint`      | Run linters                               |
| `make build`     | Build the application                     |
| `make psql`      | Access PostgreSQL CLI                     |
| `make redis-cli` | Access Redis CLI                          |

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contact

For questions or feedback, please open an issue on the GitHub repository.


