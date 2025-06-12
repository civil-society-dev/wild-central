import { Status, LoadingState, Messages } from '../types';
import { formatTimestamp } from '../utils/formatters';
import { Message } from './Message';

interface SystemStatusProps {
  status: Status | null;
  loading: LoadingState;
  messages: Messages;
  onRefreshStatus: () => void;
  onCheckHealth: () => void;
}

export const SystemStatus = ({
  status,
  loading,
  messages,
  onRefreshStatus,
  onCheckHealth
}: SystemStatusProps) => {
  return (
    <div className="section">
      <h2>System Status</h2>
      <div className="section-content">
        <div className="button-group">
          <button onClick={onRefreshStatus} disabled={loading.status}>
            {loading.status ? '⏳ Checking...' : '🔄 Refresh Status'}
          </button>
          <button onClick={onCheckHealth} disabled={loading.health}>
            {loading.health ? '⏳ Checking...' : '🏥 Check Health'}
          </button>
        </div>
        
        <Message message={messages.health} />
        
        {status && (
          <div className="status-display">
            <div className="status-grid">
              <div className="status-item">
                <span className="status-label">Status:</span>
                <span className={`status-value ${status.status === 'running' ? 'running' : 'stopped'}`}>
                  {status.status}
                </span>
              </div>
              <div className="status-item">
                <span className="status-label">Version:</span>
                <span className="status-value">{status.version}</span>
              </div>
              <div className="status-item">
                <span className="status-label">Uptime:</span>
                <span className="status-value">{status.uptime}</span>
              </div>
              <div className="status-item">
                <span className="status-label">Last Updated:</span>
                <span className="status-value">{formatTimestamp(status.timestamp)}</span>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};