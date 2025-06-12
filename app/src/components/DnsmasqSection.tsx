import { Settings, RotateCcw } from 'lucide-react';
import { LoadingState, Messages } from '../types';
import { Message } from './Message';
import { Card, CardHeader, CardTitle, CardContent, Button } from './ui';

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
    <Card>
      <CardHeader>
        <CardTitle>DNS/DHCP Management</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex gap-2">
          <Button onClick={onGenerateConfig} disabled={loading.dnsmasq} variant="outline">
            <Settings className="mr-2 h-4 w-4" />
            {loading.dnsmasq ? 'Generating...' : 'Generate Dnsmasq Config'}
          </Button>
          <Button onClick={onRestart} disabled={loading.restart} variant="outline">
            <RotateCcw className={`mr-2 h-4 w-4 ${loading.restart ? 'animate-spin' : ''}`} />
            {loading.restart ? 'Restarting...' : 'Restart Dnsmasq'}
          </Button>
        </div>
        
        <Message message={messages.dnsmasq} />
        
        {dnsmasqConfig && (
          <pre className="p-4 bg-muted rounded-md text-sm overflow-auto max-h-96">
            {dnsmasqConfig}
          </pre>
        )}
      </CardContent>
    </Card>
  );
};