# Kopelko Dating App

```bash
dating-app-backend/
├── main.go              # Entry point
├── config/
│   └── config.go        # Configuration file for env vars, DB connection
├── controllers/
│   ├── auth.go          # Signup & login
│   ├── profile.go       # Profile viewing, updating
│   ├── swipe.go         # Swipe and match logic
│   ├── subscription.go  # Subscription handling
├── models/
│   ├── user.go          # User model
│   ├── profile.go       # Profile model
│   ├── swipe.go         # Swipe model
│   ├── match.go         # Match model
│   ├── subscription.go  # Subscription model
├── routes/
│   └── routes.go        # API routes setup
├── services/
│   ├── auth_service.go  # Authentication logic
│   ├── swipe_service.go # Swipe logic
│   ├── match_service.go # Match handling
│   ├── subscription_service.go # Subscription handling
└── utils/
    ├── jwt.go           # JWT utility for token generation
    ├── db.go            # Database and Redis connections
    └── cache.go         # Caching helper (Redis)
```
