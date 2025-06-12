import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '../../test/test-utils';
import { StatusSection } from '../StatusSection';
import { apiService } from '../../services/api';

// Mock the API service
vi.mock('../../services/api', () => ({
  apiService: {
    getStatus: vi.fn(),
  },
}));

describe('StatusSection', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should render loading state initially', () => {
    vi.mocked(apiService.getStatus).mockImplementation(
      () => new Promise(() => {}) // Never resolves
    );

    render(<StatusSection />);

    expect(screen.getByText('Server Status')).toBeInTheDocument();
    expect(screen.getByText('Refreshing...')).toBeInTheDocument();
  });

  it('should render status data when loaded successfully', async () => {
    const mockStatus = {
      status: 'running',
      version: '1.0.0',
      uptime: '2 hours',
      timestamp: '2024-01-01T00:00:00Z',
    };

    vi.mocked(apiService.getStatus).mockResolvedValue(mockStatus);

    render(<StatusSection />);

    await waitFor(() => {
      expect(screen.getByText('Refresh')).toBeInTheDocument();
    });

    expect(screen.getByText('running')).toBeInTheDocument();
    expect(screen.getByText('1.0.0')).toBeInTheDocument();
    expect(screen.getByText('2 hours')).toBeInTheDocument();
  });

  it('should render error state when loading fails', async () => {
    const mockError = new Error('Network error');
    vi.mocked(apiService.getStatus).mockRejectedValue(mockError);

    render(<StatusSection />);

    await waitFor(() => {
      expect(screen.getByText('Failed to fetch status: Network error')).toBeInTheDocument();
    });
  });

  it('should show loading spinner when refreshing', async () => {
    vi.mocked(apiService.getStatus).mockImplementation(
      () => new Promise(() => {}) // Never resolves
    );

    render(<StatusSection />);

    const refreshButton = screen.getByRole('button', { name: /refresh/i });
    expect(refreshButton).toBeDisabled();
    expect(screen.getByText('Refreshing...')).toBeInTheDocument();
  });

  it('should display status in JSON format', async () => {
    const mockStatus = {
      status: 'running',
      version: '1.0.0',
      uptime: '2 hours',
      timestamp: '2024-01-01T00:00:00Z',
    };

    vi.mocked(apiService.getStatus).mockResolvedValue(mockStatus);

    render(<StatusSection />);

    await waitFor(() => {
      expect(screen.getByText('Refresh')).toBeInTheDocument();
    });

    // Check that JSON is displayed in a pre element
    const jsonDisplay = screen.getByText(/"status": "running"/);
    expect(jsonDisplay).toBeInTheDocument();
    expect(jsonDisplay.closest('pre')).toBeInTheDocument();
  });

  it('should have accessible refresh button', async () => {
    const mockStatus = {
      status: 'running',
      version: '1.0.0',
      uptime: '2 hours',
      timestamp: '2024-01-01T00:00:00Z',
    };

    vi.mocked(apiService.getStatus).mockResolvedValue(mockStatus);

    render(<StatusSection />);

    await waitFor(() => {
      expect(screen.getByRole('button', { name: /refresh/i })).toBeInTheDocument();
    });

    const refreshButton = screen.getByRole('button', { name: /refresh/i });
    expect(refreshButton).not.toBeDisabled();
    expect(refreshButton).toHaveTextContent('Refresh');
  });
});