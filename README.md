# Dating App Backend

A backend system for a dating app (similar to T\*nder/B\*mble) built using Golang and the Echo framework. The app includes features like profile management, swipe functionality, subscriptions, and payments. The backend uses PostgreSQL as the primary database and GORM for ORM.

## Features

- User Registration & Authentication (JWT)
- Profile Management with verification label for premium users
- Swipe functionality with swipe quota limit and view restriction
- Premium feature subscriptions with options for verified label or no swipe quota
- Payments integration with multiple providers
- Database schema with constraints and relationships to support complex functionalities
- Custom middleware for authentication

## Prerequisites

- **Golang** >= 1.18
- **PostgreSQL** >= 13
- **Docker** (for containerized deployment)
- **Make** (optional, for convenience in running commands)

## Setup & Installation

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/fanajib5/kopelko-dating-app-backend.git
   cd kopelko-dating-app-backend
   ```

2. **Install Dependencies**:

   ```bash
   go mod download
   ```

3. **Environment Variables**:
   Create a `.env` file at the root level with the following variables:

   ```dotenv
   # Application configuration
   API_PORT=3000
   LIMIT_SWIPE=10

   # Database credentials
   DB_HOST=localhost
   DB_USER=yourusername
   DB_PASSWORD=yourpassword
   DB_NAME=yourdatabase
   DB_PORT=5432
   DB_TIMEZONE=youtimezone

   # JWT credential
   JWT_SECRET=yourjwtsecret
   ```

4. **Database Setup**:
   Run migrations to set up the PostgreSQL database tables and enum types:

   ```bash
   go run databases/migrations/migrate.go
   ```

   Run seeders for initial premium feature records:

   ```bash
   go run databases/seeders/seed.go
   ```

5. **Run the Application**:
   Start the server using:

   ```bash
   go run main.go
   ```

   The API will be available at `http://localhost:8080`.

## Code Structure

```bash
Kopelko-Dating-App/
├── .env                        # Environment variables for sensitive data
├── .env.example                # Sample environment file for reference
├── .gitignore                  # Git ignore file to avoid committing unnecessary files
├── go.mod                      # Go module file, manages project dependencies
├── go.sum                      # Dependency checksum file
├── Kopelko Dating App.postman_collection.json  # Postman collection for testing endpoints
├── main.go                     # Main entry point for the application
├── README.md                   # Documentation for the project setup, structure, and usage
├── .vscode/                    # VS Code-specific settings
│   └── launch.json             # Configuration for debugging the project in VS Code
├── config/                     # Configuration files and setup
│   └── config.go               # Loads and parses environment variables
├── controllers/                # Handles HTTP request processing and response generation
│   ├── auth.go                 # Controller for authentication endpoints
│   ├── profile.go              # Controller for profile-related endpoints
│   ├── subscription.go         # Controller for subscription-related endpoints
│   └── swipe.go                # Controller for swipe functionality
├── databases/                  # Database migration and seeding management
│   ├── migrations/             # Database migration files
│   │   ├── migration.go        # Migration setup using GORM
│   │   └── schema.sql          # SQL script for initial schema setup (optional)
│   └── seeders/                # Seed data setup for initial database population
│       ├── seeder.go           # Seed management using GORM
│       └── seeder.sql          # SQL file for initial data (optional)
├── dto/                        # Data Transfer Objects for request and response validation
│   ├── login_request.go        # DTO for login request validation
│   └── register_request.go     # DTO for register request validation
├── middlewares/                # Custom middleware functions for the app
│   └── middleware.go           # Authentication and logging middleware setup
├── models/                     # Defines database schemas using GORM models
│   ├── match.go                # Match model
│   ├── premium_feature.go      # Premium feature model
│   ├── profile.go              # Profile model
│   ├── profile_view.go         # Profile view model
│   ├── subscription.go         # Subscription model
│   ├── swipe.go                # Swipe model
│   └── user.go                 # User model
├── repositories/               # Data access layer to manage database interactions
│   ├── premium_feature.go      # Repository for premium feature operations
│   ├── profile.go              # Repository for profile operations
│   ├── profile_test.txt        # Placeholder test file for profile repository
│   ├── profile_view.go         # Repository for profile view operations
│   ├── subscription.go         # Repository for subscription operations
│   ├── swipe.go                # Repository for swipe operations
│   ├── user.go                 # Repository for user operations
│   └── user_test.txt           # Placeholder test file for user repository
├── routes/                     # API route registration
│   └── routes.go               # Define and group routes for each resource
├── services/                   # Core business logic for different functionalities
│   ├── auth.go                 # Service for authentication logic
│   ├── match.go                # Service for match-related logic
│   ├── profile.go              # Service for profile-related logic
│   ├── subscription.go         # Service for subscription-related logic
│   └── swipe.go                # Service for swipe functionality
├── tests/                      # Testing resources
│   └── Kopelko Dating App.postman_collection.json  # Postman tests for API endpoints
└── utils/                      # Utility functions for the app
    ├── db.go                   # Database connection setup
    ├── httphelper.go           # Helper functions for HTTP responses
    ├── jwt.go                  # JWT token generation and verification
    └── validator.go            # Custom validator functions
```

### Explanation of Key Components

- **Main Application (`main.go`)**:
  The main entry point sets up the Echo framework and initializes configuration, routing, and middleware. It brings together the components in `controllers`, `routes`, and `middlewares`.

- **Configuration (`config`)**:
  This package loads environment variables from `.env` and manages app configuration using `config.go`.

- **Controllers**:
  These files contain HTTP handler functions for different routes. Each controller organizes endpoint logic (authentication, profile handling, subscriptions, and swipes), calling corresponding services.

- **Database Management (`databases`)**:
  Contains migration and seeding scripts. The `migrations` folder holds database schema creation scripts, while the `seeders` folder provides initial data for testing or development purposes.

- **Data Transfer Objects (`dto`)**:
  Defines structures to handle data validation for incoming requests, allowing separation between request validation and business logic.

- **Middleware**:
  Manages cross-cutting concerns like logging, authentication, and request validation, applied globally or to specific routes.

- **Models**:
  Represents the database schema using GORM. Each model maps to a table and defines relationships (e.g., `User` model, `Profile` model).

- **Repositories**:
  The data access layer interacts with the database, separating raw data operations (e.g., queries, CRUD operations) from business logic.

- **Routes**:
  Registers and groups the application's endpoints, simplifying route management and separation between API layers.

- **Services**:
  Encapsulates business logic for each feature. For example, `auth.go` handles authentication processes (e.g., registration, login), while `profile.go` manages profile operations (e.g., viewing profiles).

- **Tests**:
The tests folder contains resources and tools for testing, including Postman collections for automated API testing.

- **Utilities**:
  Shared functions, such as database initialization (`db.go`), token management (`jwt.go`), HTTP response helpers (`httphelper.go`), and custom validators (`validator.go`), are here to avoid duplicating code.

## Testing

Unit and integration tests are located in the `tests` folder. The `sqlmock` package is used for mocking database operations with GORM.

To run tests:

```bash
go test ./...
```

## Linting

Code linting is handled by `golangci-lint`. To run linting, make sure `golangci-lint` is installed:

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

Then, run linting:

```bash
golangci-lint run
```

## Deployment

### Docker

1. **Build Docker Image**:

   ```bash
   docker build -t dating-app-backend .
   ```

2. **Run Docker Container**:

   ```bash
   docker run -d -p 8080:8080 --env-file .env dating-app-backend
   ```

   The application should now be accessible on `http://localhost:8080`.

### Kubernetes

1. **Set Up Deployment and Service YAMLs**:
   Create Kubernetes deployment and service YAML files for the application and PostgreSQL database.

2. **Apply YAML Configurations**:

   ```bash
   kubectl apply -f k8s/dating-app-deployment.yaml
   kubectl apply -f k8s/postgres-deployment.yaml
   ```

3. **Access the Application**:
   The service can be accessed through the cluster IP or set up as a LoadBalancer.

### Heroku

To deploy to Heroku (assuming you have a Heroku account and the CLI installed):

1. **Create a New Heroku App**:

   ```bash
   heroku create
   ```

2. **Add PostgreSQL Add-on**:

   ```bash
   heroku addons:create heroku-postgresql:hobby-dev
   ```

3. **Push to Heroku**:

   ```bash
   git push heroku main
   ```

4. **Set Environment Variables**:

   ```bash
   heroku config:set JWT_SECRET=your_secret_key
   ```

5. **Run Migrations**:

   ```bash
   heroku run go run migrations/migrate.go
   ```

## API Endpoints

- **Authentication**:
  - `POST /api/register`: Register a new user
  - `POST /api/login`: User login
- **Profile**:
  - `GET /api/users/profiles/me`: View the authenticated profile itselg
  - `GET /api/users/profiles/random`: Get the other profile randomly
- **Swipe**:
  - `POST /users/swipes/:target_user_id`: Swipe left (pass) or right (like) on a profile
- **Subscription**:
  - `POST /users/subscriptions`: Purchase a premium feature

## Sample Requests

### Swipe Example

To swipe on a profile:

```http
POST /users/swipes/2
Content-Type: application/json
Authorization: Bearer <JWT>

{
    "type": "like"
}
```

### Subscription Example

To purchase a subscription:

```http
POST /users/subscriptions
Content-Type: application/json
Authorization: Bearer <JWT>

{
    "feature_id": 1
}
```

## License

This project is licensed under the MIT License. See the `LICENSE` file for more details.

---

With this `README.md`, you should have a comprehensive guide covering all aspects of setup, usage, testing, deployment, and contributing guidelines for the dating app backend.
