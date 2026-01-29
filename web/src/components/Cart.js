import React from 'react';

const Cart = ({ cart, onUpdateQuantity, onRemoveFromCart, total }) => {
  const cartItems = Object.values(cart);

  if (cartItems.length === 0) {
    return null;
  }

  return (
    <div className="cart">
      <h3>Shopping Cart</h3>
      {cartItems.map(item => (
        <div key={item.id} className="cart-item">
          <div>
            <div>{item.name}</div>
            <small>${item.price.toFixed(2)}</small>
          </div>
          <div className="quantity-controls">
            <button 
              className="quantity-btn"
              onClick={() => onUpdateQuantity(item.id, item.quantity - 1)}
            >
              -
            </button>
            <span>{item.quantity}</span>
            <button 
              className="quantity-btn"
              onClick={() => onUpdateQuantity(item.id, item.quantity + 1)}
            >
              +
            </button>
            <button 
              className="quantity-btn"
              onClick={() => onRemoveFromCart(item.id)}
              style={{ marginLeft: '0.5rem', color: '#e74c3c' }}
            >
              ×
            </button>
          </div>
        </div>
      ))}
      <div className="cart-total">
        Total: ${total.toFixed(2)}
      </div>
    </div>
  );
};

export default Cart;