import { Message as MessageType } from '../types';

interface MessageProps {
  message?: MessageType;
}

export const Message = ({ message }: MessageProps) => {
  if (!message) return null;
  
  return (
    <div className={`message ${message.type}`}>
      {message.message}
    </div>
  );
};