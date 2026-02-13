import { writable, derived } from 'svelte/store';

export interface User {
  andrewId: string;
  name: string;
  email: string;
  isAdmin: boolean;
}

export interface AuthState {
  isAuthenticated: boolean;
  isLoading: boolean;
  user: User | null;
  accessToken: string | null;
  error: string | null;
}

const initialState: AuthState = {
  isAuthenticated: false,
  isLoading: true,
  user: null,
  accessToken: null,
  error: null,
};

function createAuthStore() {
  const { subscribe, set, update } = writable<AuthState>(initialState);

  return {
    subscribe,
    setUser: (user: User, accessToken: string) => {
      update(state => ({
        ...state,
        isAuthenticated: true,
        isLoading: false,
        user,
        accessToken,
        error: null,
      }));
    },
    setLoading: (isLoading: boolean) => {
      update(state => ({ ...state, isLoading }));
    },
    setError: (error: string) => {
      update(state => ({
        ...state,
        isAuthenticated: false,
        isLoading: false,
        user: null,
        accessToken: null,
        error,
      }));
    },
    logout: () => {
      set(initialState);
      update(state => ({ ...state, isLoading: false }));
    },
    reset: () => set(initialState),
  };
}

export const authStore = createAuthStore();

// Derived stores for convenience
export const isAuthenticated = derived(authStore, $auth => $auth.isAuthenticated);
export const isLoading = derived(authStore, $auth => $auth.isLoading);
export const currentUser = derived(authStore, $auth => $auth.user);
export const isAdmin = derived(authStore, $auth => $auth.user?.isAdmin ?? false);

// Route store for SPA navigation
export type Route = 'login' | 'callback' | 'home' | 'admin';
export const currentRoute = writable<Route>('login');
