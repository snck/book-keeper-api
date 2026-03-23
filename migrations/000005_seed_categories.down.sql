DELETE FROM categories
WHERE user_id = '00000000-0000-0000-0000-000000000001'
  AND category_name IN ('Food', 'Commute', 'Health', 'Shopping', 'Entertainment');
