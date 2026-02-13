// User info from Keycloak token
export interface User {
  andrewId: string;
  name: string;
  email: string;
  isAdmin: boolean;
}

// Auth state
export interface AuthState {
  isAuthenticated: boolean;
  isLoading: boolean;
  user: User | null;
  error: string | null;
}

// Post submission
export interface PostSubmission {
  andrewId: string;
  name: string;
  caption: string;
  files: File[];
}

// API response types
export interface AuthCheckResponse {
  isAdmin: boolean;
  andrewId: string;
  name: string;
}

export interface PostSubmitResponse {
  success: boolean;
  message: string;
  postId?: number;
}

// File with preview and order
export interface UploadFile {
  id: string;
  file: File;
  preview: string;
  order: number;
}

// Keycloak token payload
export interface KeycloakTokenPayload {
  sub: string;
  preferred_username: string;
  given_name?: string;
  family_name?: string;
  name?: string;
  email?: string;
  exp: number;
  iat: number;
}
