import { cn } from '../../lib/utils';

interface BadgeProps extends React.HTMLAttributes<HTMLSpanElement> {
  variant?: 'default' | 'success' | 'danger' | 'warning' | 'info' | 'secondary';
  size?: 'sm' | 'md';
}

export function Badge({ className, variant = 'default', size = 'md', children, ...props }: BadgeProps) {
  const variants = {
    default: 'bg-background-hover text-content-secondary',
    success: 'bg-accent-success/20 text-accent-success',
    danger: 'bg-accent-danger/20 text-accent-danger',
    warning: 'bg-accent-warning/20 text-accent-warning',
    info: 'bg-accent-primary/20 text-accent-primary',
    secondary: 'bg-accent-secondary/20 text-accent-secondary',
  };

  const sizes = {
    sm: 'px-2 py-0.5 text-xs',
    md: 'px-2.5 py-1 text-xs',
  };

  return (
    <span
      className={cn(
        'inline-flex items-center gap-1 rounded-full font-medium',
        variants[variant],
        sizes[size],
        className
      )}
      {...props}
    >
      {children}
    </span>
  );
}

