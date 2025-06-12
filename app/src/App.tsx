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
      <div className="container mx-auto px-4 py-8 max-w-7xl">
        <header className="mb-12">
          <div className="flex items-center gap-4 mb-3">
            <div className="p-2 bg-primary/10 rounded-lg">
              <Cloud className="h-8 w-8 text-primary" />
            </div>
            <div>
              <h1 className="text-3xl font-bold tracking-tight text-foreground">
                Wild Cloud Central
              </h1>
              <p className="text-muted-foreground text-lg">
                Infrastructure Management Dashboard
              </p>
            </div>
          </div>
        </header>
        
        <div className="grid gap-8 lg:gap-12">
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
    </div>
  );
}

export default App;