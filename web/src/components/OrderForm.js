import React, { useState } from 'react';

const OrderForm = ({ cartItems, total, onPlaceOrder }) => {
  const [couponCode, setCouponCode] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsSubmitting(true);
    
    try {
      await onPlaceOrder({ couponCode });
      setCouponCode('');
    } catch (error) {
      // Error handling is done in parent component
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="order-form">
      <h3>Place Order</h3>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Order Summary</label>
          <div style={{ backgroundColor: '#f8f9fa', padding: '1rem', borderRadius: '4px' }}>
            {cartItems.map(item => (
              <div key={item.id} style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.5rem' }}>
                <span>{item.name} x{item.quantity}</span>
                <span>${(item.price * item.quantity).toFixed(2)}</span>
              </div>
            ))}
            <div style={{ borderTop: '1px solid #dee2e6', marginTop: '0.5rem', paddingTop: '0.5rem', fontWeight: 'bold' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <span>Total:</span>
                <span>${total.toFixed(2)}</span>
              </div>
            </div>
          </div>
        </div>

        <div className="form-group">
          <label htmlFor="couponCode">Coupon Code (Optional)</label>
          <input
            type="text"
            id="couponCode"
            value={couponCode}
            onChange={(e) => setCouponCode(e.target.value)}
            placeholder="Enter coupon code"
          />
        </div>

        <button 
          type="submit" 
          className="btn btn-success"
          disabled={isSubmitting || cartItems.length === 0}
          style={{ width: '100%' }}
        >
          {isSubmitting ? 'Placing Order...' : 'Place Order'}
        </button>
      </form>
    </div>
  );
};

export default OrderForm;