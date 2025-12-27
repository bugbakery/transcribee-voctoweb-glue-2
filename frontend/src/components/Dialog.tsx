import { Dialog as HeadlessDialog, DialogPanel, DialogTitle } from '@headlessui/react';

export function Dialog(props: {
  isOpen: boolean;
  onClose: () => void;
  children: React.ReactNode;
  title: string;
}) {
  return (
    <HeadlessDialog open={props.isOpen} onClose={() => props.onClose()} className="relative z-50">
      <div className="fixed inset-0 flex w-screen items-center justify-center p-4 bg-black/90">
        <DialogPanel className="max-w-lg space-y-4 border bg-black p-12 rounded-2xl">
          <DialogTitle className="font-bold">{props.title}</DialogTitle>
          {props.children}
        </DialogPanel>
      </div>
    </HeadlessDialog>
  );
}
