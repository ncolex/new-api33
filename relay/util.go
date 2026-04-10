package relay

import (
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
)

func keyPrefixForLog(key string) string {
	trimmed := strings.TrimSpace(key)
	if trimmed == "" {
		return "********..."
	}
	if len(trimmed) <= 8 {
		return trimmed + "..."
	}
	return trimmed[:8] + "..."
}

func getChannelPoolKeys(channel *model.Channel) []string {
	if channel == nil {
		return nil
	}
	if len(channel.Keys) > 0 {
		keys := make([]string, 0, len(channel.Keys))
		for _, key := range channel.Keys {
			trimmed := strings.TrimSpace(key)
			if trimmed != "" {
				keys = append(keys, trimmed)
			}
		}
		return keys
	}
	if trimmed := strings.TrimSpace(channel.Key); trimmed != "" {
		return []string{trimmed}
	}
	return nil
}

func IsKeyCoolingDown(channel *model.Channel, key string) bool {
	if channel == nil || strings.TrimSpace(key) == "" {
		return false
	}
	value, ok := channel.CooldownKeys.Load(key)
	if !ok {
		return false
	}
	expireAt, ok := value.(time.Time)
	if !ok {
		channel.CooldownKeys.Delete(key)
		return false
	}
	if time.Now().After(expireAt) {
		channel.CooldownKeys.Delete(key)
		return false
	}
	return true
}

func SetKeyCooldown(channel *model.Channel, key string) {
	if channel == nil {
		return
	}
	trimmed := strings.TrimSpace(key)
	if trimmed == "" {
		return
	}
	cooldownSeconds := common.KeyCooldownSeconds
	if cooldownSeconds <= 0 {
		cooldownSeconds = 60
	}
	expireAt := time.Now().Add(time.Duration(cooldownSeconds) * time.Second)
	channel.CooldownKeys.Store(trimmed, expireAt)
	common.SysLog(fmt.Sprintf("Key %s del canal %d en cooldown por %d segundos", keyPrefixForLog(trimmed), channel.Id, cooldownSeconds))
}

func GetNextAvailableKey(channel *model.Channel) (string, error) {
	if channel == nil {
		return "", errors.New("channel_is_nil")
	}
	keys := getChannelPoolKeys(channel)
	if len(keys) == 0 {
		return "", errors.New("channel_key_empty")
	}
	for i := 0; i < len(keys); i++ {
		idx := int(atomic.AddInt32(&channel.RotationIndex, 1)-1) % len(keys)
		if idx < 0 {
			idx = 0
		}
		candidate := keys[idx]
		if !IsKeyCoolingDown(channel, candidate) {
			return candidate, nil
		}
	}
	return "", errors.New("all_keys_in_cooldown")
}
