const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

async function apiRequest(endpoint, options = {}) {
  const url = `${API_BASE_URL}${endpoint}`;
  const config = {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
  };

  const response = await fetch(url, config);

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Request failed' }));
    throw new Error(error.error || `HTTP error! status: ${response.status}`);
  }

  if (response.status === 204) {
    return null;
  }

  return response.json();
}

// Product API
export const productApi = {
  getAll: () => apiRequest('/products'),
  getById: (id) => apiRequest(`/products/${id}`),
  create: (data) => apiRequest('/products', {
    method: 'POST',
    body: JSON.stringify(data),
  }),
  update: (id, data) => apiRequest(`/products/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  }),
  updateStock: (id, stock) => apiRequest(`/products/${id}/stock`, {
    method: 'PATCH',
    body: JSON.stringify({ stock }),
  }),
  delete: (id) => apiRequest(`/products/${id}`, {
    method: 'DELETE',
  }),
};

// Basket API
export const basketApi = {
  create: () => apiRequest('/baskets', { method: 'POST' }),
  getById: (id) => apiRequest(`/baskets/${id}`),
  addItem: (id, productId, quantity) => apiRequest(`/baskets/${id}/items`, {
    method: 'POST',
    body: JSON.stringify({ product_id: productId, quantity }),
  }),
  updateItemQuantity: (id, productId, quantity) => apiRequest(`/baskets/${id}/items/${productId}`, {
    method: 'PATCH',
    body: JSON.stringify({ quantity }),
  }),
  removeItem: (id, productId) => apiRequest(`/baskets/${id}/items/${productId}`, {
    method: 'DELETE',
  }),
  clear: (id) => apiRequest(`/baskets/${id}/items`, {
    method: 'DELETE',
  }),
};

// Order API
export const orderApi = {
  create: (basketId) => apiRequest('/orders', {
    method: 'POST',
    body: JSON.stringify({ basket_id: basketId }),
  }),
  getAll: () => apiRequest('/orders'),
  getById: (id) => apiRequest(`/orders/${id}`),
  confirm: (id) => apiRequest(`/orders/${id}/confirm`, { method: 'POST' }),
  ship: (id) => apiRequest(`/orders/${id}/ship`, { method: 'POST' }),
  deliver: (id) => apiRequest(`/orders/${id}/deliver`, { method: 'POST' }),
  cancel: (id) => apiRequest(`/orders/${id}/cancel`, { method: 'POST' }),
};
