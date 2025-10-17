# Frontend - React Application

Clean and minimal React application for the e-commerce platform.

## Structure

```
src/
├── api/                  # API client layer
│   └── client.js         # HTTP client with API methods
├── components/           # React components
│   ├── ProductList.jsx   # Product catalog display
│   ├── Basket.jsx        # Shopping basket UI
│   └── Admin.jsx         # Admin panel for management
├── App.jsx               # Main application component
├── App.css               # Application styles
└── main.jsx              # Entry point
```

## Component Architecture

### API Client (`api/client.js`)
- Centralized API communication
- Separated by domain (products, baskets, orders)
- Handles errors and response formatting

### ProductList Component
**Responsibilities**:
- Display all products
- Handle "Add to Basket" action
- Show product availability

**Props**:
- `onAddToBasket`: Callback when user adds product to basket

### Basket Component
**Responsibilities**:
- Display basket contents
- Update item quantities
- Remove items
- Clear basket
- Checkout

**Props**:
- `basketId`: Current basket ID
- `onCheckout`: Callback for checkout action

### Admin Component
**Responsibilities**:
- Product CRUD operations
- View all orders
- Manage product stock

**State Management**:
- Local component state with useState
- Form handling for product creation/editing

### App Component
**Responsibilities**:
- Global state management (basket ID, current view)
- Routing between views
- Message notifications

## State Management Strategy

This application uses React's built-in state management:
- `useState`: For component-local state
- `useEffect`: For side effects (API calls, initialization)
- `localStorage`: For basket persistence across sessions

For a larger application, consider:
- Context API for global state
- React Query for server state
- Redux for complex state management

## Development

```bash
# Install dependencies
npm install

# Start dev server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Lint
npm run lint
```

## Component Patterns

### Data Fetching Pattern
```javascript
useEffect(() => {
  loadData();
}, []);

const loadData = async () => {
  try {
    setLoading(true);
    const data = await api.getData();
    setData(data);
  } catch (err) {
    setError(err.message);
  } finally {
    setLoading(false);
  }
};
```

### Error Handling
- Try/catch in async functions
- Display error messages to user
- Graceful degradation

### Loading States
- Show loading indicators
- Prevent duplicate requests
- Disable buttons during operations

## Best Practices

1. **Component Separation**: Each component has a single responsibility
2. **Props Over State**: Pass data down, emit events up
3. **Error Boundaries**: Handle component errors gracefully
4. **Accessibility**: Use semantic HTML, proper ARIA labels
5. **Performance**: Avoid unnecessary re-renders with proper dependencies

## Future Enhancements

- Add authentication
- Implement shopping cart persistence
- Add product search and filters
- Image upload for products
- Order history for users
- Payment integration
