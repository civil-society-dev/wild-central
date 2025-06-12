import { useState } from 'react';
import { FileText, Check, AlertCircle } from 'lucide-react';
import { useConfig } from '../hooks';
import { parseSimpleYaml } from '../utils/yamlParser';
import { Card, CardHeader, CardTitle, CardContent, Button } from './ui';

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

export const ConfigSectionNew = () => {
  const [configText, setConfigText] = useState<string>(defaultConfig);
  const { 
    config, 
    isConfigured, 
    showConfigSetup, 
    isLoading, 
    isCreating, 
    error, 
    createConfig, 
    refetch 
  } = useConfig();

  const handleCreateConfig = () => {
    try {
      const configObj = parseSimpleYaml(configText);
      createConfig(configObj);
    } catch (err) {
      console.error('Failed to parse YAML:', err);
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Configuration (React Query)</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <Button onClick={() => refetch()} disabled={isLoading} variant="outline">
          <FileText className="mr-2 h-4 w-4" />
          {isLoading ? 'Loading...' : 'Reload Configuration'}
        </Button>
        
        {error && (
          <div className="p-3 bg-red-50 dark:bg-red-950 border border-red-200 dark:border-red-800 rounded-md flex items-start gap-2">
            <AlertCircle className="h-4 w-4 text-red-600 mt-0.5" />
            <div>
              <p className="text-sm font-medium text-red-800 dark:text-red-200">Configuration Error</p>
              <p className="text-sm text-red-700 dark:text-red-300">{error.message}</p>
            </div>
          </div>
        )}

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
            <Button onClick={handleCreateConfig} disabled={isCreating}>
              <Check className="mr-2 h-4 w-4" />
              {isCreating ? 'Creating...' : 'Create Configuration'}
            </Button>
          </div>
        )}
        
        {config && isConfigured && (
          <div className="space-y-2">
            <div className="p-3 bg-green-50 dark:bg-green-950 border border-green-200 dark:border-green-800 rounded-md">
              <p className="text-sm text-green-800 dark:text-green-200">
                ✓ Configuration loaded successfully
              </p>
            </div>
            <pre className="p-4 bg-muted rounded-md text-sm overflow-auto max-h-96">
              {JSON.stringify(config, null, 2)}
            </pre>
          </div>
        )}
        
        <div className="text-xs text-muted-foreground">
          React Query Status: isLoading={isLoading.toString()}, isConfigured={isConfigured.toString()}, showSetup={showConfigSetup.toString()}
        </div>
      </CardContent>
    </Card>
  );
};