import { useState, useEffect } from 'react';
import { basketApi } from '../api/client';

function Basket({ basketId, onCheckout }) {
  const [basket, setBasket] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (basketId) {
      loadBasket();
    }
  }, [basketId]);

  const loadBasket = async () => {
    try {
      setLoading(true);
      const data = await basketApi.getById(basketId);
      setBasket(data);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const updateQuantity = async (productId, newQuantity) => {
    try {
      const data = await basketApi.updateItemQuantity(basketId, productId, newQuantity);
      setBasket(data);
    } catch (err) {
      setError(err.message);
    }
  };

  const removeItem = async (productId) => {
    try {
      const data = await basketApi.removeItem(basketId, productId);
      setBasket(data);
    } catch (err) {
      setError(err.message);
    }
  };

  const clearBasket = async () => {
    try {
      const data = await basketApi.clear(basketId);
      setBasket(data);
    } catch (err) {
      setError(err.message);
    }
  };

  const formatPrice = (amount, currency) => {
    const dollars = amount / 100;
    return `${currency} $${dollars.toFixed(2)}`;
  };

  if (!basketId) {
    return <div className="basket">No basket created</div>;
  }

  if (loading && !basket) return <div className="loading">Loading basket...</div>;
  if (error) return <div className="error">Error: {error}</div>;
  if (!basket) return null;

  return (
    <div className="basket">
      <h2>Shopping Basket</h2>
      {basket.items.length === 0 ? (
        <p>Your basket is empty</p>
      ) : (
        <>
          <div className="basket-items">
            {basket.items.map((item) => (
              <div key={item.product_id} className="basket-item">
                <div className="item-details">
                  <p className="item-price">{formatPrice(item.price, item.currency)}</p>
                  <p className="item-subtotal">
                    Subtotal: {formatPrice(item.subtotal, item.currency)}
                  </p>
                </div>
                <div className="item-controls">
                  <button onClick={() => updateQuantity(item.product_id, item.quantity - 1)}>
                    -
                  </button>
                  <span>{item.quantity}</span>
                  <button onClick={() => updateQuantity(item.product_id, item.quantity + 1)}>
                    +
                  </button>
                  <button onClick={() => removeItem(item.product_id)} className="remove-btn">
                    Remove
                  </button>
                </div>
              </div>
            ))}
          </div>
          <div className="basket-summary">
            <p className="total">Total: {formatPrice(basket.total, basket.currency)}</p>
            <p className="item-count">Items: {basket.item_count}</p>
            <div className="basket-actions">
              <button onClick={clearBasket} className="clear-btn">
                Clear Basket
              </button>
              <button onClick={() => onCheckout(basketId)} className="checkout-btn">
                Checkout
              </button>
            </div>
          </div>
        </>
      )}
    </div>
  );
}

export default Basket;
