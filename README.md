🌌 Nebula
Nebula is a secure, full-stack web application designed to connect innovators worldwide through cutting-edge cybersecurity features, technology news, and an interactive community.
The project demonstrates backend capabilities using Go, PostgreSQL, Redis, and Echo, coupled with a frontend enhanced by Bootstrap.

🚀 Features
User Authentication:

Secure registration, login, and logout flow.
CSRF protection and session management using Redis.
Password hashing with bcrypt for enhanced security.
Dynamic Web Pages:

Homepage with a customizable Matrix rain background.
Dedicated sections for:
News: Latest technology and cybersecurity updates.
TTPS (Tactics, Techniques, and Procedures).
Hacks of Fame: A curated list of impactful hacks and their analysis.
WhoAmI: A place for innovators to share their stories.
MessageWall for communication with advanced sanitization.

Performance & Security:
HTTPS with secure headers (HSTS, Content Security Policy, etc.).
Protection against common vulnerabilities like XSS and CSRF.
Rate limiting to prevent brute-force attacks.
Scalable Architecture:

Backend built with Go for high performance and concurrency.
PostgreSQL for reliable and efficient data storage.
Redis for session and token management.
Responsive Design:

Mobile-friendly UI using Bootstrap.
Modern and sleek user interface.
🎯 Upcoming Features
RSS Feed for News: Stay updated with the latest cybersecurity and technology news delivered directly to your feed reader.
API for Articles and News: Access Nebula’s news and hacking articles programmatically for seamless integration into your applications or services.
🛠️ Technology Stack
Backend:

Go (Golang)
Echo framework
PostgreSQL
Redis
Frontend:

HTML5, CSS3
JavaScript (ES6+)
Bootstrap 5.3
Matrix Rain JavaScript effect for immersive design.
Other Tools:

Docker (for local development)
Air (live reload for development)
Goose (database migrations)
Sqlx (database interaction)
⚙️ Installation and Setup
Follow these steps to get Nebula up and running:


Run the Application
Start PostgreSQL and Redis.
Initialize the database:
bash
Copy code
goose up
Run the application:
bash
Copy code
air
The server will start on http://localhost:7777.

🌐 API Endpoints
Here’s a quick overview of the available and upcoming routes:

Public Routes
Method	Endpoint	Description
GET	/	Home page
GET	/news	Technology news page
GET	/ttps	Tactics, Techniques, and Procedures
GET	/hof	Hacks of Fame
GET	/who	Innovator stories
GET	/login	Login page
GET	/register	Registration page
Authenticated Routes
Method	Endpoint	Description
POST	/register	Register a new user
POST	/login	Log in as an existing user
GET	/logout	Log out the current user
Future API Endpoints
Method	Endpoint	Description
GET	/api/news	Retrieve the latest news articles
GET	/api/articles	Retrieve hacking articles

🛡️ Security Features
CSRF Protection: Middleware ensures that requests include a valid CSRF token.
Password Hashing: User passwords are hashed with bcrypt.
Secure Cookies: Session cookies are HTTP-only and use the SameSite attribute.
Secure Headers: HSTS, CSP, and other headers implemented via middleware
