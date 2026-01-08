import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Header from '../components/Header';
import ProductList from '../components/ProductList/ProductList';
import styles from './HomePage.module.css';

const HomePage: React.FC = () => {
  const navigation = useNavigate();
  const [selectedOutlet, setSelectedOutlet] = useState<string>('1');

  const handleOutletChange = (outletId: string): void => {
    setSelectedOutlet(outletId);
  };

  const handleOrders = (): void => {
    navigation('/orders');
  };

  return (
    <div className={styles.container}>
      <Header
        selectedOutlet={selectedOutlet}
        onOutletChange={handleOutletChange}
        onOrders={handleOrders}
      />
      
      <main className={styles.main}>
        <ProductList outletId={selectedOutlet} />
      </main>
      
      <footer className={styles.footer}>
        <img src={'/icons/footer.jpg'} />
      </footer>
    </div>
  );
};

export default HomePage;