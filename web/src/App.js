import React, { useState, useEffect } from 'react';
import axios from 'axios';
import ProductList from './components/ProductList';
import Cart from './components/Cart';
import OrderForm from './components/OrderForm';
import OrderList from './components/OrderList';
import './App.css';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

function App() {
  const [products, setProducts] = useState([]);
  const [cart, setCart] = useState({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [currentView, setCurrentView] = useState('products'); // 'products' or 'orders'

  useEffect(() => {
    fetchProducts();
  }, []);

  const fetchProducts = async () => {
    try {
      const response = await axios.get(`${API_BASE_URL}/api/v1/product`, {
        headers: {
          'X-API-Key': 'apitest'
        }
      });
      setProducts(response.data);
      setLoading(false);
    } catch (err) {
      setError('Failed to load products. Please try again.');
      setLoading(false);
    }
  };

  const addToCart = (product) => {
    setCart(prevCart => {
      const newCart = { ...prevCart };
      if (newCart[product.id]) {
        newCart[product.id].quantity += 1;
      } else {
        newCart[product.id] = {
          ...product,
          quantity: 1
        };
      }
      return newCart;
    });
    setSuccess('Product added to cart!');
    setTimeout(() => setSuccess(''), 2000);
  };

  const updateQuantity = (productId, quantity) => {
    if (quantity <= 0) {
      removeFromCart(productId);
      return;
    }
    setCart(prevCart => ({
      ...prevCart,
      [productId]: {
        ...prevCart[productId],
        quantity
      }
    }));
  };

  const removeFromCart = (productId) => {
    setCart(prevCart => {
      const newCart = { ...prevCart };
      delete newCart[productId];
      return newCart;
    });
  };

  const getCartItems = () => {
    return Object.values(cart);
  };

  const getCartTotal = () => {
    return Object.values(cart).reduce((total, item) => {
      return total + (item.price * item.quantity);
    }, 0);
  };

  const placeOrder = async (orderData) => {
    try {
      const orderItems = getCartItems().map(item => ({
        productId: item.id,
        quantity: item.quantity
      }));

      const payload = {
        items: orderItems,
        couponCode: orderData.couponCode
      };

      const response = await axios.post(
        `${API_BASE_URL}/api/v1/order`,
        payload,
        {
          headers: {
            'X-API-Key': 'apitest',
            'Content-Type': 'application/json'
          }
        }
      );

      setSuccess('Order placed successfully!');
      setCart({});
      setTimeout(() => setSuccess(''), 3000);
      
      return response.data;
    } catch (err) {
      setError(err.response?.data?.message || 'Failed to place order. Please try again.');
      setTimeout(() => setError(''), 3000);
      throw err;
    }
  };

  return (
    <div className="App">
      <header className="header">
        <div className="container">
          <h1>Oolio - Food Ordering</h1>
          <nav style={{ marginTop: '1rem' }}>
            <button
              onClick={() => setCurrentView('products')}
              style={{
                padding: '0.5rem 1rem',
                margin: '0 0.5rem',
                backgroundColor: currentView === 'products' ? '#007bff' : '#6c757d',
                color: 'white',
                border: 'none',
                borderRadius: '4px',
                cursor: 'pointer'
              }}
            >
              Products
            </button>
            <button
              onClick={() => setCurrentView('orders')}
              style={{
                padding: '0.5rem 1rem',
                margin: '0 0.5rem',
                backgroundColor: currentView === 'orders' ? '#007bff' : '#6c757d',
                color: 'white',
                border: 'none',
                borderRadius: '4px',
                cursor: 'pointer'
              }}
            >
              Orders
            </button>
          </nav>
        </div>
      </header>

      <main className="container">
        {error && <div className="error">{error}</div>}
        {success && <div className="success">{success}</div>}

        {currentView === 'orders' ? (
          <OrderList />
        ) : (
          <>
            {loading ? (
              <div className="loading">Loading products...</div>
            ) : (
              <div style={{ display: 'flex', gap: '2rem' }}>
                <div style={{ flex: 1 }}>
                  <ProductList 
                    products={products} 
                    onAddToCart={addToCart}
                  />
                </div>
                
                {getCartItems().length > 0 && (
                  <div>
                    <Cart 
                      cart={cart}
                      onUpdateQuantity={updateQuantity}
                      onRemoveFromCart={removeFromCart}
                      total={getCartTotal()}
                    />
                    <OrderForm 
                      cartItems={getCartItems()}
                      total={getCartTotal()}
                      onPlaceOrder={placeOrder}
                    />
                  </div>
                )}
              </div>
            )}
          </>
        )}
      </main>
    </div>
  );
}

export default App;