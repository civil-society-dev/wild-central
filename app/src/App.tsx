import { useEffect } from 'react';
import { Cloud } from 'lucide-react';
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
    <div className="min-h-screen bg-background">
      <div className="container mx-auto p-6 space-y-6">
        <div className="flex items-center gap-3 mb-8">
          <Cloud className="h-8 w-8 text-primary" />
          <h1 className="text-3xl font-bold">Wild Cloud Central Management</h1>
        </div>
        
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