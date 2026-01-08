import React, { useState, useEffect } from 'react';
import api from '../../api/api';
import { Order } from '../../types';
import styles from './AdminOrders.module.css'

const AdminOrders: React.FC = () => {
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [statusChanging, setStatusChanging] = useState<string | null>(null);

  useEffect(() => {
    fetchAllOrders();
  }, []);

  const fetchAllOrders = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await api.getAllOrders();

      const data = response.data === null ? [] : response.data;
    
      if (Array.isArray(data)) {
        setOrders(data);
      } else {
        console.error('Ожидался массив заказов, но получено:', data);
        setOrders([]);
      }
    } catch (err: any) {
      console.error('Ошибка загрузки заказов:', err);
      setError('Не удалось загрузить заказы. Пожалуйста, попробуйте позже.');
    } finally {
      setLoading(false);
    }
  };

  const handleStatusChange = async (orderId: string, newStatus: string) => {
    try {
      setStatusChanging(orderId);
      await api.updateOrderStatus(orderId, newStatus);
      
      setOrders(prevOrders =>
        prevOrders.map(order =>
          order.id === orderId ? { ...order, status_name: newStatus } : order
        )
      );
      
      alert('Статус заказа обновлен!');
    } catch (err: any) {
      console.error('Ошибка обновления статуса:', err);
      alert('Не удалось обновить статус заказа');
    } finally {
      setStatusChanging(null);
    }
  };

  const getStatusOptions = () => {
    return [
      { value: 'ordered', label: 'Забронировано' },
      { value: 'delivered', label: 'Доставлено' },
      { value: 'cancelled', label: 'Отменено' },
    ];
  };

  const formatDate = (dateString: string): string => {
    try {
      const date = new Date(dateString);
      return date.toLocaleDateString('ru-RU', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
        timeZone: 'UTC'
      });
    } catch (e) {
      return 'Неверная дата';
    }
  };

    const getStatusClass = (status: string) => {
    switch(status) {
      case 'ordered': return styles.statusOrdered;
      case 'cancelled': return styles.statusCancelled;
      case 'delivered': return styles.statusDelivered;
      default: return styles.statusNew;
    }
  };

  const getStatusText = (status: string): string => {
    const statusMap: Record<string, string> = {
      'ordered': 'Забронировано',
      'delivered': 'Доставлено',
      'cancelled': 'Отменено',
    };
    return statusMap[status.toLowerCase()] || status;
  };

  if (loading) {
    return (
      <div className={styles.loadingContainer}>
        <div className={styles.spinner}></div>
        <p>Загрузка заказов...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className={styles.errorContainer}>
        <p className={styles.errorText}>{error}</p>
        <button onClick={fetchAllOrders} className={styles.retryButton}>
          Попробовать снова
        </button>
      </div>
    );
  }

  const displayOrders = Array.isArray(orders) ? orders : [];

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <h2 className={styles.title}>Управление заказами</h2>
        <button onClick={fetchAllOrders} className={styles.refreshButton}>
          Обновить
        </button>
      </header>

      <div className={styles.tableContainer}>
        {displayOrders.length === 0 ? (
          <div className={styles.emptyState}>
            <p>Заказы не найдены</p>
          </div>
        ) : (
          <table className={styles.table}>
            <thead>
              <tr className={styles.tableHeader}>
                <th className={styles.th}>ID</th>
                <th className={styles.th}>Точка</th>
                <th className={styles.th}>Дата</th>
                <th className={styles.th}>Статус</th>
                <th className={styles.th}></th>
              </tr>
            </thead>
            <tbody>
              {displayOrders.map((order) => (
                <tr key={order.id} className={styles.tableRow}>
                  <td className={styles.td}>{order.id}</td>
                  <td className={styles.td}>{order.sales_outlet_id || 'Не указана'}</td>
                  <td className={styles.td}>{formatDate(order.created_at)}</td>
                  <td className={styles.td}>
                    <span className={`${styles.statusBadge} ${getStatusClass(order.status_name)}`}>
                      {getStatusText(order.status_name)}
                    </span>
                  </td>
                  <td className={styles.td}>
                    <div className={styles.actions}>
                      <select
                        value={order.status_name ? order.status_name.toLowerCase() : 'ordered'}
                        onChange={(e) => handleStatusChange(order.id, e.target.value)}
                        disabled={statusChanging === order.id}
                        className={styles.statusSelect}
                      >
                        {getStatusOptions().map(option => (
                          <option key={option.value} value={option.value}>
                            {option.label}
                          </option>
                        ))}
                      </select>
                      {statusChanging === order.id && (
                        <span className={styles.loadingText}>...</span>
                      )}
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
};

export default AdminOrders;