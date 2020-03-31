package lua

// for Redis Lua Script
// for middleware ipLimitIntercept
// running script will be atomic
const SCRIPT = `
local key = KEYS[1]
local now = tonumber(ARGV[1])
local ipLimit = tonumber(ARGV[2])
local period = tonumber(ARGV[3])
local userInfo = redis.call('HGETALL', key)
local reset = tonumber(userInfo[4])
local result = {}
if #userInfo == 0 or reset < now then
    reset = now + period
    redis.call('HMSET', key, "count", 1, "reset", reset)
    result[1] = ipLimit - 1
    result[2] = reset
    return result
end

local count = tonumber(userInfo[2])
if count < ipLimit then
    local newCount = redis.call('HINCRBY', key, "count", 1)	
    result[1] = ipLimit - newCount
    result[2] = reset
    return result
else
    result[1] = -1
    result[2] = reset
    return result
end
`
