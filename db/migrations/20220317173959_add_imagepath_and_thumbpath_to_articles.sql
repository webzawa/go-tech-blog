
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
alter table articles
  add column image_path mediumtext NOT NULL,
  add column thumb_path mediumtext NOT NULL;


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

alter table articles
  drop column image_path,
  drop column thumb_path;
