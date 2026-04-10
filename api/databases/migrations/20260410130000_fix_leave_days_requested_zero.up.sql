-- Backfill incorrect leave_requests.days_requested values (<= 0)
-- Rule:
-- - half-day records => 0.5
-- - full-day records => inclusive date diff (end_date - start_date + 1)
-- - guard invalid ranges by forcing minimum 1 for full-day

UPDATE leave_requests
SET days_requested = CASE
  WHEN COALESCE(is_half_day, false) = true THEN 0.5
  ELSE GREATEST((end_date::date - start_date::date) + 1, 1)
END
WHERE COALESCE(days_requested, 0) <= 0;
