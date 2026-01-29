import React from 'react';

const ProductList = ({ products, onAddToCart }) => {
  return (
    <div className="products-grid">
      {products.map(product => (
        <div key={product.id} className="product-card">
          <h3 className="product-name">{product.name}</h3>
          <div className="product-price">${product.price.toFixed(2)}</div>
          <div className="product-description">
            {product.description || 'Delicious food item prepared with care.'}
          </div>
          <button 
            className="btn btn-primary"
            onClick={() => onAddToCart(product)}
          >
            Add to Cart
          </button>
        </div>
      ))}
    </div>
  );
};

export default ProductList;