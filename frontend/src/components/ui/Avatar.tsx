import { useState } from 'react';
import { User } from 'lucide-react';
import { cn } from '../../lib/utils';

interface AvatarProps {
  src?: string;
  alt: string;
  size?: 'sm' | 'md' | 'lg' | 'xl';
  className?: string;
}

const sizes = {
  sm: 'w-10 h-10',
  md: 'w-16 h-16',
  lg: 'w-20 h-20',
  xl: 'w-24 h-24',
};

const iconSizes = {
  sm: 'w-5 h-5',
  md: 'w-8 h-8',
  lg: 'w-10 h-10',
  xl: 'w-12 h-12',
};

export function Avatar({ src, alt, size = 'md', className }: AvatarProps) {
  const [hasError, setHasError] = useState(false);

  if (!src || hasError) {
    return (
      <div
        className={cn(
          'rounded-full bg-background-secondary border-2 border-border flex items-center justify-center',
          sizes[size],
          className
        )}
      >
        <User className={cn('text-content-muted', iconSizes[size])} />
      </div>
    );
  }

  return (
    <img
      src={src}
      alt={alt}
      className={cn(
        'rounded-full object-cover border-2 border-border',
        sizes[size],
        className
      )}
      onError={() => setHasError(true)}
    />
  );
}

