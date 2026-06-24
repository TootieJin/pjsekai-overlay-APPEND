package scoremaker

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/vmihailenco/msgpack/v5"
)

type AuthData struct {
	DataVersion  string
	AssetVersion string
	AssetHash    string
}

type publishedScoreKeys struct {
	key []byte
	iv  []byte
}

var publishedScoreEncryptionKeys = map[string]publishedScoreKeys{
	"jp": {
		key: []byte("g2fcC0ZczN9MTJ61"),
		iv:  []byte("msx3IV0i9XE5uYZ1"),
	},
	"en": {
		key: []byte{0xdf, 0x38, 0x42, 0x14, 0xb2, 0x9a, 0x3a, 0xdf, 0xbf, 0x1b, 0xd9, 0xee, 0x5b, 0x16, 0xf8, 0x84},
		iv:  []byte{0x7e, 0x85, 0x6c, 0x90, 0x79, 0x87, 0xf8, 0xae, 0xc6, 0xaf, 0xc0, 0xc5, 0x47, 0x38, 0xfc, 0x7e},
	},
}

type PublishedScoreInfo struct {
	UserCustomMusicScoreID    string  `json:"userCustomMusicScoreId" msgpack:"userCustomMusicScoreId"`
	UserID                    int     `json:"userId" msgpack:"userId"`
	UserName                  string  `json:"userName" msgpack:"userName"`
	MusicID                   int     `json:"musicId" msgpack:"musicId"`
	Title                     string  `json:"title" msgpack:"title"`
	UserCustomMusicScorePath  string  `json:"userCustomMusicScorePath" msgpack:"userCustomMusicScorePath"`
	MusicDifficultyType       string  `json:"musicDifficultyType" msgpack:"musicDifficultyType"`
	PlayLevel                 int     `json:"playLevel" msgpack:"playLevel"`
	CreatedAt                 int64   `json:"createdAt" msgpack:"createdAt"`
	UpdatedAt                 int64   `json:"updatedAt" msgpack:"updatedAt"`
	Description               string  `json:"description" msgpack:"description"`
	FavoriteCount             int     `json:"favoriteCount" msgpack:"favoriteCount"`
	Version                   int     `json:"version" msgpack:"version"`
	IsPublished               bool    `json:"isPublished" msgpack:"isPublished"`
	PublishedAt               int64   `json:"publishedAt" msgpack:"publishedAt"`
	PlayCount                 int     `json:"playCount" msgpack:"playCount"`
	ReviewCount               int     `json:"reviewCount" msgpack:"reviewCount"`
	FullComboRate             float64 `json:"fullComboRate" msgpack:"fullComboRate"`
	CustomMusicScoreSortValue int     `json:"customMusicScoreSearchSortValue" msgpack:"customMusicScoreSearchSortValue"`
	CustomMusicScoreTags      []int   `json:"customMusicScoreTags" msgpack:"customMusicScoreTags"`
	IsDerivativeAllowed       bool    `json:"isDerivativeAllowed" msgpack:"isDerivativeAllowed"`
	IsReviewAllowed           bool    `json:"isReviewAllowed" msgpack:"isReviewAllowed"`
	IsReviewed                bool    `json:"isReviewed" msgpack:"isReviewed"`
	PreviewStartTimeSec       int     `json:"previewStartTimeSec" msgpack:"previewStartTimeSec"`
}

type publishedScoreInner struct {
	MusicID                  int    `json:"musicId" msgpack:"musicId"`
	Title                    string `json:"title" msgpack:"title"`
	UserCustomMusicScorePath string `json:"userCustomMusicScorePath" msgpack:"userCustomMusicScorePath"`
}

type PublishedScoreInfoEnvelope struct {
	UserCustomMusicScoreInfoJson struct {
		UserCustomMusicScoreInfoJson publishedScoreInner `json:"userCustomMusicScoreInfoJson" msgpack:"userCustomMusicScoreInfoJson"`
		UserCustomMusicScoreID       int                 `json:"userCustomMusicScoreId" msgpack:"userCustomMusicScoreId"`
		UserID                       int                 `json:"userId" msgpack:"userId"`
		MusicDifficultyType          string              `json:"musicDifficultyType" msgpack:"musicDifficultyType"`
		PlayLevel                    int                 `json:"playLevel" msgpack:"playLevel"`
	} `json:"userCustomMusicScoreInfoJson" msgpack:"userCustomMusicScoreInfoJson"`
}

func toInt(v any) int {
	switch t := v.(type) {
	case int:
		return t
	case int8:
		return int(t)
	case int16:
		return int(t)
	case int32:
		return int(t)
	case int64:
		return int(t)
	case uint:
		return int(t)
	case uint8:
		return int(t)
	case uint16:
		return int(t)
	case uint32:
		return int(t)
	case uint64:
		return int(t)
	case float32:
		return int(t)
	case float64:
		return int(t)
	case string:
		i, _ := strconv.Atoi(strings.TrimSpace(t))
		return i
	default:
		return 0
	}
}

func toString(v any) string {
	switch t := v.(type) {
	case string:
		return strings.TrimSpace(t)
	case nil:
		return ""
	default:
		return strings.TrimSpace(fmt.Sprint(t))
	}
}

func toMap(v any) map[string]any {
	if m, ok := v.(map[string]any); ok {
		return m
	}
	return map[string]any{}
}

func toInt64(v any) int64 {
	switch t := v.(type) {
	case int:
		return int64(t)
	case int8:
		return int64(t)
	case int16:
		return int64(t)
	case int32:
		return int64(t)
	case int64:
		return t
	case uint:
		return int64(t)
	case uint8:
		return int64(t)
	case uint16:
		return int64(t)
	case uint32:
		return int64(t)
	case uint64:
		return int64(t)
	case float32:
		return int64(t)
	case float64:
		return int64(t)
	case string:
		i, _ := strconv.ParseInt(strings.TrimSpace(t), 10, 64)
		return i
	default:
		return 0
	}
}

func toBool(v any) bool {
	switch t := v.(type) {
	case bool:
		return t
	case int:
		return t != 0
	case int8:
		return t != 0
	case int16:
		return t != 0
	case int32:
		return t != 0
	case int64:
		return t != 0
	case uint:
		return t != 0
	case uint8:
		return t != 0
	case uint16:
		return t != 0
	case uint32:
		return t != 0
	case uint64:
		return t != 0
	case float32:
		return t != 0
	case float64:
		return t != 0
	case string:
		s := strings.ToLower(strings.TrimSpace(t))
		return s == "true" || s == "1" || s == "yes"
	default:
		return false
	}
}

func toFloat64(v any) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case float32:
		return float64(t)
	case int:
		return float64(t)
	case int8:
		return float64(t)
	case int16:
		return float64(t)
	case int32:
		return float64(t)
	case int64:
		return float64(t)
	case uint:
		return float64(t)
	case uint8:
		return float64(t)
	case uint16:
		return float64(t)
	case uint32:
		return float64(t)
	case uint64:
		return float64(t)
	case string:
		f, _ := strconv.ParseFloat(strings.TrimSpace(t), 64)
		return f
	default:
		return 0
	}
}

func toIntSlice(v any) []int {
	switch t := v.(type) {
	case []any:
		out := make([]int, 0, len(t))
		for _, item := range t {
			out = append(out, toInt(item))
		}
		return out
	case []int:
		out := make([]int, len(t))
		copy(out, t)
		return out
	default:
		return nil
	}
}

func findValueByKeys(v any, keys ...string) (any, bool) {
	if len(keys) == 0 {
		return nil, false
	}
	keyset := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		k = strings.ToLower(strings.TrimSpace(k))
		if k != "" {
			keyset[k] = struct{}{}
		}
	}
	if len(keyset) == 0 {
		return nil, false
	}

	var walk func(any) (any, bool)
	walk = func(node any) (any, bool) {
		switch t := node.(type) {
		case map[string]any:
			for k, val := range t {
				if _, ok := keyset[strings.ToLower(strings.TrimSpace(k))]; ok {
					return val, true
				}
			}
			for _, val := range t {
				if got, ok := walk(val); ok {
					return got, true
				}
			}
		case []any:
			for _, item := range t {
				if got, ok := walk(item); ok {
					return got, true
				}
			}
		case string:
			s := strings.TrimSpace(t)
			if s == "" {
				return nil, false
			}
			if strings.HasPrefix(s, "{") || strings.HasPrefix(s, "[") {
				var decoded any
				if json.Unmarshal([]byte(s), &decoded) == nil {
					if got, ok := walk(decoded); ok {
						return got, true
					}
				}
			}
		}
		return nil, false
	}

	return walk(v)
}

func extractPublishedScoreInfo(raw []byte) (PublishedScoreInfo, error) {
	var root map[string]any
	if err := msgpack.Unmarshal(raw, &root); err != nil {
		if jsonErr := json.Unmarshal(raw, &root); jsonErr != nil {
			return PublishedScoreInfo{}, err
		}
	}

	outer := toMap(root["userCustomMusicScoreInfoJson"])
	if len(outer) == 0 {
		outer = toMap(root["userCustomMusicScoreInfo"])
	}
	if len(outer) == 0 {
		outer = toMap(root["data"])
	}
	if len(outer) == 0 {
		outer = root
	}
	inner := toMap(outer["userCustomMusicScoreInfoJson"])
	if len(inner) == 0 {
		inner = toMap(outer["userCustomMusicScoreInfo"])
	}

	info := PublishedScoreInfo{
		UserCustomMusicScoreID:    toString(outer["userCustomMusicScoreId"]),
		UserID:                    toInt(outer["userId"]),
		UserName:                  toString(outer["userName"]),
		MusicID:                   toInt(inner["musicId"]),
		Title:                     toString(inner["title"]),
		UserCustomMusicScorePath:  toString(inner["userCustomMusicScorePath"]),
		MusicDifficultyType:       toString(outer["musicDifficultyType"]),
		PlayLevel:                 toInt(outer["playLevel"]),
		CreatedAt:                 toInt64(outer["createdAt"]),
		UpdatedAt:                 toInt64(outer["updatedAt"]),
		Description:               toString(outer["description"]),
		FavoriteCount:             toInt(outer["favoriteCount"]),
		Version:                   toInt(outer["version"]),
		IsPublished:               toBool(outer["isPublished"]),
		PublishedAt:               toInt64(outer["publishedAt"]),
		PlayCount:                 toInt(outer["playCount"]),
		ReviewCount:               toInt(outer["reviewCount"]),
		FullComboRate:             toFloat64(outer["fullComboRate"]),
		CustomMusicScoreSortValue: toInt(outer["customMusicScoreSearchSortValue"]),
		CustomMusicScoreTags:      toIntSlice(outer["customMusicScoreTags"]),
		IsDerivativeAllowed:       toBool(outer["isDerivativeAllowed"]),
		IsReviewAllowed:           toBool(outer["isReviewAllowed"]),
		IsReviewed:                toBool(outer["isReviewed"]),
		PreviewStartTimeSec:       toInt(outer["previewStartTimeSec"]),
	}

	if info.UserCustomMusicScorePath == "" {
		info.UserCustomMusicScorePath = toString(outer["userCustomMusicScorePath"])
	}
	if info.UserCustomMusicScorePath == "" {
		info.UserCustomMusicScorePath = toString(root["userCustomMusicScorePath"])
	}
	if info.MusicID == 0 {
		info.MusicID = toInt(outer["musicId"])
	}
	if info.MusicID == 0 {
		info.MusicID = toInt(root["musicId"])
	}
	if info.Title == "" {
		info.Title = toString(outer["title"])
	}
	if info.Title == "" {
		info.Title = toString(root["title"])
	}
	if info.UserID == 0 {
		info.UserID = toInt(root["userId"])
	}
	if info.UserName == "" {
		info.UserName = toString(inner["userName"])
	}
	if info.UserName == "" {
		info.UserName = toString(root["userName"])
	}
	if info.UserCustomMusicScoreID == "" {
		info.UserCustomMusicScoreID = toString(root["userCustomMusicScoreId"])
	}
	if info.UserCustomMusicScoreID == "" {
		if v, ok := findValueByKeys(root,
			"userCustomMusicScoreId",
			"publishedUserCustomMusicScoreId",
			"customMusicScoreId",
		); ok {
			info.UserCustomMusicScoreID = toString(v)
		}
	}
	if info.MusicDifficultyType == "" {
		info.MusicDifficultyType = toString(inner["musicDifficultyType"])
	}
	if info.MusicDifficultyType == "" {
		info.MusicDifficultyType = toString(root["musicDifficultyType"])
	}
	if info.PlayLevel == 0 {
		info.PlayLevel = toInt(inner["playLevel"])
	}
	if info.PlayLevel == 0 {
		info.PlayLevel = toInt(root["playLevel"])
	}

	// additional metadata fields with fallbacks
	if info.CreatedAt == 0 {
		info.CreatedAt = toInt64(inner["createdAt"])
	}
	if info.CreatedAt == 0 {
		info.CreatedAt = toInt64(root["createdAt"])
	}
	if info.CreatedAt == 0 {
		if v, ok := findValueByKeys(root, "createdAt", "publishedAt", "postAt", "publishedAtMs"); ok {
			info.CreatedAt = toInt64(v)
		}
	}

	if info.UpdatedAt == 0 {
		info.UpdatedAt = toInt64(inner["updatedAt"])
	}
	if info.UpdatedAt == 0 {
		info.UpdatedAt = toInt64(root["updatedAt"])
	}
	if info.UpdatedAt == 0 {
		if v, ok := findValueByKeys(root, "updatedAt", "modifiedAt", "updatedAtMs"); ok {
			info.UpdatedAt = toInt64(v)
		}
	}
	if info.UpdatedAt == 0 {
		if v, ok := findValueByKeys(root, "publishedAt", "postAt", "publishedAtMs"); ok {
			info.UpdatedAt = toInt64(v)
		}
	}
	if info.PublishedAt == 0 {
		if v, ok := findValueByKeys(root, "publishedAt", "postAt", "publishedAtMs"); ok {
			info.PublishedAt = toInt64(v)
		}
	}

	if info.Description == "" {
		info.Description = toString(inner["description"])
	}
	if info.Description == "" {
		info.Description = toString(root["description"])
	}

	if info.FavoriteCount == 0 {
		info.FavoriteCount = toInt(inner["favoriteCount"])
	}
	if info.FavoriteCount == 0 {
		info.FavoriteCount = toInt(root["favoriteCount"])
	}
	if info.FavoriteCount == 0 {
		if v, ok := findValueByKeys(root,
			"favoriteCount",
			"favouriteCount",
			"favoritesCount",
			"likesCount",
		); ok {
			info.FavoriteCount = toInt(v)
		}
	}

	if info.Version == 0 {
		info.Version = toInt(inner["version"])
	}
	if info.Version == 0 {
		info.Version = toInt(root["version"])
	}
	if info.Version == 0 {
		if v, ok := findValueByKeys(root, "version", "revision", "seq", "generation"); ok {
			info.Version = toInt(v)
		}
	}

	if !info.IsPublished {
		info.IsPublished = toBool(inner["isPublished"])
	}
	if !info.IsPublished {
		info.IsPublished = toBool(root["isPublished"])
	}
	if !info.IsPublished {
		if v, ok := findValueByKeys(root, "isPublished", "published"); ok {
			info.IsPublished = toBool(v)
		}
	}

	if info.PlayCount == 0 {
		if v, ok := findValueByKeys(root, "playCount"); ok {
			info.PlayCount = toInt(v)
		}
	}
	if info.ReviewCount == 0 {
		if v, ok := findValueByKeys(root, "reviewCount"); ok {
			info.ReviewCount = toInt(v)
		}
	}
	if info.FullComboRate == 0 {
		if v, ok := findValueByKeys(root, "fullComboRate"); ok {
			info.FullComboRate = toFloat64(v)
		}
	}
	if info.CustomMusicScoreSortValue == 0 {
		if v, ok := findValueByKeys(root, "customMusicScoreSearchSortValue"); ok {
			info.CustomMusicScoreSortValue = toInt(v)
		}
	}
	if len(info.CustomMusicScoreTags) == 0 {
		if v, ok := findValueByKeys(root, "customMusicScoreTags"); ok {
			info.CustomMusicScoreTags = toIntSlice(v)
		}
	}
	if !info.IsDerivativeAllowed {
		if v, ok := findValueByKeys(root, "isDerivativeAllowed"); ok {
			info.IsDerivativeAllowed = toBool(v)
		}
	}
	if !info.IsReviewAllowed {
		if v, ok := findValueByKeys(root, "isReviewAllowed"); ok {
			info.IsReviewAllowed = toBool(v)
		}
	}
	if !info.IsReviewed {
		if v, ok := findValueByKeys(root, "isReviewed"); ok {
			info.IsReviewed = toBool(v)
		}
	}
	if info.PreviewStartTimeSec == 0 {
		if v, ok := findValueByKeys(root, "previewStartTimeSec"); ok {
			info.PreviewStartTimeSec = toInt(v)
		}
	}

	if info.MusicID == 0 {
		return PublishedScoreInfo{}, fmt.Errorf("published score info missing musicId")
	}
	return info, nil
}

func getEncryptionSet(region string) publishedScoreKeys {
	region = strings.ToLower(region)
	if region == "en" {
		return publishedScoreEncryptionKeys["en"]
	}
	return publishedScoreEncryptionKeys["jp"]
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, fmt.Errorf("invalid padded data length")
	}
	padLen := int(data[len(data)-1])
	if padLen == 0 || padLen > blockSize || padLen > len(data) {
		return nil, fmt.Errorf("invalid padding")
	}
	for _, b := range data[len(data)-padLen:] {
		if int(b) != padLen {
			return nil, fmt.Errorf("invalid padding")
		}
	}
	return data[:len(data)-padLen], nil
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padLen := blockSize - (len(data) % blockSize)
	if padLen == 0 {
		padLen = blockSize
	}
	padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(data, padding...)
}

func decryptPayload(ciphertext []byte, region string) ([]byte, error) {
	keys := getEncryptionSet(region)
	block, err := aes.NewCipher(keys.key)
	if err != nil {
		return nil, err
	}
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext length must be a multiple of block size")
	}
	plaintext := make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(block, keys.iv).CryptBlocks(plaintext, ciphertext)
	return pkcs7Unpad(plaintext, aes.BlockSize)
}

func encryptPayload(data any, region string) ([]byte, error) {
	packed, err := msgpack.Marshal(data)
	if err != nil {
		return nil, err
	}
	padded := pkcs7Pad(packed, aes.BlockSize)

	keys := getEncryptionSet(region)
	block, err := aes.NewCipher(keys.key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, len(padded))
	cipher.NewCBCEncrypter(block, keys.iv).CryptBlocks(ciphertext, padded)
	return ciphertext, nil
}

func newDeviceID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uint32(b[0])<<24|uint32(b[1])<<16|uint32(b[2])<<8|uint32(b[3]),
		uint16(b[4])<<8|uint16(b[5]),
		uint16(b[6])<<8|uint16(b[7]),
		uint16(b[8])<<8|uint16(b[9]),
		uint64(b[10])<<40|uint64(b[11])<<32|uint64(b[12])<<24|uint64(b[13])<<16|uint64(b[14])<<8|uint64(b[15]),
	), nil
}

func extractSessionToken(payload any) string {
	switch v := payload.(type) {
	case map[string]any:
		for k, val := range v {
			if strings.EqualFold(k, "sessionToken") {
				if s, ok := val.(string); ok {
					return s
				}
			}
			if nested := extractSessionToken(val); nested != "" {
				return nested
			}
		}
	case []any:
		for _, item := range v {
			if nested := extractSessionToken(item); nested != "" {
				return nested
			}
		}
	}
	return ""
}

func RegisterGuestAccount(region string, appInfo AppInfo) (int, string, error) {
	if strings.TrimSpace(appInfo.AppVersion) == "" || strings.TrimSpace(appInfo.AppHash) == "" {
		return 0, "", fmt.Errorf("app info is required")
	}
	baseURL, ok := apiEndpoints[strings.ToLower(region)]
	if !ok {
		return 0, "", fmt.Errorf("unsupported region: %s", region)
	}

	body, err := encryptPayload(map[string]any{
		"platform":        defaultPlatform,
		"deviceModel":     defaultDeviceModel,
		"operatingSystem": defaultOperatingSystem,
	}, region)
	if err != nil {
		return 0, "", err
	}

	issueCookie, err := FetchIssueCookie(region, nil)
	if err != nil {
		return 0, "", fmt.Errorf("failed to fetch issue cookie: %w", err)
	}

	headers, err := BuildHeaders(DownloadConfig{
		Region:      region,
		AppVersion:  appInfo.AppVersion,
		AppHash:     appInfo.AppHash,
		IssueCookie: issueCookie,
	})
	if err != nil {
		return 0, "", err
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"/api/user", bytes.NewReader(body))
	if err != nil {
		return 0, "", err
	}
	req.Header = headers

	client := &http.Client{}
	if strings.TrimSpace(issueCookie) != "" {
		client.Transport = &cookieRoundTripper{base: http.DefaultTransport, cookie: issueCookie}
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", err
	}
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return 0, "", fmt.Errorf("register request failed: status=%d %s", resp.StatusCode, decodePublishedErrorBody(raw, region))
	}

	plain, err := decryptPayload(raw, region)
	if err != nil {
		return 0, "", fmt.Errorf("failed to decrypt register response: %w", err)
	}

	var decoded map[string]any
	if err := msgpack.Unmarshal(plain, &decoded); err != nil {
		if jsonErr := json.Unmarshal(plain, &decoded); jsonErr != nil {
			return 0, "", fmt.Errorf("failed to parse register response: %w", err)
		}
	}

	// api response shape: { "userRegistration": { "userId": ..., "credential": ... }, ... }
	// credential may also appear at the top level
	userID := 0
	credential := ""

	if reg, ok := decoded["userRegistration"].(map[string]any); ok {
		userID = toInt(reg["userId"])
		credential = toString(reg["credential"])
	}
	if credential == "" {
		credential = toString(decoded["credential"])
	}
	if userID == 0 {
		userID = toInt(decoded["userId"])
	}

	if userID == 0 {
		return 0, "", fmt.Errorf("register response missing userId")
	}
	if credential == "" {
		return 0, "", fmt.Errorf("register response missing credential")
	}

	return userID, credential, nil
}

func AuthenticateWithCredential(userID int, credential string, region string, appInfo AppInfo) (string, AuthData, error) {
	var zero AuthData
	if userID <= 0 {
		return "", zero, fmt.Errorf("user id is required")
	}
	if strings.TrimSpace(credential) == "" {
		return "", zero, fmt.Errorf("credential is required")
	}
	if strings.TrimSpace(appInfo.AppVersion) == "" || strings.TrimSpace(appInfo.AppHash) == "" {
		return "", zero, fmt.Errorf("app info is required")
	}

	baseURL, ok := apiEndpoints[strings.ToLower(region)]
	if !ok {
		return "", zero, fmt.Errorf("unsupported region: %s", region)
	}

	deviceID, err := newDeviceID()
	if err != nil {
		return "", zero, err
	}
	body, err := encryptPayload(map[string]any{
		"credential":      credential,
		"deviceId":        deviceID,
		"authTriggerType": "normal",
	}, region)
	if err != nil {
		return "", zero, err
	}

	issueCookie, err := FetchIssueCookie(region, nil)
	if err != nil {
		return "", zero, fmt.Errorf("failed to fetch issue cookie: %w", err)
	}

	headers, err := BuildHeaders(DownloadConfig{
		Region:      region,
		AppVersion:  appInfo.AppVersion,
		AppHash:     appInfo.AppHash,
		IssueCookie: issueCookie,
	})
	if err != nil {
		return "", zero, err
	}

	url := fmt.Sprintf("%s/api/user/%d/auth", baseURL, userID)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return "", zero, err
	}
	req.Header = headers

	client := &http.Client{}
	if strings.TrimSpace(issueCookie) != "" {
		client.Transport = &cookieRoundTripper{base: http.DefaultTransport, cookie: issueCookie}
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", zero, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", zero, err
	}
	if resp.StatusCode != http.StatusOK {
		return "", zero, fmt.Errorf("authenticate request failed: status=%d %s", resp.StatusCode, decodePublishedErrorBody(raw, region))
	}

	headerToken := strings.TrimSpace(resp.Header.Get("X-Session-Token"))

	plain, decErr := decryptPayload(raw, region)
	if decErr != nil {
		if headerToken != "" {
			return headerToken, zero, nil
		}
		return "", zero, decErr
	}

	var root map[string]any
	if err := msgpack.Unmarshal(plain, &root); err != nil {
		_ = json.Unmarshal(plain, &root)
	}

	token := headerToken
	if token == "" {
		var decoded any
		if err := msgpack.Unmarshal(plain, &decoded); err == nil {
			token = extractSessionToken(decoded)
		}
		if token == "" {
			var decodedJSON any
			if err := json.Unmarshal(plain, &decodedJSON); err == nil {
				token = extractSessionToken(decodedJSON)
			}
		}
	}
	if token == "" {
		return "", zero, fmt.Errorf("authenticate succeeded but session token was not found in response")
	}

	var authData AuthData
	if root != nil {
		authData.DataVersion = toString(root["dataVersion"])
		authData.AssetVersion = toString(root["assetVersion"])
		authData.AssetHash = toString(root["assetHash"])
	}

	return token, authData, nil
}

func decodePublishedErrorBody(raw []byte, region string) string {
	tryDecode := func(payload []byte) string {
		var asMap map[string]any
		if err := msgpack.Unmarshal(payload, &asMap); err == nil {
			if code, ok := asMap["errorCode"].(string); ok {
				if msg, ok := asMap["message"].(string); ok && strings.TrimSpace(msg) != "" {
					return fmt.Sprintf("errorCode=%s message=%s", code, msg)
				}
				return fmt.Sprintf("errorCode=%s", code)
			}
			if asJSON, err := json.Marshal(asMap); err == nil {
				return string(asJSON)
			}
		}

		asText := strings.TrimSpace(string(payload))
		if asText != "" {
			return asText
		}
		return "<binary body>"
	}

	if decrypted, err := decryptPayload(raw, region); err == nil {
		if decoded := tryDecode(decrypted); decoded != "" {
			return decoded
		}
	}
	return tryDecode(raw)
}

func requestPublishedScoreInfo(
	method string,
	url string,
	region string,
	sessionToken string,
	authData AuthData,
	appInfo AppInfo,
) ([]byte, string, error) {
	issueCookie, err := FetchIssueCookie(region, nil)
	if err != nil {
		return nil, sessionToken, fmt.Errorf("failed to fetch issue cookie: %w", err)
	}

	headers, err := BuildHeaders(DownloadConfig{
		Region:       region,
		AppVersion:   appInfo.AppVersion,
		AppHash:      appInfo.AppHash,
		IssueCookie:  issueCookie,
		SessionToken: sessionToken,
	})
	if err != nil {
		return nil, "", err
	}
	if strings.TrimSpace(authData.DataVersion) != "" {
		headers.Set("X-Data-Version", authData.DataVersion)
	}
	if strings.TrimSpace(authData.AssetVersion) != "" {
		headers.Set("X-Asset-Version", authData.AssetVersion)
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header = headers

	client := &http.Client{}
	if strings.TrimSpace(issueCookie) != "" {
		client.Transport = &cookieRoundTripper{
			base:   http.DefaultTransport,
			cookie: issueCookie,
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	newToken := resp.Header.Get("X-Session-Token")
	if newToken == "" {
		newToken = sessionToken
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, newToken, err
	}
	if resp.StatusCode != http.StatusOK {
		decoded := decodePublishedErrorBody(raw, region)
		if strings.Contains(strings.ToLower(decoded), "session_error") {
			return nil, newToken, fmt.Errorf("published score request failed: status=%d %s (session token may be expired)", resp.StatusCode, decoded)
		}
		return nil, newToken, fmt.Errorf("published score request failed: status=%d %s", resp.StatusCode, decoded)
	}

	plain, err := decryptPayload(raw, region)
	if err != nil {
		return nil, newToken, err
	}
	return plain, newToken, nil
}

func extractUserNameFromSuite(raw []byte) string {
	var root map[string]any
	if err := msgpack.Unmarshal(raw, &root); err != nil {
		if jsonErr := json.Unmarshal(raw, &root); jsonErr != nil {
			return ""
		}
	}

	if gd := toMap(root["userGamedata"]); len(gd) > 0 {
		if name := toString(gd["name"]); name != "" {
			return name
		}
	}
	if before := toMap(root["beforeUserGamedata"]); len(before) > 0 {
		if name := toString(before["name"]); name != "" {
			return name
		}
	}
	if after := toMap(root["afterUserGamedata"]); len(after) > 0 {
		if name := toString(after["name"]); name != "" {
			return name
		}
	}
	if name := toString(root["name"]); name != "" {
		return name
	}

	return ""
}

func fetchUserNameByUserID(
	lookupUserID int,
	region string,
	sessionToken string,
	authData AuthData,
	appInfo AppInfo,
) (string, string, error) {
	if lookupUserID <= 0 {
		return "", sessionToken, fmt.Errorf("user id is required")
	}
	baseURL, ok := apiEndpoints[strings.ToLower(region)]
	if !ok {
		return "", sessionToken, fmt.Errorf("unsupported region: %s", region)
	}

	url := fmt.Sprintf("%s/api/suite/user/%d", baseURL, lookupUserID)
	raw, newToken, err := requestPublishedScoreInfo(http.MethodGet, url, region, sessionToken, authData, appInfo)
	if err != nil {
		return "", newToken, err
	}

	return extractUserNameFromSuite(raw), newToken, nil
}

func FetchPublishedScoreInfo(
	scoreID string,
	userID int,
	sessionToken string,
	region string,
	authData AuthData,
	appInfo AppInfo,
) (PublishedScoreInfo, [][]byte, string, error) {
	if strings.TrimSpace(scoreID) == "" {
		return PublishedScoreInfo{}, nil, sessionToken, fmt.Errorf("score id is required")
	}
	if userID <= 0 {
		return PublishedScoreInfo{}, nil, sessionToken, fmt.Errorf("user id is required")
	}
	if strings.TrimSpace(sessionToken) == "" {
		return PublishedScoreInfo{}, nil, sessionToken, fmt.Errorf("session token is required")
	}
	if strings.TrimSpace(appInfo.AppVersion) == "" || strings.TrimSpace(appInfo.AppHash) == "" {
		return PublishedScoreInfo{}, nil, sessionToken, fmt.Errorf("app info is required")
	}

	baseURL, ok := apiEndpoints[strings.ToLower(region)]
	if !ok {
		return PublishedScoreInfo{}, nil, sessionToken, fmt.Errorf("unsupported region: %s", region)
	}

	url := fmt.Sprintf("%s/api/user/%d/custom-music-score/published/search/%s", baseURL, userID, scoreID)
	raw, newToken, err := requestPublishedScoreInfo(http.MethodGet, url, region, sessionToken, authData, appInfo)
	if err != nil {
		return PublishedScoreInfo{}, nil, newToken, err
	}

	info, err := extractPublishedScoreInfo(raw)
	if err != nil {
		return PublishedScoreInfo{}, [][]byte{raw}, newToken, err
	}
	if strings.TrimSpace(info.UserCustomMusicScoreID) == "" {
		info.UserCustomMusicScoreID = strings.TrimSpace(scoreID)
	}
	if !info.IsPublished {
		info.IsPublished = true
	}

	lookupUserID := info.UserID
	if lookupUserID <= 0 {
		lookupUserID = userID
	}
	if info.UserName == "" && lookupUserID > 0 {
		if userName, newestToken, userErr := fetchUserNameByUserID(lookupUserID, region, newToken, authData, appInfo); userErr == nil {
			info.UserName = userName
			newToken = newestToken
		}
	}

	var supplementaryRaws [][]byte
	if lookupUserID > 0 && (info.FavoriteCount == 0 || info.Version == 0) {
		supplementURL := fmt.Sprintf("%s/api/user/%d/custom-music-score/published", baseURL, lookupUserID)
		if suppRaw, suppToken, suppErr := requestPublishedScoreInfo(http.MethodGet, supplementURL, region, newToken, authData, appInfo); suppErr == nil {
			newToken = suppToken
			supplementaryRaws = append(supplementaryRaws, suppRaw)
			mergePublishedListEntry(suppRaw, scoreID, &info)
		}
	}

	return info, append([][]byte{raw}, supplementaryRaws...), newToken, nil
}

func mergePublishedListEntry(raw []byte, scoreID string, info *PublishedScoreInfo) map[string]any {
	var root map[string]any
	if err := msgpack.Unmarshal(raw, &root); err != nil {
		if err2 := json.Unmarshal(raw, &root); err2 != nil {
			return nil
		}
	}

	// The list is typically at root["userCustomMusicScoreInfoJsons"] (plural) or root["data"]
	var entries []any
	for _, key := range []string{"userCustomMusicScoreInfoJsons", "userCustomMusicScoreInfoJson", "data", "list", "scores"} {
		if v, ok := root[key]; ok {
			switch t := v.(type) {
			case []any:
				entries = t
			case map[string]any:
				for _, sub := range t {
					if arr, ok := sub.([]any); ok {
						entries = arr
						break
					}
				}
			}
			if len(entries) > 0 {
				break
			}
		}
	}
	// Fall back to treating root itself as a single entry if no list found
	if len(entries) == 0 {
		entries = []any{root}
	}

	for _, raw := range entries {
		entry := toMap(raw)
		entryID := toString(entry["userCustomMusicScoreId"])
		if entryID != scoreID {
			// Also check inside nested InfoJson
			inner := toMap(entry["userCustomMusicScoreInfoJson"])
			entryID = toString(inner["userCustomMusicScoreId"])
			if entryID != scoreID {
				continue
			}
		}
		if info.FavoriteCount == 0 {
			if fc := toInt(entry["favoriteCount"]); fc != 0 {
				info.FavoriteCount = fc
			}
		}
		if info.Version == 0 {
			if v := toInt(entry["version"]); v != 0 {
				info.Version = v
			}
		}
		if info.UpdatedAt == 0 {
			if ua := toInt64(entry["updatedAt"]); ua != 0 {
				info.UpdatedAt = ua
			}
		}
		return root
	}
	return root
}
