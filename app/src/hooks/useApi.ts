import { useState, useCallback } from 'react';
import { apiService } from '../services/api';
import { Status, Config, LoadingState, Messages } from '../types';

export const useApi = () => {
  const [status, setStatus] = useState<Status | null>(null);
  const [config, setConfig] = useState<Config | null>(null);
  const [dnsmasqConfig, setDnsmasqConfig] = useState<string>('');
  const [messages, setMessages] = useState<Messages>({});
  const [loading, setLoading] = useState<LoadingState>({ status: true });
  const [showConfigSetup, setShowConfigSetup] = useState<boolean>(false);

  const setLoadingState = (key: string, value: boolean) => {
    setLoading(prev => ({ ...prev, [key]: value }));
  };

  const setMessage = (key: string, message: string | null, type: 'info' | 'success' | 'error' = 'info') => {
    if (message === null) {
      setMessages(prev => {
        const newMessages = { ...prev };
        delete newMessages[key];
        return newMessages;
      });
    } else {
      setMessages(prev => ({ ...prev, [key]: { message, type } }));
    }
  };

  const fetchStatus = useCallback(async () => {
    try {
      setLoadingState('status', true);
      const data = await apiService.getStatus();
      setStatus(data);
      setMessage('status', null);
    } catch (err) {
      const error = err as Error;
      setMessage('status', `Failed to fetch status: ${error.message}`, 'error');
    } finally {
      setLoadingState('status', false);
    }
  }, []);

  const fetchHealth = useCallback(async () => {
    try {
      setLoadingState('health', true);
      const data = await apiService.getHealth();
      setMessage('health', `Service: ${data.service} - Status: ${data.status}`, 'success');
    } catch (err) {
      const error = err as Error;
      setMessage('health', `Health check failed: ${error.message}`, 'error');
    } finally {
      setLoadingState('health', false);
    }
  }, []);

  const fetchConfig = useCallback(async () => {
    try {
      setLoadingState('config', true);
      const data = await apiService.getConfig();
      
      console.log('fetchConfig response:', data);
      console.log('data.configured:', data.configured);
      console.log('data.config:', data.config);
      
      if (data.configured === false) {
        console.log('Config not configured, showing setup');
        setShowConfigSetup(true);
        setMessage('config', data.message || 'No configuration found', 'error');
        setConfig(null);
      } else {
        console.log('Config is configured, setting config data');
        setShowConfigSetup(false);
        setConfig(data.config || null);
        setMessage('config', 'Configuration loaded successfully', 'success');
      }
    } catch (err) {
      const error = err as Error;
      console.error('fetchConfig error:', error);
      setMessage('config', `Failed to load config: ${error.message}`, 'error');
    } finally {
      setLoadingState('config', false);
    }
  }, []);

  const createConfig = async (configObj: Config) => {
    try {
      setLoadingState('createConfig', true);
      const data = await apiService.createConfig(configObj);
      setMessage('config', `Configuration created successfully! Status: ${data.status}`, 'success');
      setShowConfigSetup(false);
      await fetchConfig();
    } catch (err) {
      const error = err as Error;
      setMessage('config', `Error creating config: ${error.message}`, 'error');
    } finally {
      setLoadingState('createConfig', false);
    }
  };

  const generateDnsmasqConfig = async () => {
    try {
      setLoadingState('dnsmasq', true);
      const data = await apiService.getDnsmasqConfig();
      setDnsmasqConfig(data);
      setMessage('dnsmasq', 'Dnsmasq config generated successfully', 'success');
    } catch (err) {
      const error = err as Error;
      setMessage('dnsmasq', `Failed to generate dnsmasq config: ${error.message}`, 'error');
    } finally {
      setLoadingState('dnsmasq', false);
    }
  };

  const restartDnsmasq = async () => {
    try {
      setLoadingState('restart', true);
      const data = await apiService.restartDnsmasq();
      setMessage('dnsmasq', `Dnsmasq restart: ${data.status}`, 'success');
    } catch (err) {
      const error = err as Error;
      setMessage('dnsmasq', `Failed to restart dnsmasq: ${error.message}`, 'error');
    } finally {
      setLoadingState('restart', false);
    }
  };

  const downloadAssets = async () => {
    try {
      setLoadingState('assets', true);
      const data = await apiService.downloadPXEAssets();
      setMessage('assets', `PXE Assets: ${data.status}`, 'success');
    } catch (err) {
      const error = err as Error;
      setMessage('assets', `Failed to download assets: ${error.message}`, 'error');
    } finally {
      setLoadingState('assets', false);
    }
  };

  return {
    // State
    status,
    config,
    dnsmasqConfig,
    messages,
    loading,
    showConfigSetup,
    
    // Actions
    fetchStatus,
    fetchHealth,
    fetchConfig,
    createConfig,
    generateDnsmasqConfig,
    restartDnsmasq,
    downloadAssets,
    setMessage,
    setShowConfigSetup
  };
};