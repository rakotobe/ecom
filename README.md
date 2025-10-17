# E-Commerce Application

A minimal e-commerce application demonstrating **Clean Architecture** and **Clean Code** principles with Domain-Driven Design (DDD).

## Architecture Overview

This application follows **Clean Architecture** principles with strict layer separation and dependency rules. Dependencies always point inward, from outer layers to inner layers.

```
┌─────────────────────────────────────────────────────────┐
│                     API Layer                            │
│  (HTTP Handlers, Request/Response, Routing)              │
│  Dependencies: Application Layer                         │
└──────────────────┬──────────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────────┐
│                Application Layer                         │
│  (Use Cases, Services, DTOs)                            │
│  Dependencies: Domain Layer                              │
└──────────────────┬──────────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────────┐
│                  Domain Layer                            │
│  (Entities, Value Objects, Repository Interfaces)       │
│  Dependencies: NONE (Pure business logic)                │
└─────────────────────────────────────────────────────────┘
                   ▲
┌──────────────────┴──────────────────────────────────────┐
│              Infrastructure Layer                        │
│  (Database, Repositories Implementation)                 │
│  Dependencies: Domain Layer (implements interfaces)      │
└─────────────────────────────────────────────────────────┘
```

### Layer Responsibilities

#### Domain Layer (`/backend/domain`)
- **Pure business logic** with no external dependencies
- Contains:
  - **Entities** (`entity/`): Core business objects (Product, Basket, Order)
  - **Value Objects** (`value/`): Immutable values (Money, Quantity)
  - **Repository Interfaces** (`repository/`): Contracts for data persistence
- **Key Principle**: Domain layer knows nothing about databases, HTTP, or frameworks

#### Application Layer (`/backend/application`)
- **Orchestrates** domain logic to fulfill use cases
- Contains:
  - **Services** (`service/`): Use case implementations
  - **DTOs** (`dto/`): Data transfer objects for API communication
- **Dependencies**: Only depends on domain layer
- **Responsibility**: Application-specific business rules and workflows

#### Infrastructure Layer (`/backend/infrastructure`)
- **Technical implementations** of domain interfaces
- Contains:
  - **Database** (`database/`): PostgreSQL connection and migrations
  - **Persistence** (`persistence/`): Repository implementations
- **Dependencies**: Depends on domain interfaces
- **Responsibility**: Database access, external services integration

#### API Layer (`/backend/api`)
- **Thin HTTP layer** that delegates to application services
- Contains:
  - **Handlers** (`handler/`): HTTP request handlers
  - **Middleware** (`middleware/`): CORS, logging
  - **Router** (`router/`): Route configuration
- **Responsibility**: HTTP concerns only (validation, serialization)

## Features

- **Product Management**: CRUD operations for products
- **Shopping Basket**: Add/remove items, update quantities
- **Checkout**: Create orders from basket
- **Order Management**: Track order status
- **Admin Panel**: Product and order management UI

## Technology Stack

### Backend
- **Language**: Go 1.23+
- **Database**: PostgreSQL 16
- **Router**: Gorilla Mux
- **Architecture**: Clean Architecture + DDD

### Frontend
- **Framework**: React 18 with Hooks
- **Build Tool**: Vite
- **Styling**: Plain CSS
- **State Management**: React useState/useEffect

## Prerequisites

- Docker and Docker Compose (recommended)
- OR:
  - Go 1.23+
  - Node.js 20+
  - PostgreSQL 16

## Quick Start with Docker

1. **Clone and navigate to the project**:
   ```bash
   cd ecom
   ```

2. **Start all services**:
   ```bash
   docker-compose up --build
   ```

3. **Access the application**:
   - Frontend: http://localhost:3333
   - Backend API: http://localhost:8888
   - Health Check: http://localhost:8888/health

4. **Stop services**:
   ```bash
   docker-compose down
   ```

## Local Development Setup

### Backend

1. **Start PostgreSQL**:
   ```bash
   docker run --name ecom-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=ecom -p 5555:5432 -d postgres:16-alpine
   ```

2. **Configure environment**:
   ```bash
   cd backend
   cp .env.example .env
   # Edit .env if needed
   ```

3. **Install dependencies**:
   ```bash
   go mod download
   ```

4. **Run the server**:
   ```bash
   go run cmd/main.go
   ```

5. **Run tests**:
   ```bash
   # Unit tests for domain layer
   go test ./domain/...

   # All tests
   go test ./...
   ```

### Frontend

1. **Navigate to frontend**:
   ```bash
   cd frontend
   ```

2. **Configure environment**:
   ```bash
   cp .env.example .env
   # Edit .env if needed
   ```

3. **Install dependencies**:
   ```bash
   npm install
   ```

4. **Start development server**:
   ```bash
   npm run dev
   ```

5. **Build for production**:
   ```bash
   npm run build
   ```

## API Documentation

### Base URL
```
http://localhost:8888/api/v1
```

### Products

#### Create Product
```http
POST /products
Content-Type: application/json

{
  "name": "Product Name",
  "description": "Product description",
  "price": 1999,        // price in cents
  "currency": "USD",
  "stock": 100
}
```

#### Get All Products
```http
GET /products
```

#### Get Product by ID
```http
GET /products/{id}
```

#### Update Product
```http
PUT /products/{id}
Content-Type: application/json

{
  "name": "Updated Name",
  "description": "Updated description",
  "price": 2499,
  "currency": "USD"
}
```

#### Update Stock
```http
PATCH /products/{id}/stock
Content-Type: application/json

{
  "stock": 50
}
```

#### Delete Product
```http
DELETE /products/{id}
```

### Baskets

#### Create Basket
```http
POST /baskets
```

#### Get Basket
```http
GET /baskets/{id}
```

#### Add Item to Basket
```http
POST /baskets/{id}/items
Content-Type: application/json

{
  "product_id": "product-uuid",
  "quantity": 2
}
```

#### Update Item Quantity
```http
PATCH /baskets/{id}/items/{productId}
Content-Type: application/json

{
  "quantity": 3
}
```

#### Remove Item
```http
DELETE /baskets/{id}/items/{productId}
```

#### Clear Basket
```http
DELETE /baskets/{id}/items
```

### Orders

#### Create Order (Checkout)
```http
POST /orders
Content-Type: application/json

{
  "basket_id": "basket-uuid"
}
```

#### Get All Orders
```http
GET /orders
```

#### Get Order by ID
```http
GET /orders/{id}
```

## Code Organization Rationale

### Why Clean Architecture?

1. **Independence**: Business logic is independent of frameworks, databases, and UI
2. **Testability**: Domain logic can be tested without any infrastructure
3. **Maintainability**: Changes in one layer don't affect others
4. **Flexibility**: Easy to swap implementations (e.g., change database)

### Key Design Decisions

#### 1. Value Objects for Business Concepts
```go
// Money ensures currency consistency and prevents invalid states
type Money struct {
    amount   int64  // cents to avoid floating point issues
    currency string
}
```

**Benefits**:
- Type safety (can't accidentally add different currencies)
- Encapsulated validation
- Clear intent

#### 2. Entity Reconstruction Pattern
```go
// NewProduct creates a new product (generates ID, timestamps)
func NewProduct(...) (*Product, error)

// ReconstructProduct rebuilds from persistence (uses existing ID, timestamps)
func ReconstructProduct(...) *Product
```

**Benefits**:
- Clear distinction between creation and reconstruction
- Preserves domain rules
- No setters needed

#### 3. Repository Interfaces in Domain
```go
// Domain defines the contract
type ProductRepository interface {
    Save(ctx context.Context, product *Product) error
    FindByID(ctx context.Context, id string) (*Product, error)
}

// Infrastructure implements it
type ProductRepositoryImpl struct { ... }
```

**Benefits**:
- Domain owns its persistence needs
- Easy to mock for testing
- Infrastructure is a plugin

#### 4. Dependency Injection
```go
// main.go wires everything together
productRepo := persistence.NewProductRepository(db)
productService := service.NewProductService(productRepo)
productHandler := handler.NewProductHandler(productService)
```

**Benefits**:
- No global state
- Explicit dependencies
- Testable components

## Testing Strategy

### Unit Tests
- **Location**: `domain/` and `application/` layers
- **Focus**: Business logic in isolation
- **No dependencies**: No database, no HTTP, no external services
- **Example**: Money calculations, product stock reduction, basket total

```bash
go test ecom-backend/domain/...
```

### Integration Tests
- **Location**: `infrastructure/` and `api/` layers
- **Focus**: Integration with external systems
- **Dependencies**: Test database, HTTP server
- **Example**: Repository operations, HTTP endpoints

```bash
go test ecom-backend/infrastructure/...
go test ecom-backend/api/...
```

### Test Coverage Goals
- Domain layer: 90%+
- Application layer: 80%+
- Infrastructure layer: 70%+
- API layer: 70%+

## SOLID Principles Applied

### Single Responsibility Principle
- Each layer has one reason to change
- Handlers only handle HTTP concerns
- Services only orchestrate business logic
- Entities only contain business rules

### Open/Closed Principle
- Repository interfaces allow new implementations without changing domain
- Middleware pattern allows adding features without modifying handlers

### Liskov Substitution Principle
- Any repository implementation can replace another
- Mock repositories work interchangeably with real ones

### Interface Segregation Principle
- Small, focused repository interfaces
- Handlers depend only on what they need

### Dependency Inversion Principle
- High-level domain doesn't depend on low-level infrastructure
- Both depend on abstractions (interfaces)

## Project Structure

```
ecom/
├── backend/
│   ├── domain/              # Business logic core
│   │   ├── entity/          # Business entities
│   │   ├── value/           # Value objects
│   │   └── repository/      # Repository interfaces
│   ├── application/         # Use cases
│   │   ├── dto/             # Data transfer objects
│   │   └── service/         # Application services
│   ├── infrastructure/      # Technical implementations
│   │   ├── database/        # DB connection & migrations
│   │   └── persistence/     # Repository implementations
│   ├── api/                 # HTTP layer
│   │   ├── handler/         # HTTP handlers
│   │   ├── middleware/      # Middleware
│   │   └── router/          # Routing configuration
│   ├── cmd/                 # Application entry point
│   │   └── main.go          # Main file with DI setup
│   ├── Dockerfile           # Backend container
│   └── go.mod               # Go dependencies
├── frontend/
│   ├── src/
│   │   ├── api/             # API client
│   │   ├── components/      # React components
│   │   ├── App.jsx          # Main app component
│   │   └── App.css          # Styles
│   ├── Dockerfile           # Frontend container
│   └── package.json         # Node dependencies
├── docker-compose.yml       # Service orchestration
└── README.md                # This file
```

## Common Operations

### Adding a New Feature

1. **Start with Domain**: Define entities, value objects, interfaces
2. **Application Layer**: Create DTOs and services
3. **Infrastructure**: Implement repository
4. **API**: Add handlers and routes
5. **Tests**: Write unit and integration tests
6. **Frontend**: Create components if needed

### Database Migrations

Migrations run automatically on startup. To add a new migration:

1. Edit `backend/infrastructure/database/postgres.go`
2. Add migration SQL to the `migrations` slice
3. Restart the backend

### Environment Variables

#### Backend
- `DB_HOST`: PostgreSQL host (default: localhost)
- `DB_PORT`: PostgreSQL port (default: 5555 external, 5432 internal)
- `DB_USER`: Database user (default: postgres)
- `DB_PASSWORD`: Database password (default: postgres)
- `DB_NAME`: Database name (default: ecom)
- `DB_SSLMODE`: SSL mode (default: disable)
- `PORT`: Server port (default: 8888 external, 8080 internal)

#### Frontend
- `VITE_API_URL`: Backend API URL (default: http://localhost:8888/api/v1)

## Troubleshooting

### Backend won't start
- Check PostgreSQL is running: `docker ps` or `psql -U postgres -p 5555`
- Verify environment variables are set correctly
- Check logs: `docker-compose logs backend`

### Frontend can't connect to backend
- Verify VITE_API_URL is correct
- Check CORS settings in backend
- Ensure backend is running: `curl http://localhost:8888/health`

### Tests failing
- Ensure no other services are using ports 5555, 8888, or 3333
- Run `go mod tidy` to sync dependencies
- Check if database is accessible
