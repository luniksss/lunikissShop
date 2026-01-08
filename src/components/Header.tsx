import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../api/api';
import { SalesOutlet } from '../types';
import LoginModal from './LoginModal/LoginModal';
import UserProfileModal from './UserProfileModal/UserProfileModal';
import OutletOptions from './OutletOptions';
import styles from './Header.module.css'

interface HeaderProps {
  selectedOutlet: string;
  onOutletChange: (outletId: string) => void;
  onOrders: () => void;
}

const Header: React.FC<HeaderProps> = ({
  selectedOutlet,
  onOutletChange,
  onOrders,
}) => {
  const navigation = useNavigate();
  const [outlets, setOutlets] = useState<SalesOutlet[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [isLoggedIn, setIsLoggedIn] = useState<boolean>(false);
  const [isLoginModalOpen, setIsLoginModalOpen] = useState<boolean>(false);   
  const [isProfileModalOpen, setIsProfileModalOpen] = useState<boolean>(false);
  const [userRole, setUserRole] = useState<string | null>(null);

  useEffect(() => {
    const token = localStorage.getItem('token');
    const user = localStorage.getItem('user');
    setIsLoggedIn(!!token);
  
    if (user) {
      const userData = JSON.parse(user);
      setUserRole(userData.role);
    }

    loadOutlets();
  }, []);

  const loadOutlets = async (): Promise<void> => {
    try {
      const response = await api.getSalesOutlets();
      setOutlets(response.data);
    } catch (error) {
      console.error('Ошибка загрузки точек:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleLogin = (): void => {
    setIsLoginModalOpen(true);
  };

  const handleLoginSuccess = (): void => {
    const token = localStorage.getItem('token');
    setIsLoggedIn(!!token);

    const user = localStorage.getItem('user');
    if (user) {
        const userData = JSON.parse(user);
        setUserRole(userData.role);
    }
  };

  const handleLogout = (): void => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    setIsLoggedIn(false);
    setUserRole(null);
    navigation('/');
  };

  const handleOrdersClick = (): void => {
    if (isLoggedIn) {
      navigation('/orders');
    } else {
      handleLogin();
    }
  };

   const handleProfileClick = (): void => {
    if (isLoggedIn) {
      setIsProfileModalOpen(true);
    } else {
      setIsLoginModalOpen(true);
    }
  };

  return (
    <>
      <header className={styles.header}>
        <div className={styles.logo}>
          <h1 className={styles.logoText}>Lunikiss</h1>
        </div>
        
        <div className={styles.controls}>
          <div className={styles.outletSelector}>
            <select
              id="outlet-select"
              value={selectedOutlet}
              onChange={(e: React.ChangeEvent<HTMLSelectElement>) => 
                onOutletChange(e.target.value)
              }
              className={styles.select}
              disabled={loading}
            >
              <OutletOptions outlets={outlets} />
            </select>
          </div>

           <div className={styles.buttons}>
            <button 
              onClick={handleOrdersClick} 
              className={styles.iconButton}
              title="Мои заказы"
            >
              <img className={styles.iconImg} src='/icons/user-orders.png'/>
            </button>

            <button 
              onClick={handleProfileClick}
              className={styles.iconButton}
              title={isLoggedIn ? "Мой профиль" : "Войти"}
            >
              <img className={styles.iconImg} src='/icons/user-account.png'/>
              {isLoggedIn && (
                <div className={styles.userIndicator}></div>
              )}
            </button>

            {(userRole === 'admin' || userRole === 'seller') && (
              <button 
                onClick={() => navigation('/admin')}
                className={styles.iconButton}
                title="Админка"
              >
                <img className={styles.iconImg} src='/icons/admin.png'/>
              </button>
            )}
            
            {isLoggedIn ? (
              <button 
                onClick={handleLogout} 
                className={`${styles.button} ${styles.logoutButton}`}
              >
                Выйти
              </button>
            ) : (
              <button 
                onClick={handleLogin} 
                className={styles.button}
              >
                Войти
              </button>
            )}
          </div>
        </div>
      </header>

        <LoginModal
        isOpen={isLoginModalOpen}
        onClose={() => {
            setIsLoginModalOpen(false);
        }}
        onSuccess={() => {
            handleLoginSuccess();
        }}
        />

        <UserProfileModal
        isOpen={isProfileModalOpen}
        onClose={() => {
          setIsProfileModalOpen(false);
        }}
      />
    </>
  );
};

export default Header;