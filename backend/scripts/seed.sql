-- CareerHubV2 Sample Jobs Seed — Sprint 1
-- Run this AFTER docs/init.sql has been applied.
-- docs/init.sql already seeds Categories, JobTypes, Roles, Permissions, and RolePermissions.
-- This script only inserts sample job listings for local development and testing.

USE careerhubv2;
GO

INSERT INTO Jobs (
    JobTitle, CompanyName, CategoryID, JobTypeID,
    IsActive, HasSalary, SalaryMin, SalaryMax,
    DeadlineAt, State, City, PositionCount,
    JobDescription, Responsibilities, Requirements, AdditionalInformation, HowToApply
) VALUES
(
    'Junior Software Engineer',
    'TechNova Sdn Bhd',
    (SELECT CategoryID FROM Categories WHERE CategoryName = 'Computing & IT'),
    (SELECT JobTypeID FROM JobTypes WHERE TypeName = 'Full-Time'),
    1, 1, 3000.00, 4500.00,
    '2026-09-30', 'Selangor', 'Petaling Jaya', 2,
    'Join our product team to build scalable web applications used by thousands of Malaysian SMEs.',
    'Develop and maintain backend REST APIs;
Participate in code reviews and sprint planning;
Write unit and integration tests;
Collaborate with frontend and DevOps teams.',
    'Diploma or degree in Computer Science or related field;
Proficiency in at least one backend language (Go, Python, or Node.js);
Basic understanding of SQL databases;
Good communication skills.',
    'Hybrid work arrangement (3 days onsite);
Medical and dental benefits included.',
    'Send your CV and GitHub profile to careers@technova.com.my with subject line "Junior SE Application".'
),
(
    'Software Engineering Intern',
    'Axiata Digital Labs',
    (SELECT CategoryID FROM Categories WHERE CategoryName = 'Computing & IT'),
    (SELECT JobTypeID FROM JobTypes WHERE TypeName = 'Internship'),
    1, 1, 1000.00, 1500.00,
    '2026-08-15', 'Kuala Lumpur', 'Bangsar South', 3,
    'A 3-6 month internship supporting the platform engineering team on cloud infrastructure and API development.',
    'Assist in developing microservices using Go and Docker;
Help maintain CI/CD pipelines;
Write documentation for internal tools;
Participate in daily standups.',
    'Currently pursuing a degree in Software Engineering or Computer Science;
Familiar with Git version control;
Exposure to any cloud platform (AWS, Azure, or GCP) is a plus.',
    'Flexible working hours;
Internship allowance provided;
Free lunch on Fridays.',
    'Apply via our careers portal at axiatadigital.com/careers or email internship@axiatadigital.com.'
),
(
    'Civil & Structural Engineer',
    'Gamuda Engineering Berhad',
    (SELECT CategoryID FROM Categories WHERE CategoryName = 'Engineering'),
    (SELECT JobTypeID FROM JobTypes WHERE TypeName = 'Full-Time'),
    1, 1, 4500.00, 7000.00,
    '2026-10-01', 'Selangor', 'Shah Alam', 1,
    'Support the design and delivery of large-scale infrastructure projects including highways, LRT extensions, and drainage systems.',
    'Perform structural analysis and design calculations;
Prepare engineering drawings using AutoCAD;
Liaise with contractors and project managers on site;
Ensure designs comply with local building codes.',
    'Degree in Civil or Structural Engineering;
Graduate Engineer (GE) registration with BEM preferred;
Proficient in AutoCAD and STAAD.Pro;
Minimum 1 year experience (fresh graduates encouraged to apply).',
    'Project site allowance provided;
Career progression to Senior Engineer within 2 years.',
    'Email your resume to talent@gamuda.com.my with the subject "Civil Engineer Application".'
),
(
    'Event Executive',
    'Lavin Pharma (M) Sdn Bhd',
    (SELECT CategoryID FROM Categories WHERE CategoryName = 'Hospitality & Tourism'),
    (SELECT JobTypeID FROM JobTypes WHERE TypeName = 'Internship'),
    1, 1, 700.00, 1200.00,
    '2026-07-31', 'Selangor', 'Shah Alam', 2,
    'We are seeking energetic Event Executives to support beauty brand launch activations and roadshow campaigns across Klang Valley.',
    'Plan event booth layouts and logistics;
Manage social media posting campaigns during events;
Coordinate with caterers, vendors, and venue staff;
Prepare post-event reports with photos and attendance data.',
    'Pursuing or completed a diploma or degree in PR, Marketing, or Events Management;
Excellent communication and organizing skills;
Willing to work on weekends during event periods;
Own transport preferred.',
    'Flexible working hours;
Product allowance and skincare goodie bags provided.',
    'Send your resume and portfolio to careers@lavinpharma.com.my with subject "Event Executive Application".'
),
(
    'Business Analyst (Finance)',
    'CIMB Group',
    (SELECT CategoryID FROM Categories WHERE CategoryName = 'Business & Finance'),
    (SELECT JobTypeID FROM JobTypes WHERE TypeName = 'Full-Time'),
    1, 1, 4000.00, 6000.00,
    '2026-09-15', 'Kuala Lumpur', 'Bukit Bintang', 1,
    'Work closely with the digital banking product team to analyse business processes and translate requirements into technical specifications.',
    'Gather and document business requirements from stakeholders;
Map current and future-state process flows;
Produce detailed functional specification documents;
Support UAT planning and defect tracking.',
    'Degree in Business Administration, Finance, or IT;
Strong analytical and problem-solving skills;
Experience with process mapping tools (Visio, Lucidchart);
Knowledge of banking or fintech domain is an advantage.',
    'Hybrid work model (Tuesday-Thursday onsite);
Annual performance bonus.',
    'Apply through CIMB careers portal at cimb.com/careers or email talent@cimb.com.'
),
(
    'Part-Time Barista',
    'The Coffee Bean & Tea Leaf',
    (SELECT CategoryID FROM Categories WHERE CategoryName = 'Hospitality & Tourism'),
    (SELECT JobTypeID FROM JobTypes WHERE TypeName = 'Part-Time'),
    1, 1, 1200.00, 1600.00,
    '2026-08-31', 'Kuala Lumpur', 'KLCC', 4,
    'Looking for enthusiastic part-time baristas to join our flagship KLCC outlet. Flexible weekend and evening shifts available.',
    'Prepare and serve coffee and tea beverages to brand standards;
Maintain cleanliness of workstation and equipment;
Handle cashier duties and POS transactions;
Assist with stock-taking and restocking.',
    'No experience required - full training provided;
Friendly and customer-oriented personality;
Available to work at least 3 shifts per week including weekends.',
    'Free meals during shift;
Staff discount on all beverages.',
    'Walk-in interviews welcome at our KLCC outlet, or WhatsApp 012-3456789 to schedule.'
);
GO

PRINT 'Sample jobs inserted successfully.';
GO
