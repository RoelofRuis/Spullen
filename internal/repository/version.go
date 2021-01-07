package repository

import "strconv"

type StorableVersionRepository struct {
	appVersion int
	storedVersion int
}

func NewVersionRepository() *StorableVersionRepository {
	return &StorableVersionRepository{appVersion: 1, storedVersion: 0}
}

func (s *StorableVersionRepository) GetApplicationVersion() int {
	return s.appVersion
}

func (s *StorableVersionRepository) GetStoredVersion() int {
	return s.storedVersion
}

// --- LOADING AND SAVING
// Ensuring it is a Storable
func (s *StorableVersionRepository) IsDirty() bool {
	return false
}

func (s *StorableVersionRepository) AfterPersist() {}

func (s *StorableVersionRepository) Instantiate(data []byte) error {
	versionBigint, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	version := int(versionBigint)

	s.storedVersion = version

	return nil
}

func (s *StorableVersionRepository) ToRaw() ([]byte, error) {
	return []byte(strconv.Itoa(s.appVersion)), nil
}
