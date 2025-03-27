# Nebula Blog - Project Rules and Guidelines

## Project Structure
- `cmd/` - Application entry points and main server code
- `controllers/` - Request handlers that use business logic
- `handlers/` - Page rendering and route handlers
- `middlewares/` - HTTP middleware components
- `models/` - Database models and struct definitions
- `repositories/` - Database access layer
- `services/` - Business logic services
- `migrations/` - Database migration files
- `static/` - Static assets (CSS, JS, images)
- `templates/` - HTML templates
- `tmp/` - Temporary files for hot reload (not committed)

## Technology Stack
- **Backend**: Go with Echo framework
- **Hot Reloading**: Air
- **Frontend**: Minimal JS with Bootstrap CSS
- **Databases**: PostgreSQL (primary storage), Redis (caching/sessions)
- **Containerization**: Docker and Docker Compose

## Coding Standards

### Go Code
- Follow standard Go project layout
- Use consistent error handling with proper logging
- Implement middleware for auth, logging, and error handling
- Implement repository pattern for database access
- Use environment variables for configuration
- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Maximum line length of 100 characters
- Use dependency injection for better testability
- Implement idiomatic Go error handling
- Use context for cancellation and timeouts

### Database
- All schema changes must have migrations
- Use prepared statements for database queries
- Implement proper connection pooling
- Use UUIDs for primary keys
- Use foreign key constraints
- Include indexes for frequently queried columns
- Use transactions for multi-step operations
- Implement database health checks

### Redis Usage
- Use for caching frequently accessed data
- Implement for session management
- Use for rate limiting
- Set appropriate TTL for cached items
- Monitor Redis memory usage

### Frontend
- Use Bootstrap for responsive design
- Minimize custom CSS/JS
- Use semantic HTML
- Ensure responsive images
- Implement proper meta tags for SEO

### Docker
- Use docker-compose for local development
- Use named volumes for database persistence
- Use environment variables for configuration
- Set resource limits for containers
- Implement health checks for each service

## Development Workflow
1. Use feature branches
2. Run linters before committing
3. Write tests for all new features
4. Document API endpoints
5. Review database queries for performance
6. Update documentation as needed

## Security Guidelines
- No secrets in code
- Implement proper authentication
- Validate all user inputs
- Use HTTPS in production
- Implement rate limiting
- Regular security audits
- Use prepared statements to prevent SQL injection
- Implement CSRF protection
- Set secure HTTP headers
- Use secure cookie settings
- Implement content security policy

## API Design
- RESTful API design
- Consistent error responses
- Proper HTTP status codes
- Use JSON for request/response bodies
- Implement pagination for list endpoints
- Use query parameters for filtering
- Implement robust validation
- Use meaningful error messages

## Template Patterns
- Each page should have its own template file
- Use layouts for common elements (header, footer)
- Follow a consistent naming convention
- Keep templates simple and focused
- Use partial templates for reusable components

## Project-Specific Requirements
- Implement user authentication
- Create a Matrix rain animation for the home page
- Implement TTPS, HOF, and WHO pages
- Ensure responsive design for mobile
- Implement content sanitization
- Implement CSRF protection
- Set secure HTTP headers
- Use Redis for session management
- Implement rate limiting
- Create user profiles

## Development Environment Setup
1. Clone the repository
2. Start PostgreSQL and Redis with docker-compose
3. Create a `.env` file with required environment variables
4. Run with Air for hot reloading
5. Access the application at http://localhost:7777 