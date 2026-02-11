# Audit Logging and Notifications Implementation

## Overview
This implementation adds comprehensive audit logging and notification features to the SZI Registry application.

## New Models

### AuditLog
Records all changes made to SZI records:
- User who made the change
- Action type (CREATE, UPDATE, DELETE)
- Old and new values (JSON format)
- Timestamp and description

### Notification
Manages user notifications:
- Certificate expiration alerts (3 levels: EXPIRED, CRITICAL, EXPIRING_SOON)
- Priority-based notifications (1-4)
- Read/unread status tracking
- Linked to specific SZI records

## Key Features

### 1. Automatic Audit Logging
All SZI record operations (create, update, delete) are automatically logged with:
- Full record snapshots (before/after)
- User identification
- Timestamps

### 2. Smart Notifications
- **EXPIRED**: Certificate has already expired (Priority 4, critical)
- **CRITICAL**: Less than 7 days until expiration (Priority 3)
- **EXPIRING_SOON**: 7-30 days until expiration (Priority 2)

### 3. Notification System
- Automatic generation at startup
- Periodic checks every 24 hours
- Avoids duplicate notifications
- Badge counter on dashboard button
- Full notification management UI

### 4. Database Schema
Foreign key relationships established:
- SZIRecord -> User (CASCADE on delete)
- Notification -> User (CASCADE on delete)
- Notification -> SZIRecord (CASCADE on delete)
- AuditLog -> User (SET NULL on delete)

## Usage

### Accessing Notifications
1. Login to the application
2. Click the "🔔 Уведомления" button on the dashboard
3. Unread count is displayed in the button text

### Notification Management
- Mark individual notifications as read
- Mark all notifications as read at once
- Delete notifications
- Notifications are sorted by priority and date

### Audit Trail
The audit trail is automatically maintained and can be queried through:
- `GetUserAuditLogs()` - User-specific actions
- `GetRecordAuditLogs()` - Record-specific history
- `GetAllAuditLogs()` - Complete audit trail (admin)

## Technical Details

### New Files
- `models/gorm_models.go` - Updated with new models
- `repositories/audit_log_repository.go` - Audit log data access
- `repositories/notification_repository.go` - Notification data access
- `services/audit_service.go` - Audit logging business logic
- `services/notification_service.go` - Notification management
- `ui/notifications_window.go` - Notification UI

### Modified Files
- `services/szi_service.go` - Added audit logging to CRUD operations
- `ui/dashboard_window.go` - Added notification button
- `repositories/db_manager.go` - Added new models to AutoMigrate
- `main.go` - Added notification generation on startup and periodic checks

## Testing
The code has been verified with:
- `go vet` - No syntax errors
- `gofmt` - Proper code formatting
- Manual code review

Full testing requires:
1. Running PostgreSQL database
2. GUI environment (X11) for the Fyne UI
3. Test user accounts and SZI records

## Future Enhancements
Possible improvements:
- Email notifications for critical alerts
- Configurable notification thresholds
- Export audit logs to external systems
- Advanced filtering in notification UI
