import { useId } from 'react';
import clsx from 'clsx';

import { FiTerminal } from 'react-icons/fi';
import { HiOutlineLightBulb } from 'react-icons/hi';
import {IoRocketOutline} from "react-icons/io5";

const icons = {
  rocket: IoRocketOutline,
  terminal: FiTerminal,
  bulb: HiOutlineLightBulb,
};

export function Icon({ icon, className, ...props }) {
  let id = useId();
  let IconComponent = icons[icon];

  return (
    <IconComponent
      id={id}
      aria-hidden="true"
      fill="none"
      className={className}
      {...props}
    />
  );
}

const gradients = {
  blue: [
    { stopColor: '#0EA5E9' },
    { stopColor: '#22D3EE', offset: '.527' },
    { stopColor: '#818CF8', offset: 1 },
  ],
  amber: [
    { stopColor: '#FDE68A', offset: '.08' },
    { stopColor: '#F59E0B', offset: '.837' },
  ],
};

export function Gradient({ color = 'blue', ...props }) {
  return (
    <radialGradient cx={0} cy={0} r={1} gradientUnits="userSpaceOnUse" {...props}>
      {gradients[color].map((stop, stopIndex) => (
        <stop key={stopIndex} {...stop} />
      ))}
    </radialGradient>
  );
}

export function LightMode({ className, ...props }) {
  return <g className={clsx('dark:hidden', className)} {...props} />;
}

export function DarkMode({ className, ...props }) {
  return <g className={clsx('hidden dark:inline', className)} {...props} />;
}
