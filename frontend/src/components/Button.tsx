import { cn } from '../cn';

export function Button({
  children,
  className,
  ...props
}: {
  children: React.ReactNode;
  className?: string;
} & React.ButtonHTMLAttributes<HTMLButtonElement>) {
  return (
    <button
      className={cn(
        'bg-white/90 border border-white text-black hover:bg-white text-sm font-semibold py-1 px-2 rounded-lg',
        className,
      )}
      {...props}
    >
      {children}
    </button>
  );
}
