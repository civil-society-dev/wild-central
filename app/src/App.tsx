import { useEffect } from 'react';
import './App.css';
import { useApi } from './hooks/useApi';
import {
  SystemStatus,
  ConfigurationSection,
  DnsmasqSection,
  PxeAssetsSection
} from './components';

function App() {
  const {
    status,
    config,
    dnsmasqConfig,
    messages,
    loading,
    showConfigSetup,
    fetchStatus,
    fetchHealth,
    fetchConfig,
    createConfig,
    generateDnsmasqConfig,
    restartDnsmasq,
    downloadAssets
  } = useApi();

  useEffect(() => {
    fetchStatus();
    fetchHealth();
    fetchConfig();
    
    // Refresh status every 30 seconds
    const interval = setInterval(fetchStatus, 30000);
    return () => clearInterval(interval);
  }, [fetchStatus, fetchHealth, fetchConfig]);

  return (
    <div className="App">
      <div className="container">
        <h1>🌩️ Wild Cloud Central Management</h1>
        
        <SystemStatus
          status={status}
          loading={loading}
          messages={messages}
          onRefreshStatus={fetchStatus}
          onCheckHealth={fetchHealth}
        />

        <ConfigurationSection
          config={config}
          loading={loading}
          messages={messages}
          showConfigSetup={showConfigSetup}
          onFetchConfig={fetchConfig}
          onCreateConfig={createConfig}
        />

        <DnsmasqSection
          dnsmasqConfig={dnsmasqConfig}
          loading={loading}
          messages={messages}
          onGenerateConfig={generateDnsmasqConfig}
          onRestart={restartDnsmasq}
        />

        <PxeAssetsSection
          loading={loading}
          messages={messages}
          onDownloadAssets={downloadAssets}
        />
      </div>
    </div>
  );
}

export default App;