# Tasks: Event Registration System

**Feature**: 001-event-registration-system  
**Branch**: `001-event-registration-system`  
**Input**: Design documents from `.specify/specs/001-event-registration-system/`  
**Prerequisites**: ✅ plan.md, ✅ spec.md, ✅ data-model.md, ✅ contracts/

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

**Timeline Estimate**: 3-5 days for solo developer (MVP scope)

---

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3, US4)
- File paths use web app structure: `backend/`, `frontend/`, `database/`

---

## Phase 1: Setup (Project Initialization) ✅ COMPLETE

**Purpose**: Create project structure and initialize dependencies

**Time Estimate**: 2-3 hours

- [x] T001 Create repository directory structure per plan.md (backend/, frontend/, database/, docs/)
- [x] T002 Initialize backend Go module in `backend/go.mod` with dependencies (gin, lib/pq, bcrypt, godotenv, jwt-go)
- [x] T003 [P] Initialize frontend Next.js project in `frontend/` with TypeScript and TailwindCSS
- [x] T004 [P] Create `.gitignore` files for backend and frontend (exclude .env files, node_modules, binaries)
- [x] T005 [P] Create `backend/.env.example` with all required environment variables
- [x] T006 [P] Create `frontend/.env.local.example` with API URL configuration
- [x] T007 [P] Create `docker-compose.yml` for local development (PostgreSQL, backend, frontend services)
- [x] T008 [P] Create `README.md` at repository root with setup instructions

**Checkpoint**: ✅ Project structure created, dependencies initialized

---

## Phase 2: Foundational Infrastructure (CRITICAL - Blocks All User Stories) ✅ COMPLETE

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**Time Estimate**: 4-6 hours

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### Database Setup

- [x] T009 Create PostgreSQL database schema in `database/migrations/001_init.sql` (participants, admins, email_logs tables)
- [x] T010 Create database seed script in `database/seeds/001_admin_seed.sql` (initial admin user)
- [x] T011 Create database connection module in `backend/internal/database/db.go` (connection pooling, migration support)
- [x] T012 Test database connection and verify schema creation works

### Backend Core Infrastructure

- [x] T013 [P] Create configuration management in `backend/config/config.go` (load from .env, validate required vars)
- [x] T014 [P] Create error handling middleware in `backend/internal/middleware/error.go` (JSON error responses)
- [x] T015 [P] Create logging utility in `backend/internal/utils/logger.go` (structured logging)
- [x] T016 [P] Create CORS middleware in `backend/internal/middleware/cors.go` (configure allowed origins)
- [x] T017 Create main server entrypoint in `backend/cmd/server/main.go` (Gin setup, route registration, graceful shutdown)

### Authentication Infrastructure

- [x] T018 Create JWT utility in `backend/internal/services/auth.go` (GenerateToken, ValidateToken, HashPassword, ComparePassword)
- [x] T019 Create JWT authentication middleware in `backend/internal/middleware/auth.go` (extract token, validate, attach admin to context)
- [x] T020 Test JWT generation and validation with sample data

### Frontend Core Infrastructure

- [x] T021 [P] Create API client utility in `frontend/src/services/api.ts` (base URL, error handling, request/response interceptors)
- [x] T022 [P] Create TypeScript types in `frontend/src/types/index.ts` (Participant, Admin, API responses)
- [x] T023 [P] Setup TailwindCSS configuration in `frontend/tailwind.config.ts` (color palette, typography)
- [x] T024 Create base layout in `frontend/src/app/layout.tsx` (global styles, metadata)

**Checkpoint**: ✅ Foundation ready - database connected, server runs, auth works, frontend scaffolded

---

## Phase 3: User Story 1 - Public Participant Registration (Priority: P1) 🎯 MVP ✅ COMPLETE

**Goal**: Allow visitors to register for the event via public landing page

**Independent Test**: Visit landing page, submit registration form, verify participant appears in database with PENDING/UNPAID status

**Time Estimate**: 6-8 hours

### Backend Implementation

- [x] T025 [P] Create Participant model in `backend/internal/models/participant.go` (struct, validation, CRUD methods)
- [x] T026 [P] Create validation utility in `backend/internal/utils/validator.go` (email format, phone format, required fields)
- [x] T027 Create participant registration handler in `backend/internal/handlers/participant.go` (POST /register endpoint)
- [x] T028 Add participant routes to main server in `backend/cmd/server/main.go` (public router group)
- [x] T029 Implement duplicate email check in participant handler
- [x] T030 Add input sanitization to prevent SQL injection
- [x] T031 Add structured logging for registration events

### Frontend Implementation

- [x] T032 [P] Create RegistrationForm component in `frontend/src/components/RegistrationForm.tsx` (form fields, validation, submission)
- [x] T033 [P] Create form validation logic with client-side hints (email format, required fields, character limits)
- [x] T034 Create public landing page in `frontend/src/app/page.tsx` (event info, registration form, success/error messages)
- [x] T035 Add API integration to RegistrationForm (call POST /register, handle responses)
- [x] T036 Style landing page with TailwindCSS (poster-like, flat design per constitution)
- [x] T037 Add loading states and error handling to registration form

### Testing & Validation

- [x] T038 Manual test: Register participant with valid data → verify 201 response and database entry
- [x] T039 Manual test: Register with duplicate email → verify 409 error
- [x] T040 Manual test: Register with invalid email → verify 400 validation error
- [x] T041 Manual test: Submit incomplete form → verify validation errors displayed

**Checkpoint**: ✅ User Story 1 complete - public registration fully functional and tested

---

## Phase 4: User Story 2 - Admin Authentication & Dashboard Access (Priority: P2) ✅ COMPLETE

**Goal**: Admin can log in and access dashboard to view participants

**Independent Test**: Navigate to admin login, enter credentials, verify successful authentication and redirect to dashboard

**Time Estimate**: 4-5 hours

### Backend Implementation

- [x] T042 [P] Create Admin model in `backend/internal/models/admin.go` (struct, FindByEmail method)
- [x] T043 Create admin login handler in `backend/internal/handlers/admin.go` (POST /login endpoint)
- [x] T044 Add admin routes to main server (admin router group with /login public, others protected)
- [x] T045 Implement login logic (verify email exists, compare password hash, generate JWT)
- [x] T046 Add login attempt logging for security audit

### Frontend Implementation

- [x] T047 [P] Create admin login page in `frontend/src/app/admin/login/page.tsx` (email/password form)
- [x] T048 [P] Create authentication context/hook for token management (store in localStorage or cookies)
- [x] T049 Add login form submission logic (call POST /login, store token, redirect to dashboard)
- [x] T050 Add error handling for invalid credentials
- [x] T051 Style admin login page (simple, professional, distinct from public page)

### Testing & Validation

- [x] T052 Manual test: Login with valid credentials → verify token received and redirect to dashboard
- [x] T053 Manual test: Login with invalid credentials → verify 401 error displayed
- [x] T054 Manual test: Access dashboard without login → verify redirect to login page

**Checkpoint**: ✅ User Story 2 complete - admin authentication works, dashboard accessible

---

## Phase 5: User Story 3 - Admin Payment Status Management (Priority: P2)

**Goal**: Admin can view participant list and update payment status

**Independent Test**: Login as admin, view participant list, update one participant to PAID, verify database reflects change

**Time Estimate**: 5-6 hours

### Backend Implementation

- [ ] T055 Create list participants handler in `backend/internal/handlers/admin.go` (GET /participants endpoint)
- [ ] T056 Create update payment status handler in `backend/internal/handlers/admin.go` (PATCH /participants/:id/payment endpoint)
- [ ] T057 Add participant routes to admin router group (both require JWT middleware)
- [ ] T058 Implement payment status validation (only PAID or UNPAID allowed)
- [ ] T059 Add updated_at timestamp update on payment status change
- [ ] T060 Add logging for payment status changes (audit trail)

### Frontend Implementation

- [ ] T061 [P] Create ParticipantList component in `frontend/src/components/ParticipantList.tsx` (table with participant details)
- [ ] T062 [P] Create PaymentStatusToggle component in `frontend/src/components/PaymentStatusToggle.tsx` (button/toggle to update status)
- [ ] T063 Create admin dashboard page in `frontend/src/app/admin/dashboard/page.tsx` (participant list with filters)
- [ ] T064 Add API integration for fetching participants (GET /participants with JWT header)
- [ ] T065 Add API integration for updating payment status (PATCH /participants/:id/payment)
- [ ] T066 Implement optimistic UI updates (update UI immediately, rollback on error)
- [ ] T067 Add loading and error states for dashboard operations
- [ ] T068 Style dashboard with table layout and clear payment status indicators

### Testing & Validation

- [ ] T069 Manual test: Login and view participant list → verify all participants displayed
- [ ] T070 Manual test: Update participant status to PAID → verify UI updates and database reflects change
- [ ] T071 Manual test: Attempt to update with invalid status → verify error handling
- [ ] T072 Manual test: Refresh dashboard → verify participant list loads correctly

**Checkpoint**: User Story 3 complete - admin can manage payment status

---

## Phase 6: User Story 4 - Automatic Email Confirmation on Payment (Priority: P1)

**Goal**: System automatically sends confirmation email when payment status changes to PAID

**Independent Test**: Update participant to PAID status, verify email received at participant's email address

**Time Estimate**: 4-5 hours

### Backend Implementation

- [ ] T073 [P] Create email service in `backend/internal/services/email.go` (SMTP connection, SendConfirmationEmail function)
- [ ] T074 [P] Create email template builder in `backend/internal/services/email.go` (HTML/plain text confirmation email)
- [ ] T075 Add email sending logic to payment update handler (trigger ONLY on UNPAID → PAID transition)
- [ ] T076 Implement async email sending with goroutine (non-blocking, don't wait for SMTP response)
- [ ] T077 Add email logging to `email_logs` table (record SUCCESS or FAILED status)
- [ ] T078 Implement email error handling (log error but don't fail payment update)
- [ ] T079 Add idempotency check (don't send duplicate emails if status is already PAID)

### Configuration & Testing

- [ ] T080 Add SMTP configuration to `backend/.env.example` (host, port, username, password, from address)
- [ ] T081 Create email configuration validator (verify SMTP settings on startup)
- [ ] T082 Manual test: Update participant to PAID → verify email sent and logged as SUCCESS
- [ ] T083 Manual test: Simulate SMTP failure → verify payment still updated and error logged
- [ ] T084 Manual test: Update PAID participant to PAID again → verify no duplicate email sent
- [ ] T085 Test email content and formatting (verify all participant details included)

**Checkpoint**: User Story 4 complete - email automation works reliably

---

## Phase 7: Integration & End-to-End Testing

**Purpose**: Verify all user stories work together seamlessly

**Time Estimate**: 3-4 hours

- [ ] T086 End-to-end test: Register participant → verify appears in database with UNPAID status
- [ ] T087 End-to-end test: Login as admin → view participant list → verify new participant appears
- [ ] T088 End-to-end test: Update participant to PAID → verify email sent and status updated
- [ ] T089 End-to-end test: Logout admin → verify cannot access dashboard
- [ ] T090 Test error scenarios: Network failures, database down, SMTP timeout
- [ ] T091 Test edge cases: Special characters in names, long addresses, multiple concurrent updates
- [ ] T092 Verify all API endpoints return consistent JSON response format
- [ ] T093 Verify frontend handles all API error codes gracefully
- [ ] T094 Test mobile responsiveness of public page and admin dashboard
- [ ] T095 Verify no business logic exists in frontend (all validations happen in backend)

**Checkpoint**: Full system integration verified

---

## Phase 8: Documentation & Polish

**Purpose**: Finalize documentation and prepare for deployment

**Time Estimate**: 2-3 hours

- [ ] T096 [P] Create API documentation in `docs/API.md` (all endpoints with examples)
- [ ] T097 [P] Create deployment guide in `docs/DEPLOYMENT.md` (production setup, environment variables)
- [ ] T098 [P] Update repository README.md with project overview and quick start
- [ ] T099 [P] Document SMTP configuration in quickstart.md (Gmail, Mailtrap, SendGrid examples)
- [ ] T100 [P] Add inline code comments for complex logic (email trigger, JWT validation)
- [ ] T101 Create database backup and restore procedures documentation
- [ ] T102 Add security checklist to deployment guide (HTTPS, JWT secret, bcrypt cost)
- [ ] T103 Verify all environment variables are documented in .env.example files
- [ ] T104 Code cleanup: Remove console.logs, fix linting issues, format code
- [ ] T105 Run through quickstart.md on fresh environment to verify accuracy

**Checkpoint**: Documentation complete, ready for deployment

---

## Phase 9: Deployment Preparation (Optional - Production Ready)

**Purpose**: Prepare application for production deployment

**Time Estimate**: 3-4 hours

- [ ] T106 [P] Create production Dockerfile for backend in `backend/Dockerfile`
- [ ] T107 [P] Create production Dockerfile for frontend in `frontend/Dockerfile`
- [ ] T108 Create production docker-compose.yml (optimized, no dev dependencies)
- [ ] T109 Setup PostgreSQL production configuration (connection limits, SSL mode)
- [ ] T110 Configure HTTPS/TLS for production backend (reverse proxy setup)
- [ ] T111 Setup environment variable management for production (secrets manager or encrypted .env)
- [ ] T112 Add health check endpoints for monitoring (backend and database connectivity)
- [ ] T113 Configure logging for production (log levels, log rotation, centralized logging)
- [ ] T114 Setup database backup automation (daily backups, retention policy)
- [ ] T115 Perform security audit (SQL injection, XSS, CSRF, rate limiting)
- [ ] T116 Load testing with 100+ concurrent users (verify performance targets)
- [ ] T117 Deploy to staging environment and run full test suite
- [ ] T118 Create rollback plan and disaster recovery procedures

**Checkpoint**: Application production-ready

---

## Dependencies & Execution Order

### Phase Dependencies

1. **Setup (Phase 1)**: No dependencies - start immediately
2. **Foundational (Phase 2)**: Depends on Setup - **BLOCKS ALL USER STORIES**
3. **User Stories (Phase 3-6)**: All depend on Foundational completion
   - Can proceed in parallel (if multiple developers)
   - Or sequentially by priority: US1 (P1) → US4 (P1) → US2 (P2) → US3 (P2)
4. **Integration (Phase 7)**: Depends on all user stories being complete
5. **Documentation (Phase 8)**: Can happen in parallel with user stories
6. **Deployment (Phase 9)**: Depends on Integration completion

### User Story Dependencies

- **US1 - Public Registration (P1)**: Can start after Foundational - No dependencies on other stories
- **US2 - Admin Authentication (P2)**: Can start after Foundational - Independent from US1
- **US3 - Payment Management (P2)**: Requires US2 complete (needs admin auth) and US1 complete (needs participants)
- **US4 - Email Automation (P1)**: Requires US3 complete (payment status update triggers email)

### Recommended Execution Order (Solo Developer)

1. **Phase 1 + Phase 2**: Setup and Foundation (6-9 hours)
2. **Phase 3 (US1)**: Public registration (6-8 hours) → **First MVP checkpoint**
3. **Phase 4 (US2)**: Admin auth (4-5 hours)
4. **Phase 5 (US3)**: Payment management (5-6 hours)
5. **Phase 6 (US4)**: Email automation (4-5 hours) → **Complete MVP**
6. **Phase 7**: Integration testing (3-4 hours)
7. **Phase 8**: Documentation (2-3 hours)
8. **Phase 9**: Deployment (3-4 hours) - optional, do when ready

**Total Estimated Time**: 33-44 hours (4-5 full working days)

### Parallel Opportunities (Team of 2-3)

Once **Phase 2 (Foundational)** is complete:

- **Developer A**: US1 (Public Registration) → US4 (Email Automation)
- **Developer B**: US2 (Admin Auth) → US3 (Payment Management)
- **Developer C**: Documentation (Phase 8) + Frontend styling

**Team Time Savings**: Reduce to 2-3 days instead of 4-5 days

---

## Implementation Strategy

### MVP First (Minimum Viable Product)

**Goal**: Get to a working demo as fast as possible

1. ✅ Complete Phase 1: Setup (2-3 hours)
2. ✅ Complete Phase 2: Foundational (4-6 hours)
3. ✅ Complete Phase 3: User Story 1 - Public Registration (6-8 hours)
4. **🛑 STOP and VALIDATE**: Can participants register? Test thoroughly.
5. ✅ Complete Phase 4 + 5 + 6: Admin features and email (13-16 hours)
6. **🛑 STOP and VALIDATE**: Full flow working? Register → Admin → Pay → Email
7. ✅ Deploy basic version

**First milestone**: Public registration working (1-2 days)  
**Second milestone**: Complete MVP with admin and email (3-4 days)

### Incremental Delivery Strategy

1. **Foundation** (Phase 1 + 2) → Database and server running
2. **Public Registration** (Phase 3) → Participants can register ✅ **Demo-able**
3. **Admin Login** (Phase 4) → Admin can authenticate ✅ **Demo-able**
4. **Payment Management** (Phase 5) → Admin can update status ✅ **Demo-able**
5. **Email Automation** (Phase 6) → Full workflow automated ✅ **Production ready**
6. **Polish** (Phase 7 + 8) → Tested and documented ✅ **Deployable**

Each phase adds value without breaking previous functionality.

---

## Testing Checklist (Manual Testing - MVP)

### User Story 1: Public Registration
- [ ] ✅ Register with valid data → Success message, database entry created
- [ ] ✅ Register with duplicate email → Error: "Email already registered"
- [ ] ✅ Register with invalid email → Error: "Valid email required"
- [ ] ✅ Submit incomplete form → Validation errors displayed
- [ ] ✅ Special characters in name/address → Properly stored and displayed

### User Story 2: Admin Authentication
- [ ] ✅ Login with valid credentials → Token received, redirect to dashboard
- [ ] ✅ Login with invalid email → Error: "Invalid email or password"
- [ ] ✅ Login with wrong password → Error: "Invalid email or password"
- [ ] ✅ Access dashboard without login → Redirect to login page
- [ ] ✅ Token expires after 24 hours → Auto-logout or refresh prompt

### User Story 3: Payment Management
- [ ] ✅ View participant list as admin → All participants displayed
- [ ] ✅ Update UNPAID to PAID → UI updates, database updated
- [ ] ✅ Update PAID to PAID → No changes, no duplicate actions
- [ ] ✅ Multiple admins updating same participant → Last update wins (acceptable for MVP)

### User Story 4: Email Automation
- [ ] ✅ Update to PAID → Email sent within 5 seconds
- [ ] ✅ Email contains correct participant details
- [ ] ✅ Email failure → Payment still updated, error logged
- [ ] ✅ Update PAID to PAID again → No duplicate email sent
- [ ] ✅ Check email_logs table → All attempts logged

### Integration
- [ ] ✅ Complete flow: Register → Login → Update → Email received
- [ ] ✅ Multiple participants → All managed correctly
- [ ] ✅ Error handling → No crashes, graceful error messages

---

## Notes

- **[P] tasks**: Can run in parallel (different files, no blocking dependencies)
- **[Story] labels**: Map tasks to user stories for traceability
- **Checkpoints**: Stop and validate before proceeding
- **Constitution compliance**: All tasks align with MVP-first, state-driven principles
- **Commit strategy**: Commit after each task or logical group of tasks
- **Branch strategy**: Work on feature branch `001-event-registration-system`, merge to main when complete

---

## Quick Reference: Task Counts by Phase

| Phase | Task Count | Estimated Time |
|-------|------------|----------------|
| Phase 1: Setup | 8 tasks | 2-3 hours |
| Phase 2: Foundational | 16 tasks | 4-6 hours |
| Phase 3: US1 (Public Registration) | 17 tasks | 6-8 hours |
| Phase 4: US2 (Admin Auth) | 13 tasks | 4-5 hours |
| Phase 5: US3 (Payment Management) | 18 tasks | 5-6 hours |
| Phase 6: US4 (Email Automation) | 13 tasks | 4-5 hours |
| Phase 7: Integration Testing | 10 tasks | 3-4 hours |
| Phase 8: Documentation | 10 tasks | 2-3 hours |
| Phase 9: Deployment (Optional) | 13 tasks | 3-4 hours |
| **TOTAL** | **118 tasks** | **33-44 hours** |

---

**Document Status**: ✅ Ready for implementation  
**Next Step**: Begin Phase 1 (Setup) - Create project structure  
**First Milestone Goal**: Complete Phase 3 (Public Registration working)  
**MVP Completion Goal**: Complete Phase 6 (All user stories functional)
