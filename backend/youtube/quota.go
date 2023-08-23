package main

import (
	"github.com/gofiber/fiber/v2/log"
	"time"
)

const quotaCacheKey string = "youtube-quota"

func getQuotaResetTime() time.Time {
	nowLocal := time.Now()
	nowPST := nowLocal.In(time.FixedZone("PST", -8*60*60))

	passedMidnight := time.Date(nowPST.Year(), nowPST.Month(), nowPST.Day(), 0, 0, 0, 0, nowPST.Location())
	nextMidnight := passedMidnight.AddDate(0, 0, 1)

	return nextMidnight
}

func (s *Service) addToQuota(usedPoints int64) int64 {
	currentlyUsed, err := s.rdb.IncrBy(quotaCacheKey, usedPoints).Result()
	if err != nil {
		log.Errorf("Failed to increment quota: %s", err)
		return 0
	}

	if err := s.rdb.ExpireAt(quotaCacheKey, getQuotaResetTime()).Err(); err != nil {
		log.Errorf("Failed to set expiry on key \"%s\": %s", quotaCacheKey, err)
	}

	return currentlyUsed
}
