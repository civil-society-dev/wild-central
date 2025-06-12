import { useState } from 'react';
import { Config, LoadingState, Messages } from '../types';
import { parseSimpleYaml } from '../utils/yamlParser';
import { Message } from './Message';

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
    <div className="section">
      <h2>Configuration</h2>
      <div className="section-content">
        <button onClick={onFetchConfig} disabled={loading.config}>
          {loading.config ? '⏳ Loading...' : '📄 Load Current Config'}
        </button>
        
        <Message message={messages.config} />
        
        {showConfigSetup && (
          <div className="config-setup">
            <h3>Initial Configuration Setup</h3>
            <p>No configuration found. Please provide the initial configuration:</p>
            <textarea
              value={configText}
              onChange={(e) => setConfigText(e.target.value)}
              rows={20}
              cols={80}
              placeholder="Enter YAML configuration..."
            />
            <div className="button-group">
              <button onClick={handleCreateConfig} disabled={loading.createConfig}>
                {loading.createConfig ? '⏳ Creating...' : '✅ Create Configuration'}
              </button>
            </div>
          </div>
        )}
        
        {config && (
          <pre className="config-display">
            {JSON.stringify(config, null, 2)}
          </pre>
        )}
        
        {/* Debug info */}
        <div style={{marginTop: '10px', fontSize: '12px', color: '#666'}}>
          Debug: config={config ? 'exists' : 'null'}, showConfigSetup={showConfigSetup.toString()}
        </div>
      </div>
    </div>
  );
};