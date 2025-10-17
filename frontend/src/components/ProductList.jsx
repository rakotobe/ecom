import { useState, useEffect } from 'react';
import { productApi } from '../api/client';

function ProductList({ onAddToBasket }) {
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    loadProducts();
  }, []);

  const loadProducts = async () => {
    try {
      setLoading(true);
      const data = await productApi.getAll();
      setProducts(data || []);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const formatPrice = (amount, currency) => {
    const dollars = amount / 100;
    return `${currency} $${dollars.toFixed(2)}`;
  };

  if (loading) return <div className="loading">Loading products...</div>;
  if (error) return <div className="error">Error: {error}</div>;

  return (
    <div className="product-list">
      <h2>Products</h2>
      {products.length === 0 ? (
        <p>No products available</p>
      ) : (
        <div className="products-grid">
          {products.map((product) => (
            <div key={product.id} className="product-card">
              <h3>{product.name}</h3>
              <p className="description">{product.description}</p>
              <p className="price">{formatPrice(product.price, product.currency)}</p>
              <p className="stock">Stock: {product.stock}</p>
              {product.stock > 0 ? (
                <button onClick={() => onAddToBasket(product)}>Add to Basket</button>
              ) : (
                <button disabled>Out of Stock</button>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default ProductList;
