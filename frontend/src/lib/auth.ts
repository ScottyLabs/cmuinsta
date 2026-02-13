// OIDC Authentication Service
// Uses backend for secure token exchange with client secret

export interface UserInfo {
  andrewId: string;
  email: string;
  name: string;
  givenName?: string;
}

export interface AuthState {
  isAuthenticated: boolean;
  isLoading: boolean;
  user: UserInfo | null;
  accessToken: string | null;
  isAdmin: boolean;
}

// Storage keys
const STORAGE_KEYS = {
  ACCESS_TOKEN: 'auth_access_token',
  USER_INFO: 'auth_user_info',
  IS_ADMIN: 'auth_is_admin',
  EXPIRES_AT: 'auth_expires_at',
};

// Login - get login URL from backend and redirect
export async function login(): Promise<void> {
  try {
    const response = await fetch('/api/auth/login-url');
    const data = await response.json();

    if (data.success && data.loginUrl) {
      window.location.href = data.loginUrl;
    } else {
      throw new Error('Failed to get login URL');
    }
  } catch (error) {
    console.error('Login error:', error);
    throw error;
  }
}

// Handle OAuth callback - send code to backend for secure token exchange
export async function handleCallback(): Promise<{ user: UserInfo; isAdmin: boolean } | null> {
  const urlParams = new URLSearchParams(window.location.search);
  const code = urlParams.get('code');
  const error = urlParams.get('error');

  if (error) {
    console.error('Auth error:', error, urlParams.get('error_description'));
    throw new Error(urlParams.get('error_description') || 'Authentication failed');
  }

  if (!code) {
    return null;
  }

  // Send code to backend for secure token exchange
  const response = await fetch('/api/auth/callback', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      code: code,
      redirectUri: `${window.location.origin}/oauth2/callback`,
    }),
  });

  if (!response.ok) {
    const errorData = await response.json();
    console.error('Token exchange failed:', errorData);
    throw new Error(errorData.error || 'Failed to exchange code for tokens');
  }

  const data = await response.json();

  if (!data.success) {
    throw new Error(data.error || 'Authentication failed');
  }

  // Store auth data
  const userInfo: UserInfo = {
    andrewId: data.user.andrewId,
    email: data.user.email,
    name: data.user.name,
    givenName: data.user.givenName,
  };

  localStorage.setItem(STORAGE_KEYS.ACCESS_TOKEN, data.accessToken);
  localStorage.setItem(STORAGE_KEYS.USER_INFO, JSON.stringify(userInfo));
  localStorage.setItem(STORAGE_KEYS.IS_ADMIN, JSON.stringify(data.isAdmin));

  // Calculate expiration time
  const expiresAt = Date.now() + data.expiresIn * 1000;
  localStorage.setItem(STORAGE_KEYS.EXPIRES_AT, expiresAt.toString());

  // Clear URL params
  window.history.replaceState({}, document.title, window.location.pathname);

  return { user: userInfo, isAdmin: data.isAdmin };
}

// Get current auth state
export function getAuthState(): AuthState {
  const accessToken = localStorage.getItem(STORAGE_KEYS.ACCESS_TOKEN);
  const userInfoStr = localStorage.getItem(STORAGE_KEYS.USER_INFO);
  const isAdminStr = localStorage.getItem(STORAGE_KEYS.IS_ADMIN);
  const expiresAtStr = localStorage.getItem(STORAGE_KEYS.EXPIRES_AT);

  // Check if token is expired
  if (expiresAtStr) {
    const expiresAt = parseInt(expiresAtStr, 10);
    if (Date.now() > expiresAt) {
      // Token expired, clear storage
      clearAuthStorage();
      return {
        isAuthenticated: false,
        isLoading: false,
        user: null,
        accessToken: null,
        isAdmin: false,
      };
    }
  }

  if (!accessToken || !userInfoStr) {
    return {
      isAuthenticated: false,
      isLoading: false,
      user: null,
      accessToken: null,
      isAdmin: false,
    };
  }

  try {
    const user = JSON.parse(userInfoStr) as UserInfo;
    const isAdmin = isAdminStr ? JSON.parse(isAdminStr) : false;

    return {
      isAuthenticated: true,
      isLoading: false,
      user,
      accessToken,
      isAdmin,
    };
  } catch {
    clearAuthStorage();
    return {
      isAuthenticated: false,
      isLoading: false,
      user: null,
      accessToken: null,
      isAdmin: false,
    };
  }
}

// Clear auth storage
function clearAuthStorage(): void {
  localStorage.removeItem(STORAGE_KEYS.ACCESS_TOKEN);
  localStorage.removeItem(STORAGE_KEYS.USER_INFO);
  localStorage.removeItem(STORAGE_KEYS.IS_ADMIN);
  localStorage.removeItem(STORAGE_KEYS.EXPIRES_AT);
}

// Logout - get logout URL from backend and redirect
export async function logout(): Promise<void> {
  try {
    const response = await fetch('/api/auth/logout-url');
    const data = await response.json();

    // Clear local storage first
    clearAuthStorage();

    if (data.success && data.logoutUrl) {
      window.location.href = data.logoutUrl;
    } else {
      // Fallback to just clearing storage and going home
      window.location.href = '/';
    }
  } catch (error) {
    console.error('Logout error:', error);
    // Clear storage and redirect home anyway
    clearAuthStorage();
    window.location.href = '/';
  }
}

// Get access token (for API calls)
export function getAccessToken(): string | null {
  // Check if expired first
  const expiresAtStr = localStorage.getItem(STORAGE_KEYS.EXPIRES_AT);
  if (expiresAtStr) {
    const expiresAt = parseInt(expiresAtStr, 10);
    if (Date.now() > expiresAt) {
      clearAuthStorage();
      return null;
    }
  }
  return localStorage.getItem(STORAGE_KEYS.ACCESS_TOKEN);
}

// Check if authenticated
export function isAuthenticated(): boolean {
  return getAccessToken() !== null;
}

// Validate token with backend (optional, for checking if session is still valid)
export async function validateToken(): Promise<boolean> {
  const accessToken = getAccessToken();
  if (!accessToken) {
    return false;
  }

  try {
    const response = await fetch('/api/auth/me', {
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    });

    if (!response.ok) {
      clearAuthStorage();
      return false;
    }

    const data = await response.json();
    return data.success === true;
  } catch {
    return false;
  }
}
