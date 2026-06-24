package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/TootieJin/pjsekai-overlay-APPEND/pkg/scoremaker"
)

const DefaultCredentialFile = "credentials.json"

type CredentialEntry struct {
	UserID       int    `json:"user_id"`
	Credential   string `json:"credential"`
	Region       string `json:"region"`
	SessionToken string `json:"session_token"`
	DataVersion  string `json:"data_version,omitempty"`
	AssetVersion string `json:"asset_version,omitempty"`
	AssetHash    string `json:"asset_hash,omitempty"`
	SavedAt      int64  `json:"saved_at"`
}

type CredentialFile struct {
	Credentials []CredentialEntry `json:"credentials"`
	mu          sync.Mutex
	path        string
}

func Load(path string) (*CredentialFile, error) {
	if path == "" {
		path = DefaultCredentialFile
	}
	raw, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &CredentialFile{path: path}, nil
	}
	if err != nil {
		return nil, err
	}
	var cf CredentialFile
	if err := json.Unmarshal(raw, &cf); err != nil {
		return nil, fmt.Errorf("parse credential file: %w", err)
	}
	cf.path = path
	return &cf, nil
}

func (cf *CredentialFile) Save() error {
	cf.mu.Lock()
	defer cf.mu.Unlock()
	raw, err := json.MarshalIndent(cf, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(cf.path), 0755); err != nil {
		return err
	}
	tmp := cf.path + ".tmp"
	if err := os.WriteFile(tmp, raw, 0600); err != nil {
		return err
	}
	return os.Rename(tmp, cf.path)
}

func (cf *CredentialFile) FindForRegion(region string) *CredentialEntry {
	region = strings.ToLower(strings.TrimSpace(region))
	for i := range cf.Credentials {
		if strings.ToLower(strings.TrimSpace(cf.Credentials[i].Region)) == region {
			return &cf.Credentials[i]
		}
	}
	return nil
}

func (cf *CredentialFile) UpsertEntry(entry CredentialEntry) {
	cf.mu.Lock()
	defer cf.mu.Unlock()
	region := strings.ToLower(strings.TrimSpace(entry.Region))
	for i := range cf.Credentials {
		if cf.Credentials[i].UserID == entry.UserID &&
			strings.ToLower(strings.TrimSpace(cf.Credentials[i].Region)) == region {
			cf.Credentials[i] = entry
			return
		}
	}
	cf.Credentials = append(cf.Credentials, entry)
}

func EnsureGuest(region string, appInfo scoremaker.AppInfo, path string) (int, string, scoremaker.AuthData, *CredentialFile, error) {
	if path == "" {
		path = DefaultCredentialFile
	}
	cf, err := Load(path)
	if err != nil {
		return 0, "", scoremaker.AuthData{}, nil, fmt.Errorf("load credential file: %w", err)
	}
	entry := cf.FindForRegion(region)
	if entry == nil {
		fmt.Printf("- ゲストアカウントを新規作成中 (Registering new guest account for %s)...\n", strings.ToUpper(region))
		userID, credential, sessionToken, authData, err := registerAndAuth(region, appInfo)
		if err != nil {
			return 0, "", scoremaker.AuthData{}, nil, fmt.Errorf("register guest account: %w", err)
		}
		cf.UpsertEntry(CredentialEntry{
			UserID: userID, Credential: credential,
			Region: strings.ToLower(region), SessionToken: sessionToken,
			DataVersion: authData.DataVersion, AssetVersion: authData.AssetVersion,
			AssetHash: authData.AssetHash, SavedAt: time.Now().Unix(),
		})
		if saveErr := cf.Save(); saveErr != nil {
			fmt.Printf("- WARN: credentials.jsonファイルの保存に失敗 (Failed to save the credentials.json file): %s\n", saveErr)
		}
		fmt.Printf("- ゲストアカウント作成完了 (Guest account created): userId=%d\n", userID)
		return userID, sessionToken, authData, cf, nil
	}
	if strings.TrimSpace(entry.SessionToken) == "" {
		token, authData, authErr := scoremaker.AuthenticateWithCredential(entry.UserID, entry.Credential, region, appInfo)
		if authErr != nil {
			return 0, "", scoremaker.AuthData{}, nil, fmt.Errorf("authenticate guest: %w", authErr)
		}
		entry.SessionToken = token
		if authData.DataVersion != "" {
			entry.DataVersion = authData.DataVersion
			entry.AssetVersion = authData.AssetVersion
			entry.AssetHash = authData.AssetHash
		}
		entry.SavedAt = time.Now().Unix()
		cf.UpsertEntry(*entry)
		_ = cf.Save()
	}
	return entry.UserID, entry.SessionToken, scoremaker.AuthData{
		DataVersion: entry.DataVersion, AssetVersion: entry.AssetVersion, AssetHash: entry.AssetHash,
	}, cf, nil
}

func RefreshSession(region string, appInfo scoremaker.AppInfo, cf *CredentialFile) (int, string, scoremaker.AuthData, error) {
	entry := cf.FindForRegion(region)
	if entry == nil {
		return 0, "", scoremaker.AuthData{}, fmt.Errorf("no credential found for region %s", region)
	}
	token, authData, err := scoremaker.AuthenticateWithCredential(entry.UserID, entry.Credential, region, appInfo)
	if err != nil {
		return 0, "", scoremaker.AuthData{}, fmt.Errorf("re-authenticate: %w", err)
	}
	entry.SessionToken = token
	if authData.DataVersion != "" {
		entry.DataVersion = authData.DataVersion
		entry.AssetVersion = authData.AssetVersion
		entry.AssetHash = authData.AssetHash
	}
	entry.SavedAt = time.Now().Unix()
	cf.UpsertEntry(*entry)
	_ = cf.Save()
	return entry.UserID, token, authData, nil
}

func registerAndAuth(region string, appInfo scoremaker.AppInfo) (int, string, string, scoremaker.AuthData, error) {
	userID, credential, err := scoremaker.RegisterGuestAccount(region, appInfo)
	if err != nil {
		return 0, "", "", scoremaker.AuthData{}, err
	}
	sessionToken, authData, err := scoremaker.AuthenticateWithCredential(userID, credential, region, appInfo)
	if err != nil {
		return 0, "", "", scoremaker.AuthData{}, fmt.Errorf("authenticate after register: %w", err)
	}
	return userID, credential, sessionToken, authData, nil
}
