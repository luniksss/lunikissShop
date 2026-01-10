import React from 'react';
import styles from './SearchBar.module.css';

interface SearchBarProps {
    value: string;
    onChange: (value: string) => void;
    placeholder?: string;
    className?: string;
}

const SearchBar: React.FC<SearchBarProps> = ({
    value,
    onChange,
    placeholder = 'Поиск...',
    className = ''
}) => {
    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        onChange(e.target.value);
    };

    const handleClear = () => {
        onChange('');
    };

    return (
        <div className={`${styles.searchContainer} ${className}`}>
            <div className={styles.searchWrapper}>
                <input
                    type="text"
                    placeholder={placeholder}
                    value={value}
                    onChange={handleChange}
                    className={styles.searchInput}
                />
                <button
                    onClick={handleClear}
                    className={styles.clearButton}
                    type="button"
                    aria-label="Очистить поиск"
                >
                    ✕
                </button>
            </div>
        </div>
    );
};

export default SearchBar;