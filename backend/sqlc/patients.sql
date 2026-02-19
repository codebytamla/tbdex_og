-- name: CreatePatient :one
INSERT INTO patients (
    full_name,
    date_of_birth,
    is_dob_estimated,
    sex,
    phone,
    national_id,
    mrn,
    address,
    village,
    origin,
    status,
    created_by
) VALUES (
    @full_name,
    @date_of_birth,
    @is_dob_estimated,
    @sex,
    @phone,
    @national_id,
    @mrn,
    @address,
    @village,
    @origin,
    @status,
    @created_by
)
RETURNING id, full_name, date_of_birth, is_dob_estimated, sex, phone,
          national_id, mrn, address, village, origin, status,
          created_by, created_at, updated_at;

-- name: GetPatientByID :one
SELECT id, full_name, date_of_birth, is_dob_estimated, sex, phone,
       national_id, mrn, address, village, origin, status,
       created_by, created_at, updated_at
FROM patients
WHERE id = @id;

-- name: ListPatients :many
SELECT p.id, p.full_name, p.date_of_birth, p.is_dob_estimated, p.sex, p.phone,
       p.national_id, p.mrn, p.address, p.village, p.origin, p.status,
       p.created_by, p.created_at, p.updated_at
FROM patients p
WHERE (
    @search::text = ''
    OR p.full_name ILIKE '%' || @search::text || '%' ESCAPE '\'
    OR p.mrn ILIKE '%' || @search::text || '%' ESCAPE '\'
    OR p.national_id ILIKE '%' || @search::text || '%' ESCAPE '\'
)
AND (
    NOT @use_status_filter::bool
    OR p.status = @status::text
)
AND (
    NOT @use_cursor::bool
    OR (p.created_at, p.id) < (@cursor_created_at::timestamptz, @cursor_id::uuid)
)
ORDER BY p.created_at DESC, p.id DESC
LIMIT @limit_rows;

-- name: CountPatients :one
SELECT COUNT(*)
FROM patients p
WHERE (
    @search::text = ''
    OR p.full_name ILIKE '%' || @search::text || '%' ESCAPE '\'
    OR p.mrn ILIKE '%' || @search::text || '%' ESCAPE '\'
    OR p.national_id ILIKE '%' || @search::text || '%' ESCAPE '\'
)
AND (
    NOT @use_status_filter::bool
    OR p.status = @status::text
);

-- name: UpdatePatient :one
UPDATE patients
SET
    full_name       = COALESCE(sqlc.narg('full_name'), full_name),
    date_of_birth   = COALESCE(sqlc.narg('date_of_birth'), date_of_birth),
    is_dob_estimated = COALESCE(sqlc.narg('is_dob_estimated'), is_dob_estimated),
    sex             = COALESCE(sqlc.narg('sex'), sex),
    phone           = COALESCE(sqlc.narg('phone'), phone),
    national_id     = COALESCE(sqlc.narg('national_id'), national_id),
    address         = COALESCE(sqlc.narg('address'), address),
    village         = COALESCE(sqlc.narg('village'), village),
    origin          = COALESCE(sqlc.narg('origin'), origin),
    status          = COALESCE(sqlc.narg('status'), status),
    updated_at      = now()
WHERE id = @id
RETURNING id, full_name, date_of_birth, is_dob_estimated, sex, phone,
          national_id, mrn, address, village, origin, status,
          created_by, created_at, updated_at;

-- name: DisablePatient :one
UPDATE patients
SET status = 'inactive', updated_at = now()
WHERE id = @id
RETURNING id, full_name, date_of_birth, is_dob_estimated, sex, phone,
          national_id, mrn, address, village, origin, status,
          created_by, created_at, updated_at;

-- name: GetPatientCreatedBy :one
SELECT created_by FROM patients WHERE id = @id;

-- name: UserIsAdmin :one
SELECT EXISTS (
    SELECT 1
    FROM user_roles ur
    JOIN roles r ON r.id = ur.role_id
    WHERE ur.user_id = @user_id AND r.name = 'ADMIN'
) AS is_admin;
