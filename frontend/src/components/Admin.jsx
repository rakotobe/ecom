import { useState, useEffect } from 'react';
import { productApi, orderApi } from '../api/client';

function Admin() {
  const [products, setProducts] = useState([]);
  const [orders, setOrders] = useState([]);
  const [showForm, setShowForm] = useState(false);
  const [editingProduct, setEditingProduct] = useState(null);
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    price: '',
    currency: 'USD',
    stock: '',
  });
  const [error, setError] = useState(null);

  useEffect(() => {
    loadProducts();
    loadOrders();
  }, []);

  const loadProducts = async () => {
    try {
      const data = await productApi.getAll();
      setProducts(data || []);
    } catch (err) {
      setError(err.message);
    }
  };

  const loadOrders = async () => {
    try {
      const data = await orderApi.getAll();
      setOrders(data || []);
    } catch (err) {
      setError(err.message);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const productData = {
        name: formData.name,
        description: formData.description,
        price: parseInt(formData.price),
        currency: formData.currency,
        stock: parseInt(formData.stock),
      };

      if (editingProduct) {
        await productApi.update(editingProduct.id, productData);
      } else {
        await productApi.create(productData);
      }

      setShowForm(false);
      setEditingProduct(null);
      setFormData({ name: '', description: '', price: '', currency: 'USD', stock: '' });
      loadProducts();
    } catch (err) {
      setError(err.message);
    }
  };

  const handleEdit = (product) => {
    setEditingProduct(product);
    setFormData({
      name: product.name,
      description: product.description,
      price: product.price.toString(),
      currency: product.currency,
      stock: product.stock.toString(),
    });
    setShowForm(true);
  };

  const handleDelete = async (id) => {
    if (window.confirm('Are you sure you want to delete this product?')) {
      try {
        await productApi.delete(id);
        loadProducts();
      } catch (err) {
        setError(err.message);
      }
    }
  };

  const formatPrice = (amount, currency) => {
    const dollars = amount / 100;
    return `${currency} $${dollars.toFixed(2)}`;
  };

  const getProductName = (productId) => {
    const product = products.find(p => p.id === productId);
    return product ? product.name : `Product ${productId.substring(0, 8)}...`;
  };

  return (
    <div className="admin">
      <h2>Admin Panel</h2>
      {error && <div className="error">{error}</div>}

      <div className="admin-section">
        <h3>Products</h3>
        <button onClick={() => {
          setShowForm(!showForm);
          setEditingProduct(null);
          setFormData({ name: '', description: '', price: '', currency: 'USD', stock: '' });
        }}>
          {showForm ? 'Cancel' : 'Add New Product'}
        </button>

        {showForm && (
          <form onSubmit={handleSubmit} className="product-form">
            <input
              type="text"
              placeholder="Product Name"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              required
            />
            <textarea
              placeholder="Description"
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            />
            <input
              type="number"
              placeholder="Price (in cents)"
              value={formData.price}
              onChange={(e) => setFormData({ ...formData, price: e.target.value })}
              required
            />
            <select
              value={formData.currency}
              onChange={(e) => setFormData({ ...formData, currency: e.target.value })}
            >
              <option value="USD">USD</option>
              <option value="EUR">EUR</option>
              <option value="GBP">GBP</option>
            </select>
            <input
              type="number"
              placeholder="Stock"
              value={formData.stock}
              onChange={(e) => setFormData({ ...formData, stock: e.target.value })}
              required
            />
            <button type="submit">{editingProduct ? 'Update' : 'Create'} Product</button>
          </form>
        )}

        <div className="products-table">
          {products.map((product) => (
            <div key={product.id} className="product-row">
              <div className="product-info">
                <h4>{product.name}</h4>
                <p>{product.description}</p>
                <p>Price: {formatPrice(product.price, product.currency)}</p>
                <p>Stock: {product.stock}</p>
              </div>
              <div className="product-actions">
                <button onClick={() => handleEdit(product)}>Edit</button>
                <button onClick={() => handleDelete(product.id)}>Delete</button>
              </div>
            </div>
          ))}
        </div>
      </div>

      <div className="admin-section">
        <h3>Orders</h3>
        <div className="orders-table">
          {orders.length === 0 ? (
            <p>No orders yet</p>
          ) : (
            orders.map((order) => (
              <div key={order.id} className="order-card">
                <div className="order-header">
                  <div>
                    <strong>Order ID:</strong> {order.id.substring(0, 8)}...
                  </div>
                  <div className="order-status">{order.status}</div>
                </div>

                <div className="order-items">
                  <h4>Items:</h4>
                  {order.items.map((item, index) => (
                    <div key={index} className="order-item">
                      <div className="order-item-details">
                        <span className="item-name">{getProductName(item.product_id)}</span>
                        <span className="item-quantity">Qty: {item.quantity}</span>
                      </div>
                      <div className="order-item-price">
                        <span className="item-unit-price">{formatPrice(item.price, item.currency)} each</span>
                        <span className="item-subtotal">{formatPrice(item.subtotal, item.currency)}</span>
                      </div>
                    </div>
                  ))}
                </div>

                <div className="order-footer">
                  <div className="order-total">
                    <strong>Total:</strong> {formatPrice(order.total, order.currency)}
                  </div>
                  <div className="order-date">
                    {new Date(order.created_at).toLocaleString()}
                  </div>
                </div>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
}

export default Admin;
