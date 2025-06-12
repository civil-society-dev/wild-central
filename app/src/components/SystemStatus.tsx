import { RefreshCw, Activity } from 'lucide-react';
import { Status, LoadingState, Messages } from '../types';
import { formatTimestamp } from '../utils/formatters';
import { Message } from './Message';
import { Card, CardHeader, CardTitle, CardContent, Button, Badge } from './ui';

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
    <Card>
      <CardHeader>
        <CardTitle>System Status</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex gap-2">
          <Button onClick={onRefreshStatus} disabled={loading.status} variant="outline">
            <RefreshCw className={`mr-2 h-4 w-4 ${loading.status ? 'animate-spin' : ''}`} />
            {loading.status ? 'Checking...' : 'Refresh Status'}
          </Button>
          <Button onClick={onCheckHealth} disabled={loading.health} variant="outline">
            <Activity className="mr-2 h-4 w-4" />
            {loading.health ? 'Checking...' : 'Check Health'}
          </Button>
        </div>
        
        <Message message={messages.health} />
        
        {status && (
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
            <div className="space-y-2">
              <p className="text-sm font-medium text-muted-foreground">Status</p>
              <Badge 
                variant={status.status === 'running' ? 'default' : 'destructive'} 
                className={`text-xs font-medium ${status.status === 'running' ? 'bg-emerald-600 hover:bg-emerald-700' : ''}`}
              >
                <div className={`w-2 h-2 rounded-full mr-2 ${status.status === 'running' ? 'bg-emerald-200' : 'bg-red-200'}`} />
                {status.status}
              </Badge>
            </div>
            <div className="space-y-2">
              <p className="text-sm font-medium text-muted-foreground">Version</p>
              <p className="text-sm">{status.version}</p>
            </div>
            <div className="space-y-2">
              <p className="text-sm font-medium text-muted-foreground">Uptime</p>
              <p className="text-sm">{status.uptime}</p>
            </div>
            <div className="space-y-2">
              <p className="text-sm font-medium text-muted-foreground">Last Updated</p>
              <p className="text-sm">{formatTimestamp(status.timestamp)}</p>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
};