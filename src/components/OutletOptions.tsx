import React from 'react';
import { SalesOutlet } from '../types';

interface OutletOptionsProps {
  outlets: SalesOutlet[];
}

const OutletOptions: React.FC<OutletOptionsProps> = ({ outlets }) => {
  return (
    <>
      <option value="">Выберите точку</option>
      {outlets.map((outlet: SalesOutlet) => (
        <option key={outlet.id} value={outlet.id}>
          {outlet.address}
        </option>
      ))}
    </>
  );
};

export default OutletOptions;