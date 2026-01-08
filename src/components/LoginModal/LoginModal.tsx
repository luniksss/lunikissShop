import React, { useState } from 'react';
import api from '../../api/api';
import styles from './LoginModal.module.css';

interface LoginModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
}

interface RegisterData {
  email: string;
  password: string;
  name: string;
  surname: string;
  phone: string;
}

const LoginModal: React.FC<LoginModalProps> = ({ isOpen, onClose, onSuccess }) => {
    const [email, setEmail] = useState<string>('');
    const [password, setPassword] = useState<string>('');
  
    const [name, setName] = useState<string>('');
    const [surname, setSurname] = useState<string>('');
    const [phone, setPhone] = useState<string>('');
  
    const [isRegister, setIsRegister] = useState<boolean>(false);
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<string>('');

    const handleOverlayClick = (e: React.MouseEvent): void => {
        if (e.target === e.currentTarget) {
            handleClose();
        }
    };

    const handleSubmit = async (e: React.FormEvent): Promise<void> => {
        e.preventDefault();
        setLoading(true);
        setError('');

        try {
            let response;
    
            if (isRegister) {
                if (!name || !surname || !email || !password) {
                    setError('Заполните все обязательные поля');
                    setLoading(false);
                    return;
                }

                const registerData: RegisterData = {
                    email,
                    password,
                    name,
                    surname,
                    phone
                };
            
                response = await api.register(registerData);
            } else {
                if (!email || !password) {
                    setError('Заполните email и пароль');
                    setLoading(false);
                    return;
                }
            
                response = await api.login(email, password);
            }
        
            if (response.data.access_token) {
                localStorage.setItem('token', response.data.access_token);
                
                if (response.data.user) {
                    localStorage.setItem('user', JSON.stringify(response.data.user));
                }
                
                onSuccess(); 
                onClose();
                handleCloseForm();
            } else {
                setError('Ошибка: токен не получен');
            }
        } catch (error: any) {
            console.error('Ошибка авторизации:', error);
            setError(
                error.response?.data?.message || 
                (isRegister ? 'Ошибка регистрации' : 'Ошибка входа') + 
                '. Проверьте данные и попробуйте снова.'
            );
            window.location.href = '/';
        } finally {
            setLoading(false);
        }
};

  const handleClose = (): void => {
    handleCloseForm();
    onClose();
  };

  const handleCloseForm = (): void => {
    setEmail('');
    setPassword('');
    setName('');
    setSurname('');
    setPhone('');
    setError('');
    setIsRegister(false);
  };

  const switchMode = (): void => {
    setIsRegister(!isRegister);
    setError('');
  };

  if (!isOpen) return null;

  return (
    <div className={styles.overlay}
    onClick={handleOverlayClick}>
      <div className={styles.modal}>
        <div className={styles.header}>
          <h2 className={styles.title}>
            {isRegister ? 'Регистрация' : 'Вход в систему'}
          </h2>
          <button onClick={handleClose} className={styles.closeButton}>
            ×
          </button>
        </div>

        <form onSubmit={handleSubmit} className={styles.form}>
          <div className={styles.formGroup}>
            <label htmlFor="email" className={styles.label}>
              Email <span className={styles.required}>*</span>
            </label>
            <input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="Введите email"
              required
              className={styles.input}
              disabled={loading}
            />
          </div>

          <div className={styles.formGroup}>
            <label htmlFor="password" className={styles.label}>
              Пароль <span className={styles.required}>*</span>
            </label>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Введите пароль"
              required
              className={styles.input}
              disabled={loading}
              minLength={6}
            />
          </div>

          {isRegister && (
            <>
              <div className={styles.twoColumns}>
                <div className={styles.column}>
                  <label htmlFor="name" className={styles.label}>
                    Имя <span className={styles.required}>*</span>
                  </label>
                  <input
                    id="name"
                    type="text"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    placeholder="Имя"
                    required={isRegister}
                    className={styles.input}
                    disabled={loading}
                  />
                </div>
                
                <div className={styles.column}>
                  <label htmlFor="surname" className={styles.label}>
                    Фамилия <span className={styles.required}>*</span>
                  </label>
                  <input
                    id="surname"
                    type="text"
                    value={surname}
                    onChange={(e) => setSurname(e.target.value)}
                    placeholder="Фамилия"
                    required={isRegister}
                    className={styles.input}
                    disabled={loading}
                  />
                </div>
              </div>

              <div className={styles.formGroup}>
                <label htmlFor="phone" className={styles.label}>
                  Телефон
                </label>
                <input
                  id="phone"
                  type="tel"
                  value={phone}
                  onChange={(e) => setPhone(e.target.value)}
                  placeholder="+7 (999) 999-99-99"
                  className={styles.input}
                  disabled={loading}
                />
              </div>
            </>
          )}

          {error && (
            <div className={styles.error}>
              {error}
            </div>
          )}

          <div className={styles.buttons}>
            <button
              type="submit"
              disabled={loading}
              className={styles.submitButton}
            >
              {loading ? 'Загрузка...' : (isRegister ? 'Зарегистрироваться' : 'Войти')}
            </button>

            <button
              type="button"
              onClick={switchMode}
              className={styles.switchButton}
              disabled={loading}
            >
              {isRegister ? 'Уже есть аккаунт? Войти' : 'Нет аккаунта? Зарегистрироваться'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default LoginModal;