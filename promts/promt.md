# TASKBRIDGE - ROLE BASED PLATFORM DESIGN

You are the Lead Product Architect, Senior UX Designer, Senior Backend Architect and Full Stack Engineer for the TaskBridge platform.

Before doing anything:

Read and understand the existing project.

The project already contains:

- Authentication
- RBAC
- Users
- Categories
- Cities
- Tasks
- Applications
- Dashboard
- Search
- State Machine
- Wallet/Escrow architecture (Phase 5 in progress)

Your task is NOT to redesign the project.

Your task is to COMPLETE the platform by designing and implementing the complete role-based experience for:

1. Requester
2. Agent
3. Admin

Frontend and Backend must remain synchronized.

Every screen must have:

- UI
- Backend endpoint
- Permission rules
- Navigation placement
- State handling
- Empty states
- Loading states
- Error states

Generate everything required.

---

# PLATFORM CONCEPT

TaskBridge is a marketplace where:

Requester = Creates real-world tasks

Agent = Performs real-world tasks

Admin = Moderates and manages platform

Examples:

- University paperwork
- Government paperwork
- Medical appointments
- Document pickup
- Deliveries
- Administrative work

---

# ROLE SYSTEM

Implement complete separation between:

REQUESTER

AGENT

ADMIN

Each role must have:

- Different dashboard
- Different navigation
- Different permissions
- Different actions
- Different widgets
- Different statistics

Do not reuse the same dashboard.

---

==================================================
REQUESTER PANEL
==================================================

Goal:

Create tasks and manage them.

Navigation:

Dashboard

Create Task

My Tasks

Applications

Wallet

Notifications

Messages

Reviews

Profile

Settings

---

REQUESTER DASHBOARD

Widgets:

Active Tasks

Completed Tasks

Tasks Waiting Verification

Wallet Balance

Locked Balance

Latest Notifications

Latest Applications

Recent Activity

Task Statistics

Charts:

Tasks Created This Month

Tasks Completed This Month

Money Spent

---

REQUESTER FEATURES

Create Task

Edit Task

Delete Task

Cancel Task

View Applicants

Accept Applicant

Reject Applicant

View Task Progress

Chat With Agent

Verify Completion

Request Revision

Open Dispute

Rate Agent

Wallet Management

View Transactions

View Notifications

Manage Profile

---

TASK CREATION FORM

Fields:

Title

Description

Category

City

Address (optional)

Budget

Deadline

Priority

Attachments

Visibility

Validation Rules

Required Fields

Error States

Success States

---

TASK DETAIL PAGE

Show:

Task Information

Applications Count

Assigned Agent

Current Status

Timeline

Activity Log

Chat Shortcut

Attachments

Financial Summary

---

APPLICATION REVIEW PAGE

Show:

Agent Name

Rating

Completed Jobs

Success Rate

Verification Status

Trust Score

Proposal Message

Expected Completion Time

Accept Button

Reject Button

Chat Button

---

==================================================
AGENT PANEL
==================================================

Goal:

Find tasks and earn money.

Navigation:

Dashboard

Available Tasks

My Applications

Assigned Tasks

Completed Tasks

Wallet

Withdraw Requests

Messages

Reviews

Notifications

Profile

Settings

---

AGENT DASHBOARD

Widgets:

Available Tasks

Assigned Tasks

Tasks In Progress

Completed Tasks

Wallet Balance

Pending Withdrawals

Monthly Earnings

Total Earnings

Trust Score

Rating

Charts:

Earnings Trend

Completed Tasks Trend

Success Rate Trend

---

AGENT FEATURES

Browse Tasks

Search Tasks

Filter Tasks

Apply To Tasks

Withdraw Application

Track Application Status

Accept Assignment

Start Work

Submit Progress

Submit Completion

Upload Deliverables

Chat With Requester

Open Dispute

View Earnings

Withdraw Money

Rate Requester

Manage Profile

---

TASK MARKETPLACE PAGE

Show:

Title

Budget

Category

City

Deadline

Priority

Requester Rating

Trust Score

Application Count

Apply Button

---

APPLICATION FORM

Fields:

Proposal Message

Estimated Completion Time

Attachments

Expected Cost (future support)

Validation

---

ASSIGNED TASK PAGE

Show:

Requester

Task Details

Attachments

Timeline

Chat

Financial Summary

Milestones (future support)

Status Controls

---

==================================================
ADMIN PANEL
==================================================

Goal:

Manage and moderate the platform.

Navigation:

Dashboard

Users

Agents

Requesters

Tasks

Categories

Cities

Transactions

Wallets

Withdraw Requests

Revenue

Disputes

Reports

Notifications

Settings

Audit Logs

---

ADMIN DASHBOARD

Widgets:

Total Users

Active Users

Total Tasks

Active Tasks

Completed Tasks

Revenue

Pending Withdraws

Pending Verifications

Open Disputes

System Health

Charts:

Revenue

User Growth

Task Growth

Platform Activity

---

ADMIN FEATURES

Manage Users

Manage Agents

Manage Requesters

Manage Tasks

Manage Categories

Manage Cities

Manage Wallets

Manage Transactions

Manage Revenue

Approve Withdrawals

Reject Withdrawals

Manage Disputes

Manage Verification

Suspend Users

Ban Users

View Audit Logs

Send Platform Notifications

Generate Reports

---

USER MANAGEMENT PAGE

Show:

User Details

Verification Status

Wallet Balance

Completed Tasks

Complaints

Trust Score

Risk Score

Actions:

Approve

Suspend

Ban

Verify

Reset Account

---

DISPUTE MANAGEMENT

Admin can:

Review Evidence

Review Chat Messages

Review Attachments

Review Timeline

Decide Winner

Release Escrow

Refund Requester

Pay Agent

Close Dispute

---

==================================================
NOTIFICATION SYSTEM
==================================================

Create notification center.

Types:

Task Created

Task Assigned

Task Completed

Application Submitted

Application Accepted

Application Rejected

Payment Released

Withdraw Approved

Withdraw Rejected

Dispute Opened

System Announcement

---

==================================================
MESSAGING SYSTEM
==================================================

Task-based chat.

Each task has its own conversation.

Support:

Text

Image

File

PDF

Read Status

Delivery Status

Typing Indicator (future)

---

==================================================
REVIEW SYSTEM
==================================================

Requester rates Agent.

Agent rates Requester.

Fields:

Rating

Comment

CreatedAt

Verification

---

==================================================
TRUST SYSTEM
==================================================

Create Trust Score.

Factors:

Completed Tasks

Success Rate

Disputes

Verification Level

Reviews

Response Time

Show trust score throughout platform.

---

==================================================
PERMISSION MATRIX
==================================================

Generate complete permission matrix.

Example:

task.create

task.edit

task.delete

task.apply

task.accept

task.verify

wallet.view

wallet.withdraw

dispute.create

dispute.resolve

admin.users.manage

admin.transactions.manage

etc.

Map every permission to:

Requester

Agent

Admin

---

==================================================
BACKEND REQUIREMENTS
==================================================

Generate:

Routes

DTOs

Permissions

Policies

Services

Repository Interfaces

Swagger Documentation

Validation Rules

Database Requirements

Events

Caching Requirements

Rate Limits

---

==================================================
DELIVERABLES
==================================================

Generate:

1. Complete role-based UI structure

2. Sidebar/navigation structure for each role

3. Page hierarchy

4. User journeys

5. Permission matrix

6. Backend endpoint list

7. Missing database entities

8. Missing APIs

9. Missing services

10. Missing events

11. Missing notifications

12. Feature gap analysis

13. Recommendations for future phases

Update docs/report.md with all findings and additions.