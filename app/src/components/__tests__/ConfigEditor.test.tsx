import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '../../test/test-utils';
import { ConfigEditor } from '../ConfigEditor';
import { apiService } from '../../services/api';

// Mock the API service
vi.mock('../../services/api', () => ({
  apiService: {
    getConfigYaml: vi.fn(),
    updateConfigYaml: vi.fn(),
  },
}));

const mockYamlContent = `server:
  host: "0.0.0.0"
  port: 5055
cloud:
  domain: "wildcloud.local"
  internalDomain: "cluster.local"
  dhcpRange: "192.168.8.100,192.168.8.200"
  dns:
    ip: "192.168.8.50"
  router:
    ip: "192.168.8.1"
  dnsmasq:
    interface: "eth0"
cluster:
  endpointIp: "192.168.8.60"
  nodes:
    talos:
      version: "v1.8.0"`;

describe('ConfigEditor', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should render config button', () => {
    vi.mocked(apiService.getConfigYaml).mockResolvedValue(mockYamlContent);
    
    render(<ConfigEditor />);
    
    expect(screen.getByRole('button', { name: /config/i })).toBeInTheDocument();
  });

  it('should open dialog when config button is clicked', async () => {
    vi.mocked(apiService.getConfigYaml).mockResolvedValue(mockYamlContent);
    
    render(<ConfigEditor />);
    
    const configButton = screen.getByRole('button', { name: /config/i });
    fireEvent.click(configButton);
    
    await waitFor(() => {
      expect(screen.getByText('Configuration Editor')).toBeInTheDocument();
    });
  });

  it('should load and display YAML content', async () => {
    vi.mocked(apiService.getConfigYaml).mockResolvedValue(mockYamlContent);
    
    render(<ConfigEditor />);
    
    const configButton = screen.getByRole('button', { name: /config/i });
    fireEvent.click(configButton);
    
    await waitFor(() => {
      expect(screen.getByDisplayValue(/server:/)).toBeInTheDocument();
    });
    
    expect(apiService.getConfigYaml).toHaveBeenCalled();
  });

  it('should handle API errors gracefully', async () => {
    const errorMessage = 'Failed to fetch config';
    vi.mocked(apiService.getConfigYaml).mockRejectedValue(new Error(errorMessage));
    
    render(<ConfigEditor />);
    
    const configButton = screen.getByRole('button', { name: /config/i });
    fireEvent.click(configButton);
    
    await waitFor(() => {
      expect(screen.getByText(errorMessage)).toBeInTheDocument();
    });
  });

  it('should show missing endpoints message for 404 errors', async () => {
    vi.mocked(apiService.getConfigYaml).mockRejectedValue(new Error('HTTP error! status: 404'));
    
    render(<ConfigEditor />);
    
    const configButton = screen.getByRole('button', { name: /config/i });
    fireEvent.click(configButton);
    
    await waitFor(() => {
      expect(screen.getByText('Backend Endpoints Required')).toBeInTheDocument();
      expect(screen.getByText('GET /api/v1/config/yaml - Read config.yaml file contents')).toBeInTheDocument();
      expect(screen.getByText('PUT /api/v1/config/yaml - Write raw YAML to config.yaml file')).toBeInTheDocument();
    });
  });

  it('should disable editor when endpoints are missing', async () => {
    vi.mocked(apiService.getConfigYaml).mockRejectedValue(new Error('HTTP error! status: 404'));
    
    render(<ConfigEditor />);
    
    const configButton = screen.getByRole('button', { name: /config/i });
    fireEvent.click(configButton);
    
    await waitFor(() => {
      const textarea = screen.getByRole('textbox');
      expect(textarea).toBeDisabled();
      expect(screen.getByRole('button', { name: /update config/i })).toBeDisabled();
    });
  });

  it('should detect changes in textarea', async () => {
    vi.mocked(apiService.getConfigYaml).mockResolvedValue(mockYamlContent);
    
    render(<ConfigEditor />);
    
    const configButton = screen.getByRole('button', { name: /config/i });
    fireEvent.click(configButton);
    
    await waitFor(() => {
      expect(screen.getByDisplayValue(/server:/)).toBeInTheDocument();
    });
    
    const textarea = screen.getByRole('textbox');
    fireEvent.change(textarea, { target: { value: mockYamlContent + '\n# modified' } });
    
    await waitFor(() => {
      expect(screen.getByText('Unsaved changes')).toBeInTheDocument();
    });
  });

  it('should enable update button when changes are made', async () => {
    vi.mocked(apiService.getConfigYaml).mockResolvedValue(mockYamlContent);
    vi.mocked(apiService.updateConfigYaml).mockResolvedValue({ status: 'success' });
    
    render(<ConfigEditor />);
    
    const configButton = screen.getByRole('button', { name: /config/i });
    fireEvent.click(configButton);
    
    await waitFor(() => {
      expect(screen.getByDisplayValue(/server:/)).toBeInTheDocument();
    });
    
    const textarea = screen.getByRole('textbox');
    const updateButton = screen.getByRole('button', { name: /update config/i });
    
    // Initially disabled
    expect(updateButton).toBeDisabled();
    
    // Make changes
    fireEvent.change(textarea, { target: { value: mockYamlContent + '\n# modified' } });
    
    await waitFor(() => {
      expect(updateButton).not.toBeDisabled();
    });
  });

  it('should call updateConfigYaml when update button is clicked', async () => {
    const modifiedContent = mockYamlContent + '\n# modified';
    vi.mocked(apiService.getConfigYaml).mockResolvedValue(mockYamlContent);
    vi.mocked(apiService.updateConfigYaml).mockResolvedValue({ status: 'success' });
    
    render(<ConfigEditor />);
    
    const configButton = screen.getByRole('button', { name: /config/i });
    fireEvent.click(configButton);
    
    await waitFor(() => {
      expect(screen.getByDisplayValue(/server:/)).toBeInTheDocument();
    });
    
    const textarea = screen.getByRole('textbox');
    fireEvent.change(textarea, { target: { value: modifiedContent } });
    
    const updateButton = screen.getByRole('button', { name: /update config/i });
    fireEvent.click(updateButton);
    
    await waitFor(() => {
      expect(apiService.updateConfigYaml).toHaveBeenCalledWith(modifiedContent);
    });
  });
});