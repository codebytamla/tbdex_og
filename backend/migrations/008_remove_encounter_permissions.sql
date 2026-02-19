-- +goose Up

-- Remove deprecated encounter permissions from existing databases.
DELETE FROM permissions
WHERE name LIKE 'encounters.%';
