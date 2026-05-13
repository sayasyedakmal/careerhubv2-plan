-- CareerHubV2 Database Initialization Script
-- Target: Microsoft SQL Server

-- Create Categories Table
CREATE TABLE Categories (
    CategoryID INT IDENTITY(1,1) PRIMARY KEY,
    CategoryName NVARCHAR(100) NOT NULL
);

-- Create JobTypes Table
CREATE TABLE JobTypes (
    JobTypeID INT IDENTITY(1,1) PRIMARY KEY,
    TypeName NVARCHAR(50) NOT NULL
);

-- Create Users Table
CREATE TABLE Users (
    UserID INT IDENTITY(1,1) PRIMARY KEY,
    MicrosoftObjectID NVARCHAR(100) UNIQUE NOT NULL,
    Email NVARCHAR(255) NOT NULL,
    DisplayName NVARCHAR(255),
    CreatedAt DATETIME DEFAULT GETDATE(),
    LastLoginAt DATETIME
);

-- Create Jobs Table
CREATE TABLE Jobs (
    JobID INT IDENTITY(1,1) PRIMARY KEY,
    JobTitle NVARCHAR(255) NOT NULL,
    CompanyName NVARCHAR(255) NOT NULL,
    CategoryID INT,
    JobTypeID INT,
    HasSalary BIT DEFAULT 0,
    SalaryMin DECIMAL(18, 2),
    SalaryMax DECIMAL(18, 2),
    Deadline DATETIME,
    Location NVARCHAR(255),
    PositionCount INT DEFAULT 1,
    JobDescription NVARCHAR(MAX),
    Responsibilities NVARCHAR(MAX),
    Requirements NVARCHAR(MAX),
    AdditionalInformation NVARCHAR(MAX),
    HowToApply NVARCHAR(MAX),
    CreatedAt DATETIME DEFAULT GETDATE(),
    UpdatedAt DATETIME DEFAULT GETDATE(),
    
    CONSTRAINT FK_Jobs_Categories FOREIGN KEY (CategoryID) REFERENCES Categories(CategoryID),
    CONSTRAINT FK_Jobs_JobTypes FOREIGN KEY (JobTypeID) REFERENCES JobTypes(JobTypeID)
);

-- Seed Initial Categories
INSERT INTO Categories (CategoryName) VALUES 
('Marketing'), ('Engineering'), ('Finance'), ('Human Resources'), ('Design'), ('Sales');

-- Seed Initial Job Types
INSERT INTO JobTypes (TypeName) VALUES 
('Full-Time'), ('Part-Time'), ('Internship'), ('Contract');
