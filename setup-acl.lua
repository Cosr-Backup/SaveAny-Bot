-- Redis ACL 用户设置脚本 / Redis ACL User Setup Script
-- 用于为 SaveAny-Bot 创建专用的 ACL 用户 / Used to create dedicated ACL user for SaveAny-Bot

-- 创建 SaveAny-Bot 专用用户 / Create dedicated user for SaveAny-Bot
-- 权限：可以执行读写操作，但不能执行管理命令 / Permissions: Can execute read/write operations, but not admin commands
redis.call('ACL', 'SETUSER', 'saveany-app', 'on', 
    '>saveany-bot-password',  -- 设置密码 / Set password
    '+@read',                 -- 允许所有读操作 / Allow all read operations
    '+@write',                -- 允许所有写操作 / Allow all write operations  
    '+@string',               -- 允许字符串操作 / Allow string operations
    '+@hash',                 -- 允许哈希操作 / Allow hash operations
    '+@list',                 -- 允许列表操作 / Allow list operations
    '+@set',                  -- 允许集合操作 / Allow set operations
    '+@sortedset',            -- 允许有序集合操作 / Allow sorted set operations
    '+@stream',               -- 允许流操作 / Allow stream operations
    '+@bitmap',               -- 允许位图操作 / Allow bitmap operations
    '+@hyperloglog',          -- 允许 HyperLogLog 操作 / Allow HyperLogLog operations
    '+@geo',                  -- 允许地理位置操作 / Allow geo operations
    '+ping',                  -- 允许 ping 命令 / Allow ping command
    '+info',                  -- 允许 info 命令 / Allow info command
    '+select',                -- 允许选择数据库 / Allow database selection
    '+expire',                -- 允许设置过期时间 / Allow setting expiration
    '+ttl',                   -- 允许查询 TTL / Allow TTL queries
    '+exists',                -- 允许检查键存在 / Allow key existence checks
    '+del',                   -- 允许删除键 / Allow key deletion
    '+flushdb',               -- 允许清空当前数据库 / Allow flushing current database
    '~*'                      -- 允许访问所有键 / Allow access to all keys
)

-- 创建只读用户 (可选) / Create read-only user (optional)
redis.call('ACL', 'SETUSER', 'saveany-readonly', 'on',
    '>readonly-password',     -- 设置密码 / Set password
    '+@read',                 -- 只允许读操作 / Only allow read operations
    '+ping',                  -- 允许 ping 命令 / Allow ping command
    '+info',                  -- 允许 info 命令 / Allow info command
    '+select',                -- 允许选择数据库 / Allow database selection
    '~*'                      -- 允许访问所有键 / Allow access to all keys
)

-- 显示用户列表 / Show user list
return redis.call('ACL', 'LIST')