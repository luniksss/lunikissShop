import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import AdminOrders from '../components/AdminOrders/AdminOrders';
import AdminCatalog from '../components/AdminCatalog/AdminCatalog';
import AdminWarehouse from '../components/AdminWarehouse/AdminWarehouse';
import AdminEmployees from '../components/AdminEmployees/AdminEmployees';
import AdminOutlets from '../components/AdminOutlets/AdminOutlets';
import styles from './AdminPage.module.css'; 

const AdminPage: React.FC = () => {
  const [activeTab, setActiveTab] = useState<string>('orders');
  const [userRole, setUserRole] = useState<string>('');
  const navigate = useNavigate();

  useEffect(() => {
    const user = localStorage.getItem('user');
    if (!user) {
      navigate('/');
      return;
    }
    
    const userData = JSON.parse(user);
    const role = userData.role;
    
    if (role !== 'admin' && role !== 'seller') {
      navigate('/');
    } else {
      setUserRole(role);
    }
  }, [navigate, activeTab]);

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    navigate('/');
  };

  const renderContent = () => {
    switch (activeTab) {
      case 'catalog':
        return userRole === 'admin' ? <AdminCatalog /> : <AdminOrders />;
      case 'warehouse':
        return <AdminWarehouse />;
      case 'employees':
        return userRole === 'admin' ? <AdminEmployees /> : <AdminOrders />;
      case 'outlets':
        return userRole === 'admin' ? <AdminOutlets /> : <AdminOrders />;
      case 'orders':
        return <AdminOrders />;
    }
  };

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <div className={styles.logoSection}>
          <h1 className={styles.logo}>Lunikiss Admin</h1>
          <p className={styles.subtitle}>
            Панель {userRole === 'admin' ? 'администратора' : 'продавца'}
          </p>
        </div>
        <button onClick={handleLogout} className={styles.logoutButton}>
          Выйти
        </button>
      </header>

      <div className={styles.mainContent}>
        <nav className={styles.sidebar}>
          <ul className={styles.navList}>
            <li className={styles.navItem}>
              <button
                onClick={() => setActiveTab('orders')}
                className={`${styles.navButton} ${activeTab === 'orders' ? styles.activeNavButton : ''}`}
              >
                Заказы
              </button>
            </li>

            <li className={styles.navItem}>
              <button
                onClick={() => setActiveTab('warehouse')}
                className={`${styles.navButton} ${activeTab === 'warehouse' ? styles.activeNavButton : ''}`}
              >
                Склад
              </button>
            </li>
            
            {userRole === 'admin' && (
              <li className={styles.navItem}>
                <button
                  onClick={() => setActiveTab('catalog')}
                  className={`${styles.navButton} ${activeTab === 'catalog' ? styles.activeNavButton : ''}`}
                >
                  Каталог
                </button>
              </li>
            )}
            
            {userRole === 'admin' && (
              <li className={styles.navItem}>
                <button
                  onClick={() => setActiveTab('employees')}
                  className={`${styles.navButton} ${activeTab === 'employees' ? styles.activeNavButton : ''}`}
                >
                  Сотрудники
                </button>
              </li>
            )}
            
            {userRole === 'admin' && (
              <li className={styles.navItem}>
                <button
                  onClick={() => setActiveTab('outlets')}
                  className={`${styles.navButton} ${activeTab === 'outlets' ? styles.activeNavButton : ''}`}
                >
                  Точки
                </button>
              </li>
            )}
          </ul>
        </nav>

        <main className={styles.content}>
          {renderContent()}
        </main>
      </div>
    </div>
  );
};

export default AdminPage;