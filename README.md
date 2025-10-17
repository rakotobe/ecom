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
