import { Link as RRLink } from 'react-router';
import { cn } from '../cn';

export function Link({
  to,
  children,
  className,
  ...rest
}: { to: string; children: React.ReactNode } & React.HTMLAttributes<HTMLAnchorElement>) {
  return (
    <RRLink className={cn('hover:underline', className)} to={to} {...rest}>
      {children}
    </RRLink>
  );
}
