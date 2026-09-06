-- ============================================================
--  媒体数据一致性修复脚本（幂等，可重复执行）
--
--  背景：媒体管理此前上传统一走本地磁盘，url 落相对路径 /uploads/<path>、
--        source 落 local。改造接入 MinIO 后，为保证 DB 中
--        source / path / url 三者与实际存储、展示链路一致，
--        本脚本对存量数据进行归一化清洗（不改变文件实体）。
--
--  覆盖三类问题：
--   1) source 为空/非法           → 按 url 推断 local / external
--   2) url 为后端内部绝对地址      → 归一化为相对 /uploads/<path>
--   3) path 为空但 url 含 /uploads → 从 url 反推 path
--
--  执行方式：
--     docker exec -i sportswear-postgres psql -U postgres -d sportswear_platform < align-media-data.sql
--  或：
--     psql "$DATABASE_URL" -f align-media-data.sql
-- ============================================================
SET client_encoding TO 'UTF8';

-- 1. source 归一化：空/非法值按 url 推断（仅处理未删除记录）
UPDATE media
SET source = CASE
    WHEN url LIKE 'http://%' OR url LIKE 'https://%'
         THEN 'external'
    ELSE 'local'
END
WHERE deleted_at IS NULL
  AND (source IS NULL OR source = '' OR source NOT IN ('local','minio','external'));

-- 2. url 归一化：把后端内部绝对地址转为相对 /uploads/<path>
--    覆盖 localhost / 127.0.0.1 / backend / host.docker.internal 四种形态
UPDATE media
SET url = REGEXP_REPLACE(
    url,
    '^https?://(localhost|127\.0\.0\.1|backend|host\.docker\.internal)(:[0-9]+)?/uploads/',
    '/uploads/'
)
WHERE deleted_at IS NULL
  AND url ~ '^https?://(localhost|127\.0\.0\.1|backend|host\.docker\.internal)(:[0-9]+)?/uploads/';

-- 3. path 反推：url 为 /uploads/<path> 但 path 为空时，从 url 补 path
UPDATE media
SET path = SUBSTRING(url FROM 10)  -- '/uploads/' 共 9 字符，从第 10 位起为对象路径
WHERE deleted_at IS NULL
  AND (path IS NULL OR path = '')
  AND url LIKE '/uploads/%';

-- ============================================================
-- 4. 校验查询（只读，不修改数据）：执行后若结果为空即表示已一致
-- ============================================================
-- 4.1 source 仍为空/非法的记录数
SELECT 'source_invalid' AS check, count(*) AS cnt
FROM media
WHERE deleted_at IS NULL
  AND (source IS NULL OR source = '' OR source NOT IN ('local','minio','external'));

-- 4.2 url 仍为后端内部绝对地址的记录数
SELECT 'url_absolute' AS check, count(*) AS cnt
FROM media
WHERE deleted_at IS NULL
  AND url ~ '^https?://(localhost|127\.0\.0\.1|backend|host\.docker\.internal)(:[0-9]+)?/uploads/';

-- 4.3 非 external 但 path 为空（无法定位物理对象）的记录
SELECT 'path_empty' AS check, count(*) AS cnt
FROM media
WHERE deleted_at IS NULL
  AND source IN ('local','minio')
  AND (path IS NULL OR path = '');
