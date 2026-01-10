import React from 'react';
import styles from './FilterPanel.module.css';

interface FilterPanelProps {
  children: React.ReactNode;
  showClearButton?: boolean;
  onClear?: () => void;
  className?: string;
  clearButtonText?: string;
}

const FilterPanel: React.FC<FilterPanelProps> = ({
  children,
  showClearButton = false,
  onClear,
  className = '',
  clearButtonText = 'Сбросить фильтры'
}) => {
  return (
    <div className={`${styles.filterPanel} ${className}`}>
      <div className={styles.filterControls}>
        {children}
      </div>
      {showClearButton && onClear && (
        <button
          onClick={onClear}
          className={styles.clearButton}
          type="button"
        >
          {clearButtonText}
        </button>
      )}
    </div>
  );
};

export default FilterPanel;