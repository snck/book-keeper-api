DELETE FROM categories
WHERE user_id IS NULL
  AND category_name IN ('Food', 'Commute', 'Health', 'Shopping', 'Entertainment');
