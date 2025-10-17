# Backend - Go Clean Architecture

This backend follows Clean Architecture principles with Domain-Driven Design.

## Layer Structure

### Domain Layer (`domain/`)
Pure business logic with zero external dependencies.

**Entities** (`entity/`):
- `Product`: Product catalog item with price, stock, and metadata
- `Basket` & `BasketItem`: Shopping cart functionality
- `Order` & `OrderItem`: Order processing with status management

**Value Objects** (`value/`):
- `Money`: Represents monetary values with currency (stored in cents)
- `Quantity`: Represents item quantities with validation

**Repository Interfaces** (`repository/`):
- `ProductRepository`: Product persistence contract
- `BasketRepository`: Basket persistence contract
- `OrderRepository`: Order persistence contract

### Application Layer (`application/`)
Orchestrates domain logic to fulfill use cases.

**Services** (`service/`):
- `ProductService`: Product CRUD operations
- `BasketService`: Shopping basket management
- `OrderService`: Order creation and management (checkout)

**DTOs** (`dto/`):
- Request and response structures for API communication
- Decouples API from domain models

### Infrastructure Layer (`infrastructure/`)
Technical implementations of domain interfaces.

**Database** (`database/`):
- PostgreSQL connection management
- Database migrations
- Connection pooling

**Persistence** (`persistence/`):
- `ProductRepositoryImpl`: PostgreSQL product repository
- `BasketRepositoryImpl`: PostgreSQL basket repository
- `OrderRepositoryImpl`: PostgreSQL order repository

### API Layer (`api/`)
HTTP interface for the application.

**Handlers** (`handler/`):
- `ProductHandler`: Product endpoints
- `BasketHandler`: Basket endpoints
- `OrderHandler`: Order endpoints

**Middleware** (`middleware/`):
- CORS middleware
- Request logging

**Router** (`router/`):
- Route configuration
- Handler registration

## Running Tests

```bash
# Unit tests (domain + application)
go test ./domain/...
go test ./application/...

# Integration tests
go test ./infrastructure/...
go test ./api/...

# All tests
go test ./...

# With coverage
go test ./... -cover

# Verbose output
go test ./... -v
```

## Building

```bash
# Development
go run cmd/main.go

# Production build
go build -o ecom-backend cmd/main.go
./ecom-backend
```

## Key Design Patterns

### Repository Pattern
Domain defines interfaces, infrastructure implements them.

### Dependency Injection
All dependencies are injected through constructors in `cmd/main.go`.

### Factory Pattern
`NewX()` creates new instances, `ReconstructX()` rebuilds from persistence.

### Value Object Pattern
Immutable objects representing domain concepts (Money, Quantity).

## Best Practices

1. **Keep domain pure**: No framework dependencies in domain layer
2. **Use value objects**: Wrap primitives in meaningful types
3. **Fail fast**: Validate in constructors and return errors immediately
4. **Explicit is better**: No hidden dependencies or global state
5. **Test business logic**: Domain layer should have highest test coverage
