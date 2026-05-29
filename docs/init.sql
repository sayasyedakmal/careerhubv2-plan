-- CareerHubV2 Database Initialization Script
-- Target: Microsoft SQL Server

-- Disable constraints temporarily to make dropping tables cleaner
-- Drop tables if they exist in reverse dependency order
IF OBJECT_ID('dbo.AuditLogs', 'U') IS NOT NULL DROP TABLE dbo.AuditLogs;
IF OBJECT_ID('dbo.PasswordResets', 'U') IS NOT NULL DROP TABLE dbo.PasswordResets;
IF OBJECT_ID('dbo.UserRoles', 'U') IS NOT NULL DROP TABLE dbo.UserRoles;
IF OBJECT_ID('dbo.RolePermissions', 'U') IS NOT NULL DROP TABLE dbo.RolePermissions;
IF OBJECT_ID('dbo.Permissions', 'U') IS NOT NULL DROP TABLE dbo.Permissions;
IF OBJECT_ID('dbo.Roles', 'U') IS NOT NULL DROP TABLE dbo.Roles;
IF OBJECT_ID('dbo.Jobs', 'U') IS NOT NULL DROP TABLE dbo.Jobs;
IF OBJECT_ID('dbo.Users', 'U') IS NOT NULL DROP TABLE dbo.Users;
IF OBJECT_ID('dbo.Categories', 'U') IS NOT NULL DROP TABLE dbo.Categories;
IF OBJECT_ID('dbo.JobTypes', 'U') IS NOT NULL DROP TABLE dbo.JobTypes;

-- =========================================================================
-- 1. Create Core Tables
-- =========================================================================

-- Create Categories Table
CREATE TABLE Categories (
    CategoryID INT IDENTITY(1,1) PRIMARY KEY,
    CategoryName NVARCHAR(100) NOT NULL,
    IconName NVARCHAR(100) NULL
);

-- Create JobTypes Table
CREATE TABLE JobTypes (
    JobTypeID INT IDENTITY(1,1) PRIMARY KEY,
    TypeName NVARCHAR(50) NOT NULL
);

-- Create Users Table (Supports both local & Entra ID auth)
CREATE TABLE Users (
    UserID INT IDENTITY(1,1) PRIMARY KEY,
    MicrosoftObjectID NVARCHAR(100) UNIQUE NULL, -- Null for Alumni
    Email NVARCHAR(255) UNIQUE NOT NULL,
    PasswordHash NVARCHAR(255) NULL,            -- Null for Entra ID SSO users
    DisplayName NVARCHAR(255) NULL,
    Phone NVARCHAR(50) NULL,
    UserType NVARCHAR(50) NOT NULL,              -- 'ActiveStudent', 'Alumni', 'Staff', 'SystemAdmin'
    StudentID NVARCHAR(50) UNIQUE NULL,          -- Null for Staff/Admin
    RegistrationStatus NVARCHAR(50) DEFAULT 'N/A', -- 'Pending', 'Approved', 'Denied', 'N/A'
    CreatedAt DATETIME DEFAULT GETDATE(),
    LastLoginAt DATETIME NULL
);

-- Create PasswordResets Table
CREATE TABLE PasswordResets (
    ResetID INT IDENTITY(1,1) PRIMARY KEY,
    UserID INT NOT NULL,
    TokenHash NVARCHAR(255) NOT NULL,
    ExpiresAt DATETIME NOT NULL,
    CONSTRAINT FK_PasswordResets_Users FOREIGN KEY (UserID) REFERENCES Users(UserID) ON DELETE CASCADE
);

-- Create Jobs Table
CREATE TABLE Jobs (
    JobID INT IDENTITY(1,1) PRIMARY KEY,
    JobTitle NVARCHAR(255) NOT NULL,
    CompanyName NVARCHAR(255) NOT NULL,
    CategoryID INT NOT NULL,
    JobTypeID INT NOT NULL,
    IsActive BIT DEFAULT 1,
    HasSalary BIT DEFAULT 0,
    SalaryMin DECIMAL(18, 2) NULL,
    SalaryMax DECIMAL(18, 2) NULL,
    Deadline DATETIME NULL,
    Location NVARCHAR(255) NULL,
    PositionCount INT DEFAULT 1,
    JobDescription NVARCHAR(MAX) NULL,     -- Matches Overview Tab detail
    Responsibilities NVARCHAR(MAX) NULL,   -- Matches Details Tab narrative
    Requirements NVARCHAR(MAX) NULL,       -- Matches Requirements Tab checklist
    AdditionalInformation NVARCHAR(MAX) NULL,
    HowToApply NVARCHAR(MAX) NULL,          -- Matches Apply Tab instructions
    CreatedAt DATETIME DEFAULT GETDATE(),
    UpdatedAt DATETIME DEFAULT GETDATE(),
    
    CONSTRAINT FK_Jobs_Categories FOREIGN KEY (CategoryID) REFERENCES Categories(CategoryID),
    CONSTRAINT FK_Jobs_JobTypes FOREIGN KEY (JobTypeID) REFERENCES JobTypes(JobTypeID)
);

-- Create AuditLogs Table
CREATE TABLE AuditLogs (
    LogID INT IDENTITY(1,1) PRIMARY KEY,
    UserID INT NULL,
    Action NVARCHAR(255) NOT NULL,
    Resource NVARCHAR(255) NOT NULL,
    Details NVARCHAR(MAX) NULL,
    Timestamp DATETIME DEFAULT GETDATE(),
    CONSTRAINT FK_AuditLogs_Users FOREIGN KEY (UserID) REFERENCES Users(UserID) ON DELETE SET NULL
);

-- =========================================================================
-- 2. Create RBAC Tables
-- =========================================================================

-- Create Roles Table
CREATE TABLE Roles (
    RoleID INT IDENTITY(1,1) PRIMARY KEY,
    RoleName NVARCHAR(100) UNIQUE NOT NULL
);

-- Create Permissions Table
CREATE TABLE Permissions (
    PermissionID INT IDENTITY(1,1) PRIMARY KEY,
    PermissionName NVARCHAR(100) UNIQUE NOT NULL
);

-- Create RolePermissions Junction Table
CREATE TABLE RolePermissions (
    RoleID INT NOT NULL,
    PermissionID INT NOT NULL,
    CONSTRAINT PK_RolePermissions PRIMARY KEY (RoleID, PermissionID),
    CONSTRAINT FK_RolePermissions_Roles FOREIGN KEY (RoleID) REFERENCES Roles(RoleID) ON DELETE CASCADE,
    CONSTRAINT FK_RolePermissions_Permissions FOREIGN KEY (PermissionID) REFERENCES Permissions(PermissionID) ON DELETE CASCADE
);

-- Create UserRoles Junction Table
CREATE TABLE UserRoles (
    UserID INT NOT NULL,
    RoleID INT NOT NULL,
    CONSTRAINT PK_UserRoles PRIMARY KEY (UserID, RoleID),
    CONSTRAINT FK_UserRoles_Users FOREIGN KEY (UserID) REFERENCES Users(UserID) ON DELETE CASCADE,
    CONSTRAINT FK_UserRoles_Roles FOREIGN KEY (RoleID) REFERENCES Roles(RoleID) ON DELETE CASCADE
);

-- =========================================================================
-- 3. Seed Initial Lookup Data
-- =========================================================================

-- Seed Categories with Icon Mappings
SET IDENTITY_INSERT Categories ON;
INSERT INTO Categories (CategoryID, CategoryName, IconName) VALUES 
(1, 'Computing & IT', 'CommandLineIcon'),
(2, 'Engineering', 'WrenchIcon'),
(3, 'Business & Finance', 'BriefcaseIcon'),
(4, 'Hospitality & Tourism', 'AcademicCapIcon'),
(5, 'Communication & Creative', 'MegaphoneIcon'),
(6, 'Art & Design', 'PaintBrushIcon');
SET IDENTITY_INSERT Categories OFF;

-- Seed Job Types
SET IDENTITY_INSERT JobTypes ON;
INSERT INTO JobTypes (JobTypeID, TypeName) VALUES 
(1, 'Full-Time'),
(2, 'Part-Time'),
(3, 'Internship'),
(4, 'Contract');
SET IDENTITY_INSERT JobTypes OFF;

-- Seed Roles
SET IDENTITY_INSERT Roles ON;
INSERT INTO Roles (RoleID, RoleName) VALUES 
(1, 'System Admin'),
(2, 'SAC Department'),
(3, 'Active Student'),
(4, 'Alumni');
SET IDENTITY_INSERT Roles OFF;

-- Seed Permissions
SET IDENTITY_INSERT Permissions ON;
INSERT INTO Permissions (PermissionID, PermissionName) VALUES 
(1, 'manage_user_group'),
(2, 'manage_users'),
(3, 'approve_alumni'),
(4, 'create_job'),
(5, 'edit_job'),
(6, 'delete_job'),
(7, 'toggle_job'),
(8, 'view_jobs'),
(9, 'apply_job');
SET IDENTITY_INSERT Permissions OFF;

-- =========================================================================
-- 4. Seed Role-Permission Mappings
-- =========================================================================

-- System Admin permissions: All permissions (1-9)
INSERT INTO RolePermissions (RoleID, PermissionID) VALUES 
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6), (1, 7), (1, 8), (1, 9);

-- SAC Department permissions: Job management + alumni approval + general views
INSERT INTO RolePermissions (RoleID, PermissionID) VALUES 
(2, 3), (2, 4), (2, 5), (2, 6), (2, 7), (2, 8), (2, 9);

-- Active Student permissions: view & apply
INSERT INTO RolePermissions (RoleID, PermissionID) VALUES 
(3, 8), (3, 9);

-- Alumni permissions: view & apply
INSERT INTO RolePermissions (RoleID, PermissionID) VALUES 
(4, 8), (4, 9);
