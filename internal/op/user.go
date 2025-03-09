package op

import (
	"github.com/FXAZfung/go-cache"
	"github.com/FXAZfung/image-board/internal/db"
	"github.com/FXAZfung/image-board/internal/model"
	"github.com/FXAZfung/image-board/pkg/singleflight"
	"strconv"
	"time"
)

var userCache = cache.NewMemCache(cache.WithShards[*model.User](4))
var userListCache = cache.NewMemCache(cache.WithShards[interface{}](2))
var userG singleflight.Group[*model.User]
var userListG singleflight.Group[interface{}]

var guestUser *model.User
var adminUser *model.User

// UserCacheUpdate clears all user caches
func UserCacheUpdate() {
	userCache.Clear()
	userListCache.Clear()
}

// cacheUser helper to store user in cache
var cacheUser = func(user *model.User) {
	userCache.Set(user.Username, user, cache.WithEx[*model.User](time.Minute*10))
	userCache.Set(strconv.Itoa(int(user.ID)), user, cache.WithEx[*model.User](time.Minute*10))
	userCache.Set(strconv.Itoa(int(user.Role)), user, cache.WithEx[*model.User](time.Minute*10))
}

func GetAdmin() (*model.User, error) {
	if adminUser == nil {
		user, err := db.GetUserByRole(model.ADMIN)
		if err != nil {
			return nil, err
		}
		adminUser = user
	}
	return adminUser, nil
}

func GetGuest() (*model.User, error) {
	if guestUser == nil {
		user, err := db.GetUserByRole(model.GUEST)
		if err != nil {
			return nil, err
		}
		guestUser = user
	}
	return guestUser, nil
}

// GetUserById retrieves a user by their ID with caching
func GetUserById(id uint) (*model.User, error) {
	key := strconv.Itoa(int(id))
	if user, ok := userCache.Get(key); ok {
		return user, nil
	}

	user, err, _ := userG.Do(key, func() (*model.User, error) {
		user, err := db.GetUserById(id)
		if err != nil {
			return nil, err
		}
		cacheUser(user)
		return user, nil
	})
	return user, err
}

// GetUserByName retrieves a user by their username with caching
func GetUserByName(username string) (*model.User, error) {
	if user, ok := userCache.Get(username); ok {
		return user, nil
	}

	user, err, _ := userG.Do(username, func() (*model.User, error) {
		user, err := db.GetUserByName(username)
		if err != nil {
			return nil, err
		}
		cacheUser(user)
		return user, nil
	})
	return user, err
}

// GetUserByRole retrieves a user by their role with caching
func GetUserByRole(role int) (*model.User, error) {
	key := strconv.Itoa(role)
	if user, ok := userCache.Get(key); ok {
		return user, nil
	}

	user, err, _ := userG.Do("role_"+key, func() (*model.User, error) {
		user, err := db.GetUserByRole(role)
		if err != nil {
			return nil, err
		}
		cacheUser(user)
		return user, nil
	})
	return user, err
}

// GetUsers retrieves users with pagination and caching
func GetUsers(page, perPage int) ([]model.User, int64, error) {
	cacheKey := "users_page_" + strconv.Itoa(page) + "_" + strconv.Itoa(perPage)
	if cached, ok := userListCache.Get(cacheKey); ok {
		data := cached.(map[string]interface{})
		return data["users"].([]model.User), data["count"].(int64), nil
	}

	result, err, _ := userListG.Do(cacheKey, func() (interface{}, error) {
		users, count, err := db.GetUsers(page, perPage)
		if err != nil {
			return nil, err
		}

		// Cache individual users
		for i := range users {
			cacheUser(&users[i])
		}

		// Cache the page result
		data := map[string]interface{}{
			"users": users,
			"count": count,
		}
		userListCache.Set(cacheKey, data, cache.WithEx[interface{}](time.Minute*5))
		return data, nil
	})

	if result != nil {
		data := result.(map[string]interface{})
		return data["users"].([]model.User), data["count"].(int64), nil
	}
	return nil, 0, err
}

// CreateUser creates a new user with cache invalidation
func CreateUser(user *model.User) error {
	if err := db.CreateUser(user); err != nil {
		return err
	}

	// Update cache
	cacheUser(user)
	// Clear list cache since we added a new user
	userListCache.Clear()
	return nil
}

// UpdateUser updates a user's information and updates cache
func UpdateUser(user *model.User) error {
	if err := db.UpdateUser(user); err != nil {
		return err
	}

	// Update cache
	cacheUser(user)
	return nil
}

// DeleteUser deletes a user and removes from cache
func DeleteUser(id uint) error {
	// Get user to invalidate cache
	user, err := GetUserById(id)
	if err == nil {
		userCache.Del(user.Username)
		userCache.Del(strconv.Itoa(int(user.ID)))
		userCache.Del(strconv.Itoa(int(user.Role)))
	}

	if err := db.DeleteUser(id); err != nil {
		return err
	}

	// Clear list cache
	userListCache.Clear()
	return nil
}

// GetUserCount gets total user count with caching
func GetUserCount() (int64, error) {
	cacheKey := "user_count"
	if cached, ok := userListCache.Get(cacheKey); ok {
		return cached.(int64), nil
	}

	result, err, _ := userListG.Do(cacheKey, func() (interface{}, error) {
		count, err := db.GetUserCount()
		if err != nil {
			return int64(0), err
		}
		userListCache.Set(cacheKey, count, cache.WithEx[interface{}](time.Minute*5))
		return count, nil
	})

	if result != nil {
		return result.(int64), nil
	}
	return 0, err
}
