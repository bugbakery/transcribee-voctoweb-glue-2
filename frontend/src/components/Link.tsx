import { navigate } from 'wouter/use-browser-location';
import { cn } from '../cn';

export function Link({
  href,
  children,
  className,
  ...rest
}: { href: string; children: React.ReactNode } & React.HTMLAttributes<HTMLAnchorElement>) {
  return (
    <a
      className={cn('hover:underline', className)}
      onClick={(e) => {
        navigate(href);
        e.preventDefault();
        e.stopPropagation();
      }}
      href={href}
      {...rest}
    >
      {children}
    </a>
  );
}
