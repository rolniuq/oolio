import React, { useState, useEffect } from 'react';
import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

function OrderList() {
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [stats, setStats] = useState({});
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [expandedOrder, setExpandedOrder] = useState(null);

  useEffect(() => {
    fetchOrders();
    const interval = setInterval(fetchOrders, 10000); // Refresh every 10 seconds
    return () => clearInterval(interval);
  }, []);

  const fetchOrders = async () => {
    try {
      setError('');
      setLoading(true);
      const response = await axios.get(`${API_BASE_URL}/api/v1/order`, {
        headers: {
          'X-API-Key': 'apitest'
        }
      });
      
      // Handle different response structures
      const data = response.data;
      if (Array.isArray(data)) {
        // Response is directly an array of orders
        setOrders(data);
        setStats(calculateStats(data));
      } else if (data.orders && Array.isArray(data.orders)) {
        // Response has orders array and possibly stats
        setOrders(data.orders);
        setStats(data.stats || calculateStats(data.orders));
      } else {
        // Fallback: treat response as orders array
        setOrders(data || []);
        setStats(calculateStats(data || []));
      }
      setLoading(false);
    } catch (err) {
      console.error('Failed to fetch orders:', err);
      setError(`Failed to load orders: ${err.response?.data?.message || err.message || 'Unknown error'}`);
      setLoading(false);
    }
  };

  const calculateStats = (ordersArray) => {
    const stats = {
      total: ordersArray.length,
      pending: 0,
      processing: 0,
      completed: 0,
      failed: 0
    };
    
    ordersArray.forEach(order => {
      const status = order.status?.toLowerCase();
      if (stats.hasOwnProperty(status)) {
        stats[status]++;
      }
    });
    
    return stats;
  };

  const getStatusColor = (status) => {
    switch (status) {
      case 'completed':
        return '#28a745';
      case 'processing':
        return '#ffc107';
      case 'pending':
        return '#17a2b8';
      case 'failed':
        return '#dc3545';
      default:
        return '#6c757d';
    }
  };

  const formatDate = (timestamp) => {
    return new Date(timestamp).toLocaleString();
  };

  const filteredOrders = orders.filter(order => {
    const matchesSearch = order.id?.toString().toLowerCase().includes(searchTerm.toLowerCase()) ||
                         order.customer?.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesStatus = statusFilter === 'all' || order.status === statusFilter;
    return matchesSearch && matchesStatus;
  });

  const toggleOrderExpansion = (orderId) => {
    setExpandedOrder(expandedOrder === orderId ? null : orderId);
  };

  if (loading) {
    return (
      <div style={{ padding: '2rem', textAlign: 'center' }}>
        <div className="loading">Loading orders...</div>
      </div>
    );
  }

  return (
    <div style={{ padding: '2rem' }}>
      <h2>Order Management</h2>
      
      {error && <div className="error" style={{ marginBottom: '1rem' }}>{error}</div>}
      
      {/* Filters */}
      <div style={{ 
        display: 'flex', 
        gap: '1rem', 
        marginBottom: '2rem',
        flexWrap: 'wrap'
      }}>
        <input
          type="text"
          placeholder="Search by order ID or customer..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          style={{
            padding: '0.5rem',
            border: '1px solid #ddd',
            borderRadius: '4px',
            flex: '1',
            minWidth: '200px'
          }}
        />
        <select
          value={statusFilter}
          onChange={(e) => setStatusFilter(e.target.value)}
          style={{
            padding: '0.5rem',
            border: '1px solid #ddd',
            borderRadius: '4px',
            minWidth: '150px'
          }}
        >
          <option value="all">All Status</option>
          <option value="pending">Pending</option>
          <option value="processing">Processing</option>
          <option value="completed">Completed</option>
          <option value="failed">Failed</option>
        </select>
      </div>
      
      {/* Order Statistics */}
      <div style={{ 
        display: 'grid', 
        gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', 
        gap: '1rem', 
        marginBottom: '2rem' 
      }}>
        <div style={{ 
          padding: '1rem', 
          border: '1px solid #ddd', 
          borderRadius: '8px',
          backgroundColor: '#f8f9fa'
        }}>
          <h4 style={{ margin: '0 0 0.5rem 0', color: '#495057' }}>Total Orders</h4>
          <p style={{ fontSize: '1.5rem', fontWeight: 'bold', margin: '0' }}>
            {stats.total || 0}
          </p>
        </div>
        
        <div style={{ 
          padding: '1rem', 
          border: '1px solid #ddd', 
          borderRadius: '8px',
          backgroundColor: '#d4edda'
        }}>
          <h4 style={{ margin: '0 0 0.5rem 0', color: '#155724' }}>Completed</h4>
          <p style={{ fontSize: '1.5rem', fontWeight: 'bold', margin: '0', color: '#28a745' }}>
            {stats.completed || 0}
          </p>
        </div>
        
        <div style={{ 
          padding: '1rem', 
          border: '1px solid #ddd', 
          borderRadius: '8px',
          backgroundColor: '#fff3cd'
        }}>
          <h4 style={{ margin: '0 0 0.5rem 0', color: '#856404' }}>Processing</h4>
          <p style={{ fontSize: '1.5rem', fontWeight: 'bold', margin: '0', color: '#ffc107' }}>
            {stats.processing || 0}
          </p>
        </div>
        
        <div style={{ 
          padding: '1rem', 
          border: '1px solid #ddd', 
          borderRadius: '8px',
          backgroundColor: '#cce5ff'
        }}>
          <h4 style={{ margin: '0 0 0.5rem 0', color: '#004085' }}>Pending</h4>
          <p style={{ fontSize: '1.5rem', fontWeight: 'bold', margin: '0', color: '#17a2b8' }}>
            {stats.pending || 0}
          </p>
        </div>
        
        <div style={{ 
          padding: '1rem', 
          border: '1px solid #ddd', 
          borderRadius: '8px',
          backgroundColor: '#f8d7da'
        }}>
          <h4 style={{ margin: '0 0 0.5rem 0', color: '#721c24' }}>Failed</h4>
          <p style={{ fontSize: '1.5rem', fontWeight: 'bold', margin: '0', color: '#dc3545' }}>
            {stats.failed || 0}
          </p>
        </div>
      </div>

      {/* Order List */}
      <div style={{ 
        border: '1px solid #ddd', 
        borderRadius: '8px', 
        overflow: 'hidden',
        backgroundColor: 'white'
      }}>
        <div style={{ 
          padding: '1rem', 
          backgroundColor: '#f8f9fa', 
          borderBottom: '1px solid #ddd',
          fontWeight: 'bold',
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center'
        }}>
          <span>Order List ({filteredOrders.length} orders)</span>
          <button 
            onClick={fetchOrders}
            style={{ 
              padding: '0.25rem 0.5rem',
              backgroundColor: '#007bff',
              color: 'white',
              border: 'none',
              borderRadius: '4px',
              cursor: 'pointer',
              fontSize: '0.8rem'
            }}
          >
            Refresh
          </button>
        </div>
        
        {filteredOrders.length === 0 ? (
          <div style={{ padding: '2rem', textAlign: 'center' }}>
            <p style={{ color: '#6c757d' }}>
              {searchTerm || statusFilter !== 'all' 
                ? 'No orders found matching your criteria.' 
                : 'No orders available.'}
            </p>
          </div>
        ) : (
          <div style={{ overflowX: 'auto' }}>
            <table style={{ width: '100%', borderCollapse: 'collapse' }}>
              <thead>
                <tr style={{ backgroundColor: '#f8f9fa', borderBottom: '2px solid #ddd' }}>
                  <th style={{ padding: '0.75rem', textAlign: 'left', fontWeight: 'bold' }}>Order ID</th>
                  <th style={{ padding: '0.75rem', textAlign: 'left', fontWeight: 'bold' }}>Customer</th>
                  <th style={{ padding: '0.75rem', textAlign: 'left', fontWeight: 'bold' }}>Status</th>
                  <th style={{ padding: '0.75rem', textAlign: 'left', fontWeight: 'bold' }}>Total</th>
                  <th style={{ padding: '0.75rem', textAlign: 'left', fontWeight: 'bold' }}>Created</th>
                  <th style={{ padding: '0.75rem', textAlign: 'center', fontWeight: 'bold' }}>Actions</th>
                </tr>
              </thead>
              <tbody>
                {filteredOrders.map((order, index) => (
                  <React.Fragment key={order.id || index}>
                    <tr style={{ 
                      borderBottom: '1px solid #eee',
                      backgroundColor: index % 2 === 0 ? '#ffffff' : '#f8f9fa'
                    }}>
                      <td style={{ padding: '0.75rem', fontFamily: 'monospace' }}>
                        #{order.id || 'N/A'}
                      </td>
                      <td style={{ padding: '0.75rem' }}>
                        {order.customer || 'Unknown'}
                      </td>
                      <td style={{ padding: '0.75rem' }}>
                        <span style={{
                          padding: '0.25rem 0.5rem',
                          borderRadius: '4px',
                          fontSize: '0.8rem',
                          fontWeight: 'bold',
                          backgroundColor: getStatusColor(order.status),
                          color: 'white'
                        }}>
                          {order.status || 'unknown'}
                        </span>
                      </td>
                      <td style={{ padding: '0.75rem', fontWeight: 'bold' }}>
                        ${(order.total || 0).toFixed(2)}
                      </td>
                      <td style={{ padding: '0.75rem', fontSize: '0.9rem' }}>
                        {formatDate(order.createdAt || Date.now())}
                      </td>
                      <td style={{ padding: '0.75rem', textAlign: 'center' }}>
                        <button
                          onClick={() => toggleOrderExpansion(order.id)}
                          style={{
                            padding: '0.25rem 0.5rem',
                            backgroundColor: '#6c757d',
                            color: 'white',
                            border: 'none',
                            borderRadius: '4px',
                            cursor: 'pointer',
                            fontSize: '0.8rem'
                          }}
                        >
                          {expandedOrder === order.id ? 'Hide' : 'Details'}
                        </button>
                      </td>
                    </tr>
                    {expandedOrder === order.id && (
                      <tr>
                        <td colSpan="6" style={{ 
                          padding: '1rem',
                          backgroundColor: '#f1f3f4',
                          borderBottom: '2px solid #ddd'
                        }}>
                          <div style={{ fontSize: '0.9rem' }}>
                            <h5 style={{ margin: '0 0 0.5rem 0' }}>Order Details</h5>
                            {order.items && order.items.length > 0 ? (
                              <div>
                                <strong>Items:</strong>
                                <ul style={{ margin: '0.5rem 0', paddingLeft: '1.5rem' }}>
                                  {order.items.map((item, itemIndex) => (
                                    <li key={itemIndex}>
                                      {item.name} x{item.quantity} - ${(item.price * item.quantity).toFixed(2)}
                                    </li>
                                  ))}
                                </ul>
                              </div>
                            ) : (
                              <p style={{ color: '#6c757d' }}>No items details available</p>
                            )}
                            {order.notes && (
                              <div style={{ marginTop: '0.5rem' }}>
                                <strong>Notes:</strong> {order.notes}
                              </div>
                            )}
                          </div>
                        </td>
                      </tr>
                    )}
                  </React.Fragment>
                ))}
              </tbody>
            </table>
          </div>
        )}
        
        <div style={{ 
          padding: '0.5rem 1rem',
          borderTop: '1px solid #eee',
          fontSize: '0.8rem',
          color: '#6c757d',
          textAlign: 'right'
        }}>
          Last updated: {formatDate(Date.now())}
        </div>
      </div>

      {/* Instructions */}
      <div style={{ 
        marginTop: '2rem', 
        padding: '1rem', 
        backgroundColor: '#e9ecef', 
        borderRadius: '8px',
        fontSize: '0.9rem'
      }}>
        <h4>How to use:</h4>
        <ul style={{ margin: '0.5rem 0', paddingLeft: '1.5rem' }}>
          <li>Place orders from the main page using the product catalog</li>
          <li>Orders are processed asynchronously through a queue system</li>
          <li>Check this page for real-time order statistics</li>
          <li>Valid coupon codes: HAPPYHRS (10% off), FIFTYOFF (50% off)</li>
        </ul>
      </div>
    </div>
  );
}

export default OrderList;