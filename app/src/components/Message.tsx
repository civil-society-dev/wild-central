import { AlertCircle, CheckCircle, Info } from 'lucide-react';
import { Message as MessageType } from '../types';
import { cn } from '@/lib/utils';

interface MessageProps {
  message?: MessageType;
}

export const Message = ({ message }: MessageProps) => {
  if (!message) return null;
  
  const getIcon = () => {
    switch (message.type) {
      case 'error':
        return <AlertCircle className="h-4 w-4" />;
      case 'success':
        return <CheckCircle className="h-4 w-4" />;
      default:
        return <Info className="h-4 w-4" />;
    }
  };

  const getVariantStyles = () => {
    switch (message.type) {
      case 'error':
        return 'border-destructive/40 text-destructive bg-destructive/10 dark:bg-destructive/20';
      case 'success':
        return 'border-emerald-200 text-emerald-800 bg-emerald-50 dark:bg-emerald-950/30 dark:text-emerald-300 dark:border-emerald-700';
      default:
        return 'border-blue-200 text-blue-800 bg-blue-50 dark:bg-blue-950/30 dark:text-blue-300 dark:border-blue-700';
    }
  };
  
  return (
    <div className={cn(
      'flex items-center gap-2 p-3 rounded-md border text-sm',
      getVariantStyles()
    )}>
      {getIcon()}
      <span>{message.message}</span>
    </div>
  );
};