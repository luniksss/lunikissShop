import React from 'react';
import styles from './RoleFilter.module.css';

interface RoleOption {
  value: string;
  label: string;
}

interface RoleFilterProps {
  value: string;
  onChange: (value: string) => void;
  options: RoleOption[];
  className?: string;
}

const RoleFilter: React.FC<RoleFilterProps> = ({
  value,
  onChange,
  options,
  className = ''
}) => {
  const handleChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    onChange(e.target.value);
  };

  return (
    <div className={`${styles.filterContainer} ${className}`}>
      <select
        value={value}
        onChange={handleChange}
        className={styles.filterSelect}
      >
        {options.map((option) => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>
      <div className={styles.selectArrow}>▼</div>
    </div>
  );
};

export default RoleFilter;