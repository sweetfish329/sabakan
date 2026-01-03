/**
 * @file Public API for shared components.
 * @module shared
 */

// Status Chip
export { StatusChipComponent, type StatusType } from "./components/status-chip/status-chip";

// Loading Overlay
export { LoadingOverlayComponent } from "./components/loading-overlay/loading-overlay";

// Empty State
export { EmptyStateComponent } from "./components/empty-state/empty-state";

// Error Banner
export { ErrorBannerComponent, type ErrorSeverity } from "./components/error-banner/error-banner";

// Page Header
export { PageHeaderComponent } from "./components/page-header/page-header";

// Confirm Dialog
export {
  ConfirmDialogComponent,
  ConfirmDialogPreviewComponent,
  type ConfirmDialogData,
} from "./components/confirm-dialog/confirm-dialog";

// Info Card
export { InfoCardComponent, type InfoItem } from "./components/info-card/info-card";

// Action Button
export {
  ActionButtonComponent,
  type ActionButtonVariant,
  type ActionButtonSize,
} from "./components/action-button/action-button";

// Social Login Button
export {
  SocialLoginButtonComponent,
  type OAuthProvider,
} from "./components/social-login-button/social-login-button";

// Login Form
export { LoginFormComponent, type LoginFormData } from "./components/login-form/login-form";

// Register Form
export {
  RegisterFormComponent,
  type RegisterFormData,
} from "./components/register-form/register-form";
