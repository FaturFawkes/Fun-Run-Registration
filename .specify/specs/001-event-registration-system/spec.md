# Feature Specification: Event Registration System

**Feature Branch**: `001-event-registration-system`  
**Created**: 2025-12-31  
**Status**: Draft  
**Input**: User description: "Build a one-page event registration system for a Fun Run 5K called 'Tau-Tau Run'"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Public Participant Registration (Priority: P1)

A potential participant visits the Tau-Tau Run landing page, learns about the Fun Run 5K event, and registers by filling out a form with their personal details. Upon submission, their registration is stored in the database with PENDING status and UNPAID payment status.

**Why this priority**: Core business value. Without participant registration, the event has no attendees. This is the primary revenue and engagement driver.

**Independent Test**: Can be fully tested by visiting the public landing page, submitting a registration form with valid data, and verifying the participant appears in the database with correct initial states (PENDING/UNPAID).

**Acceptance Scenarios**:

1. **Given** a visitor is on the landing page, **When** they fill out the registration form with name, email, phone, Instagram handle, and address, **Then** their data is saved to the database with registration_status=PENDING and payment_status=UNPAID
2. **Given** a visitor submits incomplete form data, **When** they attempt to register, **Then** the system displays clear validation errors for missing required fields
3. **Given** a visitor enters an invalid email format, **When** they submit the form, **Then** the system rejects the submission with an email format error
4. **Given** a participant has already registered with an email, **When** they attempt to register again with the same email, **Then** the system prevents duplicate registration

---

### User Story 2 - Admin Authentication & Dashboard Access (Priority: P2)

An event administrator opens the admin dashboard, logs in with their credentials (email + password), and gains access to the participant management interface. The admin session is maintained securely.

**Why this priority**: Essential for event operations but depends on having participants to manage. Cannot operate the event without admin access to manage payments.

**Independent Test**: Can be fully tested by navigating to the admin login page, entering valid admin credentials, verifying successful authentication, and confirming access to the dashboard. Test with invalid credentials to verify security.

**Acceptance Scenarios**:

1. **Given** an admin is on the login page, **When** they enter valid email and password, **Then** they are authenticated and redirected to the admin dashboard
2. **Given** an admin enters invalid credentials, **When** they attempt to login, **Then** the system displays an authentication error and denies access
3. **Given** an unauthenticated user, **When** they attempt to access the admin dashboard directly, **Then** they are redirected to the login page
4. **Given** an admin is authenticated, **When** they close and reopen the dashboard, **Then** their session persists (within reasonable timeout)

---

### User Story 3 - Admin Payment Status Management (Priority: P2)

An authenticated admin views the list of all registered participants, sees their payment status, and updates a participant's payment status from UNPAID to PAID when they confirm payment has been received offline.

**Why this priority**: Critical for event operations and revenue tracking. Enables the admin to manage the payment workflow and maintain accurate records.

**Independent Test**: Can be fully tested by logging in as admin, viewing the participant list, selecting a participant with UNPAID status, updating them to PAID, and verifying the database reflects the change.

**Acceptance Scenarios**:

1. **Given** an admin is viewing the participant list, **When** they see a participant with UNPAID status, **Then** they can update the status to PAID via a clear UI action (button/toggle)
2. **Given** an admin updates a participant to PAID, **When** the update is submitted, **Then** the database immediately reflects payment_status=PAID for that participant
3. **Given** an admin views the participant list, **When** the page loads, **Then** all participants are displayed with their current name, email, phone, Instagram, address, and payment status
4. **Given** multiple participants exist, **When** an admin updates one participant's status, **Then** only that participant's status changes (no side effects on others)

---

### User Story 4 - Automatic Email Confirmation on Payment (Priority: P1)

When an admin marks a participant's payment status as PAID, the system automatically triggers an email to the participant confirming their registration and payment. The email is sent via SMTP without manual intervention.

**Why this priority**: Critical for participant experience and operational efficiency. Automates confirmation delivery and reduces manual work. Directly tied to payment processing (P2) but essential for MVP value.

**Independent Test**: Can be fully tested by updating a participant to PAID status and verifying that an email is received at the participant's registered email address. Test with multiple participants to ensure reliability.

**Acceptance Scenarios**:

1. **Given** a participant with UNPAID status, **When** an admin changes their status to PAID, **Then** an email confirmation is automatically sent to the participant's registered email address
2. **Given** a participant is updated to PAID, **When** the email is sent, **Then** the email contains confirmation of their registration and payment for the Tau-Tau Run Fun Run 5K
3. **Given** an email fails to send (SMTP error), **When** the payment status is updated, **Then** the system logs the error but still updates the payment status (payment status update and email sending are decoupled for resilience)
4. **Given** a participant is already marked PAID, **When** an admin views their status, **Then** no duplicate email is sent (email sent only on state transition from UNPAID to PAID)

---

### Edge Cases

- What happens when a participant submits a form with very long text in address fields (potential SQL injection or buffer overflow)?
- How does the system handle concurrent admin logins (multiple admins updating same participant)?
- What if SMTP server is unreachable when payment status is updated?
- How does the system prevent SQL injection attacks in registration form inputs?
- What happens if an admin tries to change payment status from PAID back to UNPAID?
- How does the system handle special characters or non-ASCII characters in participant names and addresses?
- What is the maximum number of participants the system can handle in a single event?
- How does the system handle browser refresh during registration form submission?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide a public landing page that explains the Tau-Tau Run Fun Run 5K event
- **FR-002**: System MUST allow participants to register by submitting name, email, phone number, Instagram handle, and address
- **FR-003**: System MUST validate all registration form inputs (email format, required fields, data types)
- **FR-004**: System MUST prevent duplicate registrations using the same email address
- **FR-005**: System MUST store participant data in PostgreSQL database with explicit registration_status (PENDING/CONFIRMED) and payment_status (UNPAID/PAID)
- **FR-006**: System MUST initialize new registrations with registration_status=PENDING and payment_status=UNPAID
- **FR-007**: System MUST provide an admin login page requiring email and password authentication
- **FR-008**: System MUST hash admin passwords using bcrypt (minimum 10 rounds)
- **FR-009**: System MUST provide an admin dashboard accessible only to authenticated admins
- **FR-010**: System MUST display a list of all registered participants in the admin dashboard
- **FR-011**: System MUST show participant details including name, email, phone, Instagram, address, and payment status
- **FR-012**: System MUST allow admins to update participant payment_status from UNPAID to PAID
- **FR-013**: System MUST automatically send email confirmation when payment_status changes to PAID
- **FR-014**: System MUST send email confirmation via SMTP (no external email service APIs)
- **FR-015**: System MUST include registration and payment confirmation details in the email
- **FR-016**: System MUST send email only once per UNPAID→PAID transition (no duplicate emails)
- **FR-017**: System MUST persist payment status change even if email sending fails
- **FR-018**: System MUST log email sending success/failure for operational visibility
- **FR-019**: System MUST provide a Golang REST API as the single source of truth for all data operations
- **FR-020**: Frontend MUST consume the API and contain no business logic

### Key Entities

- **Participant**: Represents a registered individual for the Fun Run 5K
  - Attributes: id, name, email, phone, instagram_handle, address, registration_status (PENDING/CONFIRMED), payment_status (UNPAID/PAID), created_at, updated_at
  - Email must be unique
  - All fields except Instagram handle are required

- **Admin**: Represents an authenticated event administrator
  - Attributes: id, email, password_hash, created_at
  - Single admin user for MVP (no role management)
  - Password stored as bcrypt hash

- **Email Event**: Represents the automatic email trigger
  - Not a stored entity, but a side effect when payment_status transitions to PAID
  - Triggered by backend state change, not manual admin action

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Participants can complete registration in under 2 minutes on the public landing page
- **SC-002**: Admin can log in and update payment status for a participant in under 30 seconds
- **SC-003**: Email confirmation is sent within 5 seconds of payment status update to PAID
- **SC-004**: System prevents 100% of duplicate email registrations
- **SC-005**: System handles at least 100 participants without performance degradation
- **SC-006**: All passwords are stored as bcrypt hashes with no plaintext exposure
- **SC-007**: Admin dashboard is inaccessible without valid authentication (0% unauthorized access)
- **SC-008**: Email sending failures do not block payment status updates (resilient architecture)
- **SC-009**: 100% of state transitions (UNPAID→PAID) are logged for audit purposes
- **SC-010**: Frontend contains 0 business logic (all logic in Golang backend)
