import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../../api/api';
import { Order } from '../../types';
import { useOutlets } from '../../hooks/useOutlets';
import OrderModal from '../OrderModal/OrderModal';
import styles from './OrderList.module.css';

const OrderList: React.FC = () => {
  const [orders, setOrders] = useState<Order[]>([]);
  const [ordersLoading, setOrdersLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedOrderId, setSelectedOrderId] = useState<string | null>(null);
  const [selectedOrderStatus, setSelectedOrderStatus] = useState<string | null>(null);
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
  const [deleting, setDeleting] = useState<string | null>(null);
  const [deleteError, setDeleteError] = useState<string | null>(null);
  
  const navigate = useNavigate();
  const { getOutletAddress, loading: outletsLoading } = useOutlets();

  useEffect(() => {
    fetchOrders();
  }, [navigate]);

  const fetchOrders = async () => {
    try {
      setOrdersLoading(true);
      setError(null);
      setDeleteError(null);
      
      const token = localStorage.getItem('token');
      const user = localStorage.getItem('user');
      
      if (!token || !user) {
        navigate('/');
        return;
      }

      const userData = JSON.parse(user);
      const response = await api.getUserOrders(userData.id);
      setOrders(Array.isArray(response.data) ? response.data : []);
    } catch (err: any) {
      console.error('Ошибка загрузки заказов:', err);
      setError('Не удалось загрузить заказы. Пожалуйста, попробуйте позже.');
    } finally {
      setOrdersLoading(false);
    }
  };

  const handleDeleteOrder = async (orderId: string, orderNumber: string) => {
    if (!window.confirm(`Вы уверены, что хотите удалить бронь №${orderNumber}? Это действие нельзя отменить.`)) {
      return;
    }

    setDeleting(orderId);
    setDeleteError(null);

    try {
      await api.deleteOrder(orderId);
    
      setOrders(prevOrders => prevOrders.filter(order => order.id !== orderId));
      if (selectedOrderId === orderId) {
        closeModal();
      }
      
      alert(`Бронь №${orderNumber} успешно удалена!`);
    } catch (err: any) {
      console.error('Ошибка удаления заказа:', err);
      setDeleteError('Не удалось удалить бронь. Пожалуйста, попробуйте позже.');
    } finally {
      setDeleting(null);
    }
  };

  const handleOrderItemDeleted = (deletedItemId: string) => {
    fetchOrders();
  };

  const handleOrderClick = (orderId: string, orderStatus: string) => {
    setSelectedOrderId(orderId);
    setSelectedOrderStatus(orderStatus);
    setIsModalOpen(true);
  };

  const closeModal = () => {
    setIsModalOpen(false);
    setSelectedOrderId(null);
    setSelectedOrderStatus(null);
  };

  const loading = ordersLoading || outletsLoading;

  const formatDate = (dateString: string): string => {
    const date = new Date(dateString);
    return date.toLocaleDateString('ru-RU', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      timeZone: 'UTC'
    });
  };

  const getStatusClass = (status: string) => {
    switch(status) {
      case 'ordered': return styles.statusOrdered;
      case 'cancelled': return styles.statusCancelled;
      case 'delivered': return styles.statusDelivered;
      default: return styles.statusNew;
    }
  };

  const canDeleteOrder = (status: string): boolean => {
    return status.toLowerCase() === 'ordered';
  };

  if (loading) {
    return (
      <div className={styles.loadingContainer}>
        <div className={styles.loadingSpinner}></div>
        <p>Загрузка данных...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className={styles.errorContainer}>
        <p className={styles.errorText}>{error}</p>
        <button 
          onClick={() => navigate('/')}
          className={styles.backButton}
        >
          Вернуться на главную
        </button>
      </div>
    );
  }

  if (orders.length === 0) {
    return (
      <div className={styles.emptyContainer}>
        <h2 className={styles.emptyTitle}>У вас пока нет заказов</h2>
        <p className={styles.emptyText}>Совершите свой первый заказ!</p>
        <button 
          onClick={() => navigate('/')}
          className={styles.shopButton}
        >
          Перейти к покупкам
        </button>
      </div>
    );
  }

  return (
    <>
    <div className={styles.container}>
        <header className={styles.header}>
        <button 
            onClick={() => navigate('/')}
            className={styles.backButton}
        >
            ← Назад
        </button>
        <h1 className={styles.title}>Мои брони</h1>
        </header>

        {deleteError && (
          <div className={styles.deleteError}>
            <p className={styles.deleteErrorText}>{deleteError}</p>
            <button 
              onClick={() => setDeleteError(null)}
              className={styles.dismissButton}
            >
              ✕
            </button>
          </div>
        )}

        <main className={styles.main}>
        <div className={styles.ordersGrid}>
            {orders.map((order) => (
            <div key={order.id} className={styles.orderCard}>
                <div className={styles.orderHeader}>
                <div className={styles.orderInfo}>
                    <h3 className={styles.orderNumber}>Бронь №{order.id}</h3>
                </div>
                <div className={styles.orderStatus}>
                    <span className={`${styles.statusBadge} ${getStatusClass(order.status_name)}`}>
                    {getStatusText(order.status_name)}
                    </span>
                </div>
                </div>
                
                <div className={styles.orderDetails}>
                <div className={styles.detailItem}>
                    <span className={styles.detailLabel}>Точка продаж:</span>
                    <span className={styles.detailValue}>
                    {getOutletAddress(order.sales_outlet_id)}
                    </span>
                </div>
                <div className={styles.detailItem}>
                    <span className={styles.detailLabel}>Дата создания:</span>
                    <span className={styles.detailValue}>
                    {formatDate(order.created_at)}
                    </span>
                </div>
                </div>

                <div className={styles.orderFooter}>
                  <button 
                    onClick={() => handleOrderClick(order.id, order.status_name)}
                    className={styles.detailsButton}
                  >
                    Подробнее
                  </button>

                  {canDeleteOrder(order.status_name) && (
                      <button 
                        onClick={() => handleDeleteOrder(order.id, order.id)}
                        disabled={deleting === order.id}
                        className={deleting === order.id ? styles.deleteButtonDisabled : styles.deleteButton}
                        title="Удалить бронь"
                      >
                        {deleting === order.id ? (
                          <span className={styles.deleteButtonText}>Удаление...</span>
                        ) : (
                          <img src='/icons/delete.png' alt="Удалить"/>
                        )}
                      </button>
                    )}
                </div>
            </div>
            ))}
        </div>
        </main>
    </div>

    {selectedOrderId && (
        <OrderModal
          isOpen={isModalOpen}
          onClose={closeModal}
          orderId={selectedOrderId}
          orderStatus={selectedOrderStatus}
          onOrderItemDeleted={handleOrderItemDeleted}
          refreshOrders={fetchOrders}
        />
      )}
    </>
    );
};

const getStatusText = (status: string): string => {
  const statusMap: Record<string, string> = {
    'ordered': 'Забронировано',
    'delivered': 'Доставлено',
    'cancelled': 'Отменено',
  };
  return statusMap[status] || status;
};

export default OrderList;