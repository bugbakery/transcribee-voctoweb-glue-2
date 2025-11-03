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
        'bg-white/80 border border-white text-black hover:bg-white text-sm font-semibold py-1 px-2 rounded-lg',
        // 'bg-white/20 border border-white/15 text-white hover:bg-white text-sm font-semibold py-1 px-2 rounded-lg',
        className,
      )}
      {...props}
    >
      {children}
    </button>
  );
}
