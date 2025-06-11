import React, { useState, useEffect } from 'react';
import './App.css';

function App() {
  const [status, setStatus] = useState(null);
  const [health, setHealth] = useState(null);
  const [config, setConfig] = useState(null);
  const [configText, setConfigText] = useState('');
  const [dnsmasqConfig, setDnsmasqConfig] = useState('');
  const [messages, setMessages] = useState({});
  const [loading, setLoading] = useState({ status: true });
  const [showConfigSetup, setShowConfigSetup] = useState(false);

  const API_BASE = 'http://localhost:5055';

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

  useEffect(() => {
    fetchStatus();
    fetchHealth();
    fetchConfig();
    
    // Refresh status every 30 seconds
    const interval = setInterval(fetchStatus, 30000);
    return () => clearInterval(interval);
  }, []);

  const setLoadingState = (key, value) => {
    setLoading(prev => ({ ...prev, [key]: value }));
  };

  const setMessage = (key, message, type = 'info') => {
    setMessages(prev => ({ ...prev, [key]: { message, type } }));
  };

  const fetchStatus = async () => {
    try {
      setLoadingState('status', true);
      const response = await fetch(`${API_BASE}/api/status`);
      if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
      const data = await response.json();
      setStatus(data);
      setMessage('status', null);
    } catch (err) {
      setMessage('status', `Failed to fetch status: ${err.message}`, 'error');
    } finally {
      setLoadingState('status', false);
    }
  };

  const fetchHealth = async () => {
    try {
      setLoadingState('health', true);
      const response = await fetch(`${API_BASE}/api/v1/health`);
      if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
      const data = await response.json();
      setHealth(data);
      setMessage('health', `Service: ${data.service} - Status: ${data.status}`, 'success');
    } catch (err) {
      setMessage('health', `Health check failed: ${err.message}`, 'error');
    } finally {
      setLoadingState('health', false);
    }
  };

  const fetchConfig = async () => {
    try {
      setLoadingState('config', true);
      const response = await fetch(`${API_BASE}/api/v1/config`);
      if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
      const data = await response.json();
      
      console.log('fetchConfig response:', data);
      console.log('data.configured:', data.configured);
      console.log('data.config:', data.config);
      
      if (data.configured === false) {
        console.log('Config not configured, showing setup');
        setShowConfigSetup(true);
        setConfigText(defaultConfig);
        setMessage('config', data.message, 'error');
        setConfig(null);
      } else {
        console.log('Config is configured, setting config data');
        setShowConfigSetup(false);
        setConfig(data.config);
        setMessage('config', 'Configuration loaded successfully', 'success');
      }
    } catch (err) {
      console.error('fetchConfig error:', err);
      setMessage('config', `Failed to load config: ${err.message}`, 'error');
    } finally {
      setLoadingState('config', false);
    }
  };

  const createConfig = async () => {
    try {
      setLoadingState('createConfig', true);
      const configObj = parseSimpleYaml(configText);
      
      const response = await fetch(`${API_BASE}/api/v1/config`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(configObj)
      });

      if (response.ok) {
        const data = await response.json();
        setMessage('config', `Configuration created successfully! Status: ${data.status}`, 'success');
        setShowConfigSetup(false);
        fetchConfig();
      } else {
        const errorText = await response.text();
        setMessage('config', `Failed to create configuration: ${errorText}`, 'error');
      }
    } catch (err) {
      setMessage('config', `Error creating config: ${err.message}`, 'error');
    } finally {
      setLoadingState('createConfig', false);
    }
  };

  const generateDnsmasqConfig = async () => {
    try {
      setLoadingState('dnsmasq', true);
      const response = await fetch(`${API_BASE}/api/v1/dnsmasq/config`);
      if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
      const data = await response.text();
      setDnsmasqConfig(data);
      setMessage('dnsmasq', 'Dnsmasq config generated successfully', 'success');
    } catch (err) {
      setMessage('dnsmasq', `Failed to generate dnsmasq config: ${err.message}`, 'error');
    } finally {
      setLoadingState('dnsmasq', false);
    }
  };

  const restartDnsmasq = async () => {
    try {
      setLoadingState('restart', true);
      const response = await fetch(`${API_BASE}/api/v1/dnsmasq/restart`, { method: 'POST' });
      if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
      const data = await response.json();
      setMessage('dnsmasq', `Dnsmasq restart: ${data.status}`, 'success');
    } catch (err) {
      setMessage('dnsmasq', `Failed to restart dnsmasq: ${err.message}`, 'error');
    } finally {
      setLoadingState('restart', false);
    }
  };

  const downloadAssets = async () => {
    try {
      setLoadingState('assets', true);
      const response = await fetch(`${API_BASE}/api/v1/pxe/assets`, { method: 'POST' });
      if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
      const data = await response.json();
      setMessage('assets', `PXE Assets: ${data.status}`, 'success');
    } catch (err) {
      setMessage('assets', `Failed to download assets: ${err.message}`, 'error');
    } finally {
      setLoadingState('assets', false);
    }
  };

  // Simple YAML to JSON parser for basic configuration
  const parseSimpleYaml = (yamlText) => {
    const config = {
      cloud: { dns: {}, router: {}, dnsmasq: {} },
      cluster: { nodes: { talos: {} } },
      server: {}
    };

    const lines = yamlText.split('\n');
    let currentSection = null;
    let currentSubsection = null;

    for (const line of lines) {
      const trimmed = line.trim();
      if (!trimmed || trimmed.startsWith('#')) continue;

      if (trimmed.startsWith('cloud:')) currentSection = 'cloud';
      else if (trimmed.startsWith('cluster:')) currentSection = 'cluster';
      else if (trimmed.startsWith('server:')) currentSection = 'server';
      else if (trimmed.startsWith('dns:')) currentSubsection = 'dns';
      else if (trimmed.startsWith('router:')) currentSubsection = 'router';
      else if (trimmed.startsWith('dnsmasq:')) currentSubsection = 'dnsmasq';
      else if (trimmed.startsWith('nodes:')) currentSubsection = 'nodes';
      else if (trimmed.startsWith('talos:')) currentSubsection = 'talos';
      else if (trimmed.includes(':')) {
        const [key, value] = trimmed.split(':').map(s => s.trim());
        const cleanValue = value.replace(/"/g, '');

        if (currentSection === 'cloud') {
          if (currentSubsection === 'dns') config.cloud.dns[key] = cleanValue;
          else if (currentSubsection === 'router') config.cloud.router[key] = cleanValue;
          else if (currentSubsection === 'dnsmasq') config.cloud.dnsmasq[key] = cleanValue;
          else config.cloud[key] = cleanValue;
        } else if (currentSection === 'cluster') {
          if (currentSubsection === 'nodes') {
            // Skip nodes level
          } else if (currentSubsection === 'talos') {
            config.cluster.nodes.talos[key] = cleanValue;
          } else {
            config.cluster[key] = cleanValue;
          }
        } else if (currentSection === 'server') {
          config.server[key] = key === 'port' ? parseInt(cleanValue) : cleanValue;
        }
      }
    }

    return config;
  };

  const formatTimestamp = (timestamp) => {
    return new Date(timestamp).toLocaleString();
  };

  const Message = ({ message, type }) => {
    if (!message) return null;
    return (
      <div className={`message ${type}`}>
        {message}
      </div>
    );
  };

  return (
    <div className="App">
      <div className="container">
        <h1>🌩️ Wild Cloud Central Management</h1>
        
        {/* System Status Section */}
        <div className="section">
          <h2>System Status</h2>
          <div className="section-content">
            <div className="button-group">
              <button onClick={fetchStatus} disabled={loading.status}>
                {loading.status ? '⏳ Checking...' : '🔄 Refresh Status'}
              </button>
              <button onClick={fetchHealth} disabled={loading.health}>
                {loading.health ? '⏳ Checking...' : '🏥 Check Health'}
              </button>
            </div>
            
            <Message message={messages.health?.message} type={messages.health?.type} />
            
            {status && (
              <div className="status-display">
                <div className="status-grid">
                  <div className="status-item">
                    <span className="status-label">Status:</span>
                    <span className={`status-value ${status.status === 'running' ? 'running' : 'stopped'}`}>
                      {status.status}
                    </span>
                  </div>
                  <div className="status-item">
                    <span className="status-label">Version:</span>
                    <span className="status-value">{status.version}</span>
                  </div>
                  <div className="status-item">
                    <span className="status-label">Uptime:</span>
                    <span className="status-value">{status.uptime}</span>
                  </div>
                  <div className="status-item">
                    <span className="status-label">Last Updated:</span>
                    <span className="status-value">{formatTimestamp(status.timestamp)}</span>
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Configuration Section */}
        <div className="section">
          <h2>Configuration</h2>
          <div className="section-content">
            <button onClick={fetchConfig} disabled={loading.config}>
              {loading.config ? '⏳ Loading...' : '📄 Load Current Config'}
            </button>
            
            <Message message={messages.config?.message} type={messages.config?.type} />
            
            {showConfigSetup && (
              <div className="config-setup">
                <h3>Initial Configuration Setup</h3>
                <p>No configuration found. Please provide the initial configuration:</p>
                <textarea
                  value={configText}
                  onChange={(e) => setConfigText(e.target.value)}
                  rows="20"
                  cols="80"
                  placeholder="Enter YAML configuration..."
                />
                <div className="button-group">
                  <button onClick={createConfig} disabled={loading.createConfig}>
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

        {/* DNS/DHCP Management Section */}
        <div className="section">
          <h2>DNS/DHCP Management</h2>
          <div className="section-content">
            <div className="button-group">
              <button onClick={generateDnsmasqConfig} disabled={loading.dnsmasq}>
                {loading.dnsmasq ? '⏳ Generating...' : '⚙️ Generate Dnsmasq Config'}
              </button>
              <button onClick={restartDnsmasq} disabled={loading.restart}>
                {loading.restart ? '⏳ Restarting...' : '🔄 Restart Dnsmasq'}
              </button>
            </div>
            
            <Message message={messages.dnsmasq?.message} type={messages.dnsmasq?.type} />
            
            {dnsmasqConfig && (
              <pre className="config-display">
                {dnsmasqConfig}
              </pre>
            )}
          </div>
        </div>

        {/* PXE Boot Assets Section */}
        <div className="section">
          <h2>PXE Boot Assets</h2>
          <div className="section-content">
            <button onClick={downloadAssets} disabled={loading.assets}>
              {loading.assets ? '⏳ Downloading...' : '⬇️ Download/Update PXE Assets'}
            </button>
            
            <Message message={messages.assets?.message} type={messages.assets?.type} />
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
