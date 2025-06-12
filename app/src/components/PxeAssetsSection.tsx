import { LoadingState, Messages } from '../types';
import { Message } from './Message';

interface PxeAssetsSectionProps {
  loading: LoadingState;
  messages: Messages;
  onDownloadAssets: () => void;
}

export const PxeAssetsSection = ({
  loading,
  messages,
  onDownloadAssets
}: PxeAssetsSectionProps) => {
  return (
    <div className="section">
      <h2>PXE Boot Assets</h2>
      <div className="section-content">
        <button onClick={onDownloadAssets} disabled={loading.assets}>
          {loading.assets ? '⏳ Downloading...' : '⬇️ Download/Update PXE Assets'}
        </button>
        
        <Message message={messages.assets} />
      </div>
    </div>
  );
};