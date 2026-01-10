import React, { useEffect, useState } from 'react';
import { OrderItem } from '../../types';
import api from '../../api/api';
import styles from './OrderModal.module.css';

interface OrderModalProps {
  isOpen: boolean;
  onClose: () => void;
  orderId: string;
  onOrderItemDeleted?: (deletedItemId: string) => void;
  refreshOrders?: () => void;
}

const OrderModal: React.FC<OrderModalProps> = ({ isOpen, onClose, orderId, onOrderItemDeleted, refreshOrders }) => {
  const [items, setItems] = useState<OrderItem[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [deletingItem, setDeletingItem] = useState<string | null>(null);
  const [deleteItemError, setDeleteItemError] = useState<string | null>(null);
  const [orderDeleted, setOrderDeleted] = useState<boolean>(false);

  useEffect(() => {
    if (isOpen && orderId && !orderDeleted) {
      fetchOrderItems();
    } else if (!isOpen) {
      resetState();
    }
  }, [isOpen, orderId, orderDeleted]);

  const resetState = () => {
    setItems([]);
    setError(null);
    setDeleteItemError(null);
    setDeletingItem(null);
    setOrderDeleted(false);
  };

  const fetchOrderItems = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const response = await api.getOrderInfo(orderId);
      setItems(response.data);
    } catch (err: any) {
      console.error('Ошибка загрузки деталей заказа:', err);
      if (err.response?.status === 404) {
        setError('Бронь не найдена. Возможно, она была удалена.');
        setOrderDeleted(true);
      } else {
        setError('Не удалось загрузить детали заказа');
      }
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteItem = async (itemId: string, productName: string) => {
    if (!window.confirm(`Вы уверены, что хотите удалить "${productName}" из брони?`)) {
      return;
    }

    setDeletingItem(itemId);
    setDeleteItemError(null);

    try {
      await api.deleteOrderItem(itemId);
      
      const updatedItems = items.filter(item => item.id !== itemId);
      setItems(updatedItems);
      
      if (onOrderItemDeleted) {
        onOrderItemDeleted(itemId);
      }

      if (updatedItems.length === 0) {
        try {
          await api.deleteOrder(orderId);
          setOrderDeleted(true);
          
          if (refreshOrders) {
            await refreshOrders(); 
          }
          
          setTimeout(() => {
            onClose();
            alert('Бронь пуста и будет удалена!');
          }, 100);
          
        } catch (err: any) {
          console.error('Ошибка удаления заказа:', err);
          setDeleteItemError('Не удалось удалить бронь. Пожалуйста, попробуйте позже.');
        }
      } else {
        alert(`Товар "${productName}" успешно удален из брони!`);
      }
    } catch (err: any) {
      console.error('Ошибка удаления товара из заказа:', err);
      setDeleteItemError('Не удалось удалить товар из брони. Пожалуйста, попробуйте позже.');
    } finally {
      setDeletingItem(null);
    }
  };

  if (orderDeleted && !isOpen) {
    return null;
  }

  if (!isOpen) return null;

  return (
    <div className={styles.modalOverlay} onClick={onClose}>
      <div className={styles.modalContent} onClick={(e) => e.stopPropagation()}>
        <div className={styles.modalHeader}>
          <h2 className={styles.modalTitle}>Детали брони №{orderId}</h2>
          <button className={styles.closeButton} onClick={onClose}>×</button>
        </div>
        
        <div className={styles.modalBody}>
          {loading ? (
            <div className={styles.loading}>
              <div className={styles.spinner}></div>
              <p>Загрузка деталей брони...</p>
            </div>
          ) : error ? (
            <div className={styles.error}>
              <p>{error}</p>
              <button 
                onClick={fetchOrderItems}
                className={styles.retryButton}
              >
                Попробовать снова
              </button>
            </div>
          ) : (!items || items.length === 0) ? (
            <div className={styles.empty}>
              <p>В брони нет товаров</p>
              <button 
                onClick={onClose}
                className={styles.retryButton}
              >
                Закрыть
              </button>
            </div>
          ) : (
            <>
              {deleteItemError && (
                <div className={styles.deleteError}>
                  <p>{deleteItemError}</p>
                  <button 
                    onClick={() => setDeleteItemError(null)}
                    className={styles.dismissButton}
                  >
                    ✕
                  </button>
                </div>
              )}

              <div className={styles.itemsTable}>
                <div className={styles.tableHeader}>
                  <div className={styles.tableCell}>Товар</div>
                  <div className={styles.tableCell}>Название</div>
                  <div className={styles.tableCell}>Размер</div>
                  <div className={styles.tableCell}>Количество</div>
                  <div className={styles.tableCell}>Цена</div>
                  <div className={styles.tableCell}>Сумма</div>
                  <div className={styles.tableCell}></div>
                </div>
                
                {items.map((item) => (
                  <div key={item.id} className={styles.tableRow}>
                    <div className={styles.tableCell}>
                      {item.product_image ? (
                        <img 
                          src={item.product_image} 
                          alt={item.product_name}
                          className={styles.productImage}
                        />
                      ) : (
                        <div className={styles.imagePlaceholder}>
                          Нет фото
                        </div>
                      )}
                    </div>
                    <div className={styles.tableCell}>
                      {item.product_name || `Товар #${item.product_id}`}
                    </div>
                    <div className={styles.tableCell}>{item.size}</div>
                    <div className={styles.tableCell}>{item.amount}</div>
                    <div className={styles.tableCell}>{item.price} ₽</div>
                    <div className={styles.tableCell}>{item.price * item.amount} ₽</div>
                    <div className={styles.tableCell}>
                      <button 
                        onClick={() => handleDeleteItem(item.id, item.product_name || `Товар #${item.product_id}`)}
                        disabled={deletingItem === item.id}
                        className={deletingItem === item.id ? styles.deleteItemButtonDisabled : styles.deleteItemButton}
                        title="Удалить товар из брони"
                      >
                        {deletingItem === item.id ? '...' : <img src='/icons/delete.png'/>}
                      </button>
                    </div>
                  </div>
                ))}
              </div>
              
              <div className={styles.totalSection}>
                <div className={styles.totalRow}>
                  <span className={styles.totalLabel}>Итого:</span>
                  <span className={styles.totalValue}>
                    {items.reduce((total, item) => total + (item.price * item.amount), 0)} ₽
                  </span>
                </div>
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  );
};

export default OrderModal;