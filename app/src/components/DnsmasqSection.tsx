import { LoadingState, Messages } from '../types';
import { Message } from './Message';

interface DnsmasqSectionProps {
  dnsmasqConfig: string;
  loading: LoadingState;
  messages: Messages;
  onGenerateConfig: () => void;
  onRestart: () => void;
}

export const DnsmasqSection = ({
  dnsmasqConfig,
  loading,
  messages,
  onGenerateConfig,
  onRestart
}: DnsmasqSectionProps) => {
  return (
    <div className="section">
      <h2>DNS/DHCP Management</h2>
      <div className="section-content">
        <div className="button-group">
          <button onClick={onGenerateConfig} disabled={loading.dnsmasq}>
            {loading.dnsmasq ? '⏳ Generating...' : '⚙️ Generate Dnsmasq Config'}
          </button>
          <button onClick={onRestart} disabled={loading.restart}>
            {loading.restart ? '⏳ Restarting...' : '🔄 Restart Dnsmasq'}
          </button>
        </div>
        
        <Message message={messages.dnsmasq} />
        
        {dnsmasqConfig && (
          <pre className="config-display">
            {dnsmasqConfig}
          </pre>
        )}
      </div>
    </div>
  );
};