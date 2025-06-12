import { useState } from 'react';
import { FileText, Check } from 'lucide-react';
import { Config, LoadingState, Messages } from '../types';
import { parseSimpleYaml } from '../utils/yamlParser';
import { Message } from './Message';
import { Card, CardHeader, CardTitle, CardContent, Button } from './ui';

interface ConfigurationSectionProps {
  config: Config | null;
  loading: LoadingState;
  messages: Messages;
  showConfigSetup: boolean;
  onFetchConfig: () => void;
  onCreateConfig: (config: Config) => void;
}

const defaultConfig = `cloud:
  domain: "wildcloud.local"
  internalDomain: "cluster.local"
  dns:
    ip: "192.168.8.50"
  router:
    ip: "192.168.8.1"
  dhcpRange: "192.168.8.100,192.168.8.200"
  dnsmasq:
    interface: "eth0"

cluster:
  endpointIp: "192.168.8.60"
  nodes:
    talos:
      version: "v1.8.0"

server:
  host: "0.0.0.0"
  port: 5055`;

export const ConfigurationSection = ({
  config,
  loading,
  messages,
  showConfigSetup,
  onFetchConfig,
  onCreateConfig
}: ConfigurationSectionProps) => {
  const [configText, setConfigText] = useState<string>(defaultConfig);

  const handleCreateConfig = () => {
    const configObj = parseSimpleYaml(configText);
    onCreateConfig(configObj);
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Configuration</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <Button onClick={onFetchConfig} disabled={loading.config} variant="outline">
          <FileText className="mr-2 h-4 w-4" />
          {loading.config ? 'Loading...' : 'Load Current Config'}
        </Button>
        
        <Message message={messages.config} />
        
        {showConfigSetup && (
          <div className="space-y-4">
            <div>
              <h3 className="text-lg font-medium">Initial Configuration Setup</h3>
              <p className="text-sm text-muted-foreground">No configuration found. Please provide the initial configuration:</p>
            </div>
            <textarea
              value={configText}
              onChange={(e) => setConfigText(e.target.value)}
              rows={20}
              className="w-full font-mono text-sm border rounded-md p-3 bg-background"
              placeholder="Enter YAML configuration..."
            />
            <Button onClick={handleCreateConfig} disabled={loading.createConfig}>
              <Check className="mr-2 h-4 w-4" />
              {loading.createConfig ? 'Creating...' : 'Create Configuration'}
            </Button>
          </div>
        )}
        
        {config && (
          <pre className="p-4 bg-muted rounded-md text-sm overflow-auto max-h-96">
            {JSON.stringify(config, null, 2)}
          </pre>
        )}
        
        {/* Debug info */}
        <div className="text-xs text-muted-foreground">
          Debug: config={config ? 'exists' : 'null'}, showConfigSetup={showConfigSetup.toString()}
        </div>
      </CardContent>
    </Card>
  );
};