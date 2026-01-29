# Oolio Web UI

A minimalist React frontend for the Oolio food ordering system.

## Features

- Browse product catalog
- Add items to shopping cart
- Adjust quantities in cart
- Apply coupon codes
- Place orders with real-time feedback
- Responsive, minimalist design
- List order items with quantities and prices

## Prerequisites

- Node.js 14+ 
- npm or yarn

## Setup

1. Navigate to the web directory:
   ```bash
   cd web
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Start the development server:
   ```bash
   npm start
   ```

The UI will be available at `http://localhost:3000`

## Configuration

Create a `.env` file in the web directory to configure the API URL:

```
REACT_APP_API_URL=http://localhost:8080
```

## Usage

1. Browse available products on the main page
2. Click "Add to Cart" to add items to your shopping cart
3. Adjust quantities using the + and - buttons in the cart
4. Optionally enter a coupon code for discounts
5. Click "Place Order" to submit your order
6. Orders are processed asynchronously via the queue system

## API Integration

The UI communicates with the backend API at:
- `GET /api/v1/product` - List products
- `GET /api/v1/order` - List orders
- `POST /api/v1/order` - Place order (requires API key)

## Build for Production

```bash
npm run build
```

The build artifacts will be in the `build/` directory.
