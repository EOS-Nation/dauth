---
--- Generated by EmmyLua(https://github.com/EmmyLua)
--- Created by cbillett.
--- DateTime: 2020-02-05 8:57 a.m.
---

redis.replicate_commands()

--- Utilities

local function slice(tbl, first, last, step)
    local sliced = {}

    for i = first or 1, last or #tbl, step or 1 do
        sliced[#sliced + 1] = tbl[i]
    end

    return sliced
end

local function computeCurrentPeriodConsumption()

    local keys = slice(KEYS, 5, 5 + windowSize)
    local values = redis.call("MGET", unpack(keys))
    local sum = 0
    for index, value in ipairs(values) do
        local v = tonumber(value)
        if v ~= nil then
            sum = sum + v
        end
    end
    return sum
end

--- Goes through the minutely consumed docs in reverse to figure out how long the user should be blocked
--[[local function computeBlockDuration(docQuota, windowSize)

    local keys = slice(KEYS, 5 + windowSize, 5, -1)
    local values = redis.call("MGET", unpack(keys))
    local sum = 0
    for index, value in ipairs(values) do
        local v = tonumber(value)
        if v ~= nil then
            sum = sum + v
            if sum >= docQuota then
                return windowSize - index
            end
        end
    end
    return 0
end]]

--- KEY DEFINITIONS
local userPeriodConsumedDoc = KEYS[1]
local blackListKey = KEYS[2]
local blackListVersionKey = KEYS[3]
local userDailyDocConsumptionKey = KEYS[4]
local userMinutelyDocConsumptionKey = KEYS[5]

--- SCRIPT INPUTS
local allocatedDocumentCount = tonumber(ARGV[1])
local consumeDocument = tonumber(ARGV[2])
local endOfWindowTimeStamp = tonumber(ARGV[3])
local blockDuration = tonumber(ARGV[4])
local windowSize = tonumber(ARGV[5])

--print("userPeriodConsumedDoc:" .. userPeriodConsumedDoc)
--print("allocatedDocumentCount:" .. allocatedDocumentCount)
--print("consumeDocument:" .. consumeDocument)
--print("endOfDayTimeStamp:" .. endOfDayTimeStamp)
--- --------------------------------------------------
--- The meat!
--- --------------------------------------------------
if consumeDocument == 0 then
    return 0
end

--- Daily consumption handling
if redis.call("INCRBY", userDailyDocConsumptionKey, consumeDocument) == consumeDocument then
    redis.call("EXPIRE", userDailyDocConsumptionKey, 86400)
end

--- Minutely consumption handling
if redis.call("INCRBY", userMinutelyDocConsumptionKey, consumeDocument) == consumeDocument then
    redis.call("EXPIRE", userMinutelyDocConsumptionKey, 60 * windowSize)
end

--- Skip paying customer
if allocatedDocumentCount == 0 then
    -- paying customer
    return "unlimited"
end

if redis.call("EXISTS", blackListKey) == 1 then
    return "already bl"
end

--- Remaining document handling
--[[local blockingDuration = computeBlockDuration(allocatedDocumentCount, 10)
if blockingDuration > 0 then
    local expireIn = blockingDuration * 60

    redis.call("SETEX", blackListKey, expireIn, "total doc exceeded")
    redis.call("INCR", blackListVersionKey)

    return "bl"
end]]

local periodConsumeDocument = redis.call("INCRBY", userPeriodConsumedDoc, consumeDocument)
if periodConsumeDocument == consumeDocument then
    -- We have a new userDailyDocConsumptionKey key
    local currentPeriodConsumption = computeCurrentPeriodConsumption()
    redis.call("SET", userPeriodConsumedDoc, currentPeriodConsumption)
    periodConsumeDocument = currentPeriodConsumption
    --local t = redis.call("TIME")
    --local expireIn = endOfWindowTimeStamp - tonumber(t[1])
    --redis.call("EXPIRE", userPeriodConsumedDoc, expireIn)
    redis.call("EXPIREAT", userPeriodConsumedDoc, endOfWindowTimeStamp)
end

if periodConsumeDocument > (allocatedDocumentCount * windowSize) then
    redis.call("SETEX", blackListKey, blockDuration, "total doc exceeded")
    redis.call("INCR", blackListVersionKey)
    --if redis.call("INCR", blackListVersionKey) == 1 then
    --    redis.call("EXPIRE", userDailyDocConsumptionKey, blockDuration)
    --end
    return "bl"
end

return "all good"
