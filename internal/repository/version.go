package repository

import "strconv"

type StorableVersionManager struct {
	latestVersion int
	storedVersion int
}

func NewVersionManager(latestVersion int) *StorableVersionManager {
	return &StorableVersionManager{latestVersion: latestVersion, storedVersion: 0}
}

func (s *StorableVersionManager) GetVersion() int {
	return s.storedVersion
}

// --- LOADING AND SAVING
// Ensuring it is a Storable
func (s *StorableVersionManager) IsDirty() bool {
	return false
}

func (s *StorableVersionManager) AfterPersist() {}

func (s *StorableVersionManager) Instantiate(data []byte) error {
	if len(data) == 0 {
		s.storedVersion = s.latestVersion
		return nil
	}

	versionBigint, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	version := int(versionBigint)

	s.storedVersion = version

	return nil
}

func (s *StorableVersionManager) ToRaw() ([]byte, error) {
	return []byte(strconv.Itoa(s.latestVersion)), nil
}
