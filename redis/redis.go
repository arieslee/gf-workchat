// Package redis
package redis

import (
	"errors"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"github.com/gomodule/redigo/redis"
	"reflect"
	"sync"
)

type Client struct {
	Group string
}

const NoResult = "redigo: nil returned"

var singleton sync.Once
var instance *Client

// GetInstance 单例
func GetInstance(name ...string) *Client {
	group := "default"
	if len(name) > 0 {
		group = name[0]
	}
	singleton.Do(func() {
		instance = &Client{
			Group: group,
		}
	})
	instance.Group = group
	return instance
}
func (r *Client) IsEmpty(err error) bool {
	return err.Error() == NoResult
}

// ValueFormat 格式化值
func (r *Client) ValueFormat(value interface{}) (string, error) {
	kind := reflect.TypeOf(value).Kind()
	//  只特殊处理数组和MAP，还有BOOL类型的值
	if kind == reflect.Slice || kind == reflect.Map || kind == reflect.Array {
		str, err := gjson.Encode(value)
		if err != nil {
			return "", err
		}
		return gconv.String(str), nil
	} else if kind == reflect.Bool {
		var val string
		if value == true {
			val = "1"
		} else {
			val = "0"
		}
		return val, nil
	}
	return gconv.String(value), nil
}

// Set 设置
func (r *Client) Set(k string, v interface{}) (interface{}, error) {
	return g.Redis(r.Group).Do("SET", k, v)
}

// Get 获取
func (r *Client) Get(k string) (interface{}, error) {
	return g.Redis(r.Group).Do("GET", k)
}
func (r *Client) GetString(k string) (string, error) {
	value, err := g.Redis(r.Group).Do("GET", k)
	if err != nil {
		return "", err
	}
	str := gconv.String(value)
	return str, nil
}
func (r *Client) GetInt(k string) (int, error) {
	value, err := g.Redis(r.Group).Do("GET", k)
	if err != nil {
		return 0, err
	}
	return redis.Int(value, err)
}
func (r *Client) GetUint(k string) (uint, error) {
	value, err := r.GetString(k)
	if err != nil {
		return 0, err
	}
	num := gconv.Uint(value)
	return num, nil
}
func (r *Client) GetInt64(k string) (int64, error) {
	value, err := g.Redis(r.Group).Do("GET", k)
	if err != nil {
		return 0, err
	}
	return redis.Int64(value, err)
}
func (r *Client) GetUint64(k string) (uint64, error) {
	value, err := g.Redis(r.Group).Do("GET", k)
	if err != nil {
		return 0, err
	}
	return redis.Uint64(value, err)
}

// Del 删除
func (r *Client) Del(k string) (interface{}, error) {
	return g.Redis(r.Group).Do("DEL", k)
}

// SetEx 设置带有效期的值
func (r *Client) SetEx(k string, v interface{}, expire int) (interface{}, error) {
	return g.Redis(r.Group).Do("SETEX", k, expire, v)
}

// INCR 增加1
func (r *Client) INCR(k string) (int, error) {
	reply, errDo := g.Redis(r.Group).Do("INCR", k)
	if reply == nil && errDo == nil {
		return 0, nil
	}
	val, err := redis.Int(reply, errDo)
	return val, err
}
func (r *Client) INCRBY(k string, num int) (int, error) {
	reply, errDo := g.Redis(r.Group).Do("INCRBY", k, num)
	if reply == nil && errDo == nil {
		return 0, nil
	}
	val, err := redis.Int(reply, errDo)
	return val, err
}

// DECR 减少1
func (r *Client) DECR(k string) (int, error) {
	reply, errDo := g.Redis(r.Group).Do("DECR", k)
	if reply == nil && errDo == nil {
		return 0, nil
	}
	val, err := redis.Int(reply, errDo)
	return val, err
}
func (r *Client) DECRBY(k string, num int) (int, error) {
	reply, errDo := g.Redis(r.Group).Do("DECRBY", k, num)
	if reply == nil && errDo == nil {
		return 0, nil
	}
	val, err := redis.Int(reply, errDo)
	return val, err
}

// Expire 为某一个key指定过期时间
func (r *Client) Expire(key string, timeOutSeconds int64) (int64, error) {
	val, err := redis.Int64(g.Redis(r.Group).Do("EXPIRE", key, timeOutSeconds))
	return val, err
}

// ****************** hash set ***********************

// HGet hget
func (r *Client) HGet(k string, field string) (string, error) {
	reply, errDo := g.Redis(r.Group).Do("HGET", k, field)
	if errDo == nil && reply == nil {
		return "", nil
	}
	val, err := redis.String(reply, errDo)
	return val, err
}

// HGetAll hgetall
func (r *Client) HGetAll(k string) (map[string]string, error) {
	reply, err := redis.StringMap(g.Redis(r.Group).Do("HGETALL", k))
	return reply, err
}
func (r *Client) HGetMap(k string) (map[string]string, error) {
	return r.HGetAll(k)
}
func (r *Client) HSetMap(key string, data map[string]string) error {
	if data == nil || len(data) <= 0 {
		return errors.New("要缓存的数据不能为空")
	}
	for k, v := range data {
		err := r.HSet(key, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
func (r *Client) HSet(k string, field string, val string) error {
	_, err := g.Redis(r.Group).Do("HSET", k, field, val)
	return err
}
func (r *Client) HIncrBy(k string, field string, increment int) (int, error) {
	val, err := redis.Int(g.Redis(r.Group).Do("HINCRBY", k, field, increment))
	return val, err
}
func (r *Client) HDel(args ...interface{}) (int64, error) {
	val, err := redis.Int64(g.Redis(r.Group).Do("HDEL", args...))
	return val, err
}
func (r *Client) HLen(k string) (int64, error) {
	val, err := redis.Int64(g.Redis(r.Group).Do("HLEN", k))
	return val, err
}
func (r *Client) HExist(k string, field string) (int, error) {
	val, err := redis.Int(g.Redis(r.Group).Do("HEXISTS", k, field))
	return val, err
}
func (r *Client) HVals(k string) (interface{}, error) {
	val, err := redis.Strings(g.Redis(r.Group).Do("HVALS", k))
	return val, err
}
func (r *Client) HKeys(k string) (re []interface{}, err error) {
	keysData, err := g.Redis(r.Group).DoVar("HKEYS", k)
	if err != nil {
		return nil, err
	}
	return gconv.SliceAny(keysData), nil
}

// ****************** list ***********************

func (r *Client) LPush(key string, value ...interface{}) (int, error) {
	ret, err := redis.Int(g.Redis(r.Group).Do("LPUSH", key, value))
	if err != nil {
		return -1, err
	} else {
		return ret, nil
	}
}

func (r *Client) LPushX(key string, value string) (int, error) {
	resp, err := redis.Int(g.Redis(r.Group).Do("LPUSHX", key, value))
	return resp, err
}

func (r *Client) LRange(key string, start int, stop int) ([]string, error) {
	resp, err := redis.Strings(g.Redis(r.Group).Do("LRANGE", key, start, stop))
	return resp, err
}

func (r *Client) LRem(key string, count int, value string) (int, error) {
	resp, err := redis.Int(g.Redis(r.Group).Do("LREM", key, count, value))
	return resp, err
}

func (r *Client) LSet(key string, index int, value string) (string, error) {
	resp, err := redis.String(g.Redis(r.Group).Do("LSET", key, index, value))
	return resp, err
}

func (r *Client) LTrim(key string, start int, stop int) (string, error) {
	resp, err := redis.String(g.Redis(r.Group).Do("LTRIM", key, start, stop))
	return resp, err
}

func (r *Client) RPop(key string) (string, error) {
	resp, err := redis.String(g.Redis(r.Group).Do("RPOP", key))
	return resp, err
}

func (r *Client) RPush(key string, value ...interface{}) (int, error) {
	args := append([]interface{}{key}, value...)
	resp, err := redis.Int(g.Redis(r.Group).Do("RPUSH", args...))
	return resp, err
}

func (r *Client) RPushX(key string, value ...interface{}) (int, error) {
	args := append([]interface{}{key}, value...)
	resp, err := redis.Int(g.Redis(r.Group).Do("RPUSHX", args...))
	return resp, err
}

func (r *Client) RPopLPush(source string, destination string) (string, error) {
	resp, err := redis.String(g.Redis(r.Group).Do("RPOPLPUSH", source, destination))
	return resp, err
}

func (r *Client) BLPop(key ...interface{}) (map[string]string, error) {
	val, err := redis.StringMap(g.Redis(r.Group).Do("BLPOP", key))
	return val, err
}

func (r *Client) BRPop(key ...interface{}) (map[string]string, error) {
	val, err := redis.StringMap(g.Redis(r.Group).Do("BRPOP", key))
	return val, err
}

func (r *Client) BRPopLPush(source string, destination string) (string, error) {
	val, err := redis.String(g.Redis(r.Group).Do("BRPOPLPUSH", source, destination))
	return val, err
}

func (r *Client) LIndex(key string, index int) (string, error) {
	val, err := redis.String(g.Redis(r.Group).Do("LINDEX", key, index))
	return val, err
}

func (r *Client) LInsertBefore(key string, pivot string, value string) (int, error) {
	val, err := redis.Int(g.Redis(r.Group).Do("LINSERT", key, "BEFORE", pivot, value))
	return val, err
}

func (r *Client) LInsertAfter(key string, pivot string, value string) (int, error) {
	val, err := redis.Int(g.Redis(r.Group).Do("LINSERT", key, "AFTER", pivot, value))
	return val, err
}

func (r *Client) LLen(key string) (int, error) {
	val, err := redis.Int(g.Redis(r.Group).Do("LLEN", key))
	return val, err
}

func (r *Client) LPop(key string) (string, error) {
	val, err := redis.String(g.Redis(r.Group).Do("LPOP", key))
	return val, err
}

// ****************** SortedSet 有序集合 ***********************
//将一个 member 元素及其 score 值加入到有序集 key 当中
func (r *Client) ZAdd(key string, score float64, member string) (int, error) {
	ret, err := redis.Int(g.Redis(r.Group).Do("ZADD", key, score, member))
	if err != nil {
		return -1, err
	} else {
		return ret, nil
	}
}

//返回有序集 key 的基数
func (r *Client) ZCard(key string) (int, error) {
	ret, err := redis.Int(g.Redis(r.Group).Do("ZCARD", key))
	if err != nil {
		return -1, err
	} else {
		return ret, nil
	}
}

//返回有序集 key 中， score 值在 min 和 max 之间(默认包括 score 值等于 min 或 max )的成员的数量。
func (r *Client) ZCount(key string, min float64, max float64) (int, error) {
	ret, err := redis.Int(g.Redis(r.Group).Do("ZCOUNT", key, min, max))
	if err != nil {
		return -1, err
	} else {
		return ret, nil
	}
}

//为有序集 key 的成员 member 的 score 值加上增量 increment
func (r *Client) ZIncrBy(key string, increment float64, member string) (float64, error) {
	ret, err := redis.Float64(g.Redis(r.Group).Do("ZINCRBY", key, increment, member))
	if err != nil {
		return -1, err
	} else {
		return ret, nil
	}
}

//返回有序集 key 中，指定区间内的成员
func (r *Client) ZRange(key string, start int, stop int, isWithScores bool) (interface{}, error) {
	if isWithScores == false {
		redisRes, err := g.Redis(r.Group).Do("ZRANGE", key, start, stop)
		if err != nil {
			return nil, err
		}
		ret, err := redis.Strings(redisRes, err)
		if err != nil {
			return nil, err
		}
		return ret, nil
	}

	// WITHSCORES 选项，来让成员和它的 score 值一并返回
	ret, err := redis.StringMap(g.Redis(r.Group).Do("ZRANGE", key, start, stop, "WITHSCORES"))
	if err != nil {
		return nil, err
	}
	return ret, nil
}

//返回有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。有序集成员按 score 值递增(从小到大)次序排列。
func (r *Client) ZRangeByScore(key string, min string, max string, isWithScores bool, offset uint, count uint) (map[string]string, error) {
	if isWithScores == false {
		ret, err := redis.StringMap(g.Redis(r.Group).Do("ZRANGEBYSCORE", key, min, max))
		if err != nil {
			return nil, err
		} else {
			return ret, nil
		}
	}

	// WITHSCORES 选项，来让成员和它的 score 值一并返回
	ret, err := redis.StringMap(g.Redis(r.Group).Do("ZRANGEBYSCORE", key, min, max, "WITHSCORES"))
	if err != nil {
		return nil, err
	} else {
		return ret, nil
	}
}

//返回有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。有序集成员按 score 值递增(从小到大)次序排列 带分页
func (r *Client) ZRangeByScoreAndOffsetCount(key string, min string, max string, isWithScores bool, offset uint, count uint) (map[string]string, error) {
	if isWithScores == false {
		ret, err := redis.StringMap(g.Redis(r.Group).Do("ZRANGEBYSCORE", key, min, max, "LIMIT "+gconv.String(offset)+","+gconv.String(count)))
		if err != nil {
			return nil, err
		} else {
			return ret, nil
		}
	}

	// WITHSCORES 选项，来让成员和它的 score 值一并返回
	ret, err := redis.StringMap(g.Redis(r.Group).Do("ZRANGEBYSCORE", key, min, max, "WITHSCORES", "LIMIT "+gconv.String(offset)+","+gconv.String(count)))
	if err != nil {
		return nil, err
	} else {
		return ret, nil
	}
}

//返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递增(从小到大)顺序排列
//如果 member 是有序集 key 的成员，返回 member 的排名。
//如果 member 不是有序集 key 的成员，返回 nil 。
func (r *Client) ZRank(key string, member string) (int, error) {
	ret, err := redis.Int(g.Redis(r.Group).Do("ZRANK", key, member))
	if err != nil {
		return -1, err
	} else {
		return ret, nil
	}
}

//移除有序集 key 中的一个或多个成员，不存在的成员将被忽略
func (r *Client) ZRem(key string, member string) (int, error) {
	ret, err := redis.Int(g.Redis(r.Group).Do("ZREM", key, member))
	if err != nil {
		return -1, err
	} else {
		return ret, nil
	}
}

//移除有序集 key 中，指定排名(rank)区间内的所有成员
func (r *Client) ZRemRangeByRank(key string, start int, stop int, isWithScores bool) (interface{}, error) {
	if isWithScores == false {
		ret, err := redis.Int(g.Redis(r.Group).Do("ZREMRANGEBYRANK ", key, start, stop))
		if err != nil {
			return -1, err
		} else {
			return ret, nil
		}
	}

	ret, err := redis.StringMap(g.Redis(r.Group).Do("ZREMRANGEBYRANK ", key, start, stop, "WITHSCORES"))
	if err != nil {
		return -1, err
	} else {
		return ret, nil
	}
}

//移除有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员
func (r *Client) ZRemRangeByScore(key string, start int, stop int, isWithScores bool) (interface{}, error) {
	if isWithScores == false {
		ret, err := redis.Int(g.Redis(r.Group).Do("ZREMRANGEBYSCORE", key, start, stop))
		if err != nil {
			return -1, err
		} else {
			return ret, nil
		}
	}

	ret, err := redis.StringMap(g.Redis(r.Group).Do("ZREMRANGEBYSCORE", key, start, stop, "WITHSCORES"))
	if err != nil {
		return -1, err
	} else {
		return ret, nil
	}
}

//返回有序集 key 中，指定区间内的成员。
func (r *Client) ZRevRange(key string, start int, stop int, isWithScores bool) (interface{}, error) {

	if isWithScores == false {
		redisRes, err := g.Redis(r.Group).Do("ZREVRANGE", key, start, stop)
		if err != nil {
			return nil, err
		}
		ret, err := redis.Strings(redisRes, err)
		if err != nil {
			return nil, err
		}
		return ret, nil
	}

	// WITHSCORES 选项，来让成员和它的 score 值一并返回
	ret, err := redis.StringMap(g.Redis(r.Group).Do("ZREVRANGE", key, start, stop, "WITHSCORES"))
	if err != nil {
		return nil, err
	}
	return ret, nil
}

//返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递减(从大到小)排序。
//如果 member 是有序集 key 的成员，返回 member 的排名。
//如果 member 不是有序集 key 的成员，返回 nil 。
func (r *Client) ZRevRank(key string, member string) (int, error) {
	ret, err := redis.Int(g.Redis(r.Group).Do("ZREVRANK", key, member))
	if err != nil {
		return -1, err
	} else {
		return ret, nil
	}
}

//返回有序集 key 中，成员 member 的 score 值。
func (r *Client) ZScore(key string, member string) (map[string]string, error) {
	ret, err := redis.StringMap(g.Redis(r.Group).Do("ZSCORE", key, member))
	if err != nil {
		return nil, err
	} else {
		return ret, nil
	}
}

//ZUNIONSTORE todo
//ZINTERSTORE todo
//ZSCAN todo

// ****************** Global functions ***********************

func (r *Client) Ping() (string, error) {
	val, err := redis.String(g.Redis(r.Group).Do("PING"))
	return val, err
}

func (r *Client) DbSize() (int64, error) {
	val, err := redis.Int64(g.Redis(r.Group).Do("DBSIZE"))
	return val, err
}

func (r *Client) FlushDB() {
	g.Redis(r.Group).Do("FLUSHALL")
}
