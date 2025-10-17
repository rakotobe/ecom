import { useState, useEffect } from 'react';
import ProductList from './components/ProductList';
import Basket from './components/Basket';
import Admin from './components/Admin';
import { basketApi, orderApi } from './api/client';
import './App.css';

function App() {
  const [basketId, setBasketId] = useState(null);
  const [currentView, setCurrentView] = useState('shop');
  const [message, setMessage] = useState(null);

  useEffect(() => {
    initializeBasket();
  }, []);

  const initializeBasket = async () => {
    let storedBasketId = localStorage.getItem('basketId');

    if (storedBasketId) {
      try {
        await basketApi.getById(storedBasketId);
        setBasketId(storedBasketId);
        return;
      } catch (err) {
        localStorage.removeItem('basketId');
      }
    }

    try {
      const newBasket = await basketApi.create();
      localStorage.setItem('basketId', newBasket.id);
      setBasketId(newBasket.id);
    } catch (err) {
      showMessage('Failed to create basket: ' + err.message, 'error');
    }
  };

  const handleAddToBasket = async (product) => {
    try {
      await basketApi.addItem(basketId, product.id, 1);
      showMessage(`Added ${product.name} to basket`, 'success');
    } catch (err) {
      showMessage('Failed to add item: ' + err.message, 'error');
    }
  };

  const handleCheckout = async (basketId) => {
    try {
      const order = await orderApi.create(basketId);
      showMessage(`Order created successfully! Order ID: ${order.id}`, 'success');
      localStorage.removeItem('basketId');
      await initializeBasket();
    } catch (err) {
      showMessage('Checkout failed: ' + err.message, 'error');
    }
  };

  const showMessage = (text, type) => {
    setMessage({ text, type });
    setTimeout(() => setMessage(null), 5000);
  };

  return (
    <div className="app">
      <header>
        <h1>E-Commerce Store</h1>
        <nav>
          <button
            className={currentView === 'shop' ? 'active' : ''}
            onClick={() => setCurrentView('shop')}
          >
            Shop
          </button>
          <button
            className={currentView === 'basket' ? 'active' : ''}
            onClick={() => setCurrentView('basket')}
          >
            Basket
          </button>
          <button
            className={currentView === 'admin' ? 'active' : ''}
            onClick={() => setCurrentView('admin')}
          >
            Admin
          </button>
        </nav>
      </header>

      {message && (
        <div className={`message ${message.type}`}>
          {message.text}
        </div>
      )}

      <main>
        {currentView === 'shop' && <ProductList onAddToBasket={handleAddToBasket} />}
        {currentView === 'basket' && <Basket basketId={basketId} onCheckout={handleCheckout} />}
        {currentView === 'admin' && <Admin />}
      </main>
    </div>
  );
}

export default App;
