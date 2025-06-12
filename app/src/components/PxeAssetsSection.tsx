import { Download } from 'lucide-react';
import { LoadingState, Messages } from '../types';
import { Message } from './Message';
import { Card, CardHeader, CardTitle, CardContent, Button } from './ui';

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
    <Card>
      <CardHeader>
        <CardTitle>PXE Boot Assets</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <Button onClick={onDownloadAssets} disabled={loading.assets} variant="outline">
          <Download className="mr-2 h-4 w-4" />
          {loading.assets ? 'Downloading...' : 'Download/Update PXE Assets'}
        </Button>
        
        <Message message={messages.assets} />
      </CardContent>
    </Card>
  );
};